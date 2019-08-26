package cron

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/modules/checkstyle"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var c = cron.New()

// SetupCron sets up all configured cron jobs.
func SetupCron() {
	go checkRepositoriesCron()
	c.AddFunc(settings.Config.CheckActiveRepoCron, checkRepositoriesCron)
	c.Start()
}

func fillFields(message string, username string,
	reportLink string, commit string) string {
	filled := message
	filled = strings.Replace(filled, "$username", username, -1)
	filled = strings.Replace(filled, "$site_title", settings.Config.SiteTitle, -1)
	filled = strings.Replace(filled, "$report_link", reportLink, -1)
	filled = strings.Replace(filled, "$commit", commit, -1)
	return filled
}

func getUsernameFromProject(project string) string {
	return strings.Split(project, "/")[0]
}

func getReportLink(project string, commit string) string {
	return fmt.Sprintf("%s/%s/report/%s", settings.Config.SiteURL, project, commit)
}

// createNewIssue creates a new issue in the GitLab project's issue
// tracker with default description.
func createNewIssue(git *gitlab.Client, project string, commit string) (*gitlab.Issue, error) {
	reportLink := getReportLink(project, commit)

	title := fillFields(settings.Config.GitLabCustomisation.IssueTitle,
		getUsernameFromProject(project), reportLink, commit)

	desc := fillFields(settings.Config.GitLabCustomisation.GeneratingBody,
		getUsernameFromProject(project), reportLink, commit)

	issueOpt := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(desc),
	}
	issue, _, err := git.Issues.CreateIssue(project, issueOpt)
	return issue, err
}

// updateIssue updates an issue posting on GitLab with a link of the generated report.
func updateIssue(git *gitlab.Client, repo *models.Repository, report *models.Report) error {
	reportLink := getReportLink(repo.RepoID, report.Commit)

	title := fillFields(settings.Config.GitLabCustomisation.IssueTitle,
		getUsernameFromProject(repo.RepoID), reportLink, report.Commit)

	desc := fillFields(settings.Config.GitLabCustomisation.CompleteBody,
		getUsernameFromProject(repo.RepoID), reportLink, report.Commit)

	updateIssue := &gitlab.UpdateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(desc),
		StateEvent:  gitlab.String("reopen"),
	}
	_, _, err := git.Issues.UpdateIssue(repo.RepoID, repo.IssueID, updateIssue)
	return err
}

// getRepoArchive returns the zip archive of a GitLab project as a byte array.
func getRepoArchive(git *gitlab.Client, project string, sha string) ([]byte, error) {
	archiveOpt := &gitlab.ArchiveOptions{
		Format: gitlab.String("zip"),
		SHA:    gitlab.String(sha),
	}
	archive, _, err := git.Repositories.Archive(project, archiveOpt)
	return archive, err
}

// checkGitLabProjectExists returns whether a GitLab project exists and accessible.
func checkGitLabProjectExists(git *gitlab.Client, project string) (err error) {
	_, resp, err := git.Projects.GetProject(project, &gitlab.GetProjectOptions{})
	if err != nil && resp != nil && resp.StatusCode != 200 {
		return fmt.Errorf("Cannot access GitLab project, returned status code: %s", resp.Status)
	} else if err != nil {
		return err
	}
	return nil
}

// getCommits returns a list of GitLab commits of a specific project.
func getCommits(git *gitlab.Client, project string) (commits []*gitlab.Commit, err error) {
	commits, _, err = git.Commits.ListCommits(project, &gitlab.ListCommitsOptions{})
	if err != nil {
		return nil, err
	} else if len(commits) == 0 {
		return nil, errors.New("Project has no commits")
	}
	return commits, err
}

// downloadArchiveZip downloads a zip archive of a GitLab project and writes in in a
// temporary directory. This function runs getRepoArchive.
func downloadArchiveZip(git *gitlab.Client, project string, commit string, tempDir string) (path string, err error) {
	var arch []byte
	arch, err = getRepoArchive(git, project, commit)
	if err != nil {
		return "", err
	}
	archPath := fmt.Sprintf("%s/archive.zip", tempDir)
	err = ioutil.WriteFile(archPath, arch, 0644)
	return archPath, err
}

// unzipProjectArchive unzips a downloaded project archive and returns the extracted path.
// This function can only be run after downloadArciveZip runs successfully.
func unzipProjectArchive(project *settings.Project, tempDir string, archPath string, commit string) (path string, err error) {
	_, err = exec.Command("unzip", "-d", tempDir, archPath).Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s-%s-%s", tempDir, project.Project, commit, commit), err
}

// runCheckstyle runs checkstyle on a project and returns a report.
func runCheckstyle(project string, report *models.Report, checkPath string, commit string) (err error) {
	var o []byte
	o, err = exec.Command("java", "-jar", settings.Config.CheckstyleJarPath, "-c", settings.Config.CheckstyleConfigPath, checkPath).Output()
	if err != nil {
		return err
	}
	report.Status = models.Complete
	report.Issues = checkstyle.GetIssues(string(o), commit, checkPath, report.ReportID)

	return err
}

// createReport creates a new empty report.
func createReport(git *gitlab.Client, project string, commit string) (report *models.Report, err error) {
	report = &models.Report{
		RepositoryID: project,
		Commit:       commit,
		Status:       models.InProgress,
	}
	if models.HasReport(project, commit) {
		err = models.UpdateReport(report)
	} else {
		err = models.AddReport(report)
	}
	return report, err
}

func genSecretKey() (key string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 12; i++ {
		key += strconv.Itoa(rand.Intn(10))
	}
	return key
}

var checkRepoLimit = make(chan struct{}, 1)

// checkRepositoriesCron checks active git repositories, cron job.
func checkRepositoriesCron() {
	checkRepoLimit <- struct{}{} // limit running instances of checkRepositoriesCron
	log.Println("Cron: Starting to check active repositories")
	// NOTE: This means that if any changes
	// are applied to the list, it is lost.
	settings.LoadActiveProjs(false)
	git := settings.GetGitLabClient()
	var err error

	var wg sync.WaitGroup
	var workers = make(chan struct{}, settings.Config.Limits.MaxCheckWorkers)
	for _, proj := range settings.ActiveProjs.Projects {
		wg.Add(1)
		workers <- struct{}{}
		go func(proj settings.Project, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				<-workers
			}()
			path := proj.GetFullPath()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Cron recovery: %s: %s\n", path, r)
				}
			}()

			// Benchmarking
			start := time.Now()
			defer func(path string, start time.Time) {
				log.Printf("Cron: %s: Done checking project (%s)\n", path,
					time.Since(start))
			}(path, start)

			log.Printf("Cron: %s: Checking project...\n", path)
			var repo *models.Repository
			// Load Repository from database
			loadRepoMutex := make(chan bool)
			go func() {
				repo, err = models.GetRepo(path)
				if err != nil && err.Error() == "Repository does not exist" {
					repo.RepoID = path
					repo.SecretKey = genSecretKey()
					models.AddRepo(repo)
				} else if err != nil {
					log.Panicf("Failed to load repository: %s\n", err)
				}
				loadRepoMutex <- true
			}()

			// Check if GitLab project exists and/or accessible.
			if err = checkGitLabProjectExists(git, path); err != nil {
				log.Panicln(err)
			}

			<-loadRepoMutex
			close(loadRepoMutex)
			// Load project's commits
			var commits []*gitlab.Commit
			commits, err = getCommits(git, path)
			if checkCronErr(err, path, "") {
				return
			}

			// Check whether report is already generated for latest commit
			latestCommit := commits[0].ID
			_ = commits
			if rep, ok := repo.GetReport(latestCommit); ok &&
				rep.Status == models.Complete {
				log.Println("Report is already latest")
				return
			}
			repo.LatestCommit = latestCommit
			models.UpdateRepo(repo)

			// Create a GitLab issue, if doesn't exist
			createIssueMutex := make(chan bool)
			go func() {
				defer func() { createIssueMutex <- true }()
				if repo.IssueID == 0 { // Assumption: issue ID generated by GitLab is always > 0, true as of July 2019
					var issue *gitlab.Issue
					issue, err = createNewIssue(git, path, latestCommit)
					if err != nil {
						log.Panicf("Failed to create a new issue: %s\n", err)
					}
					repo.IssueID = issue.IID
					models.UpdateRepo(repo)
				}
			}()
			done := false
			var report *models.Report
			report, err = createReport(git, path, latestCommit)
			if checkCronErr(err, path, "Unable to create report in DB") {
				return
			}

			// In case the following steps fail before done, mark
			// the report as failed to inform user.
			defer func() {
				if !done {
					report.Status = models.Failed
					models.UpdateReport(report)
				}
			}()

			// Create a temporary directory
			tempDir, err := ioutil.TempDir("", fmt.Sprintf("learning-bot-%s-%s-%s", proj.Namespace, proj.Project,
				latestCommit[6:]))
			if checkCronErr(err, path, "Unable to create temporary directory (full disk?)") {
				return
			}
			defer func(path string) { // Defer the deletion of the temporary directory
				exec.Command("rm", "-rf", path).Output()
			}(tempDir)

			// Download archive as zip
			var archPath string
			archPath, err = downloadArchiveZip(git, path, latestCommit, tempDir)
			if checkCronErr(err, path, "Cannot download zip archive") {
				return
			}

			// Unzip project archive
			var newPath string
			newPath, err = unzipProjectArchive(&proj, tempDir, archPath, latestCommit)
			if checkCronErr(err, path, "Unable to unzip project") {
				return
			}

			// Run checkstyle

			err = runCheckstyle(path, report, newPath, latestCommit)
			if checkCronErr(err, path, "Unable to run checkstyle") {
				report.Status = models.Failed
				return
			}

			// Update issue
			go func() {
				<-createIssueMutex
				close(createIssueMutex)
				err = updateIssue(git, repo, report)
				if checkCronErr(err, path, "Cannot update GitLab issues") {
					return
				}
			}()

			// Update repo in DB
			if models.HasReport(path, report.Commit) {
				err = models.UpdateReport(report)
			} else {
				err = models.AddReport(report)
			}
			if checkCronErr(err, path, "Unable to update repository reports in DB") {
				return
			}

			err = models.UpdateIssues(report)
			if checkCronErr(err, path, "Unable to update report issues in database") {
				return
			}
			done = true
		}(proj, &wg)

	}
	wg.Wait() // wait for all concurrent processes to finish
	close(workers)
	log.Println("Cron: End of checking active repositories")
	<-checkRepoLimit
}

func checkCronErr(err error, project string, desc string) bool {
	if err != nil {
		if len(desc) > 0 {
			log.Printf("Cron: %s: %s\n", desc, err)
		} else {
			log.Printf("Cron: %s\n", err)
		}
		return true
	}
	return false
}
