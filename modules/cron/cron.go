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
	"os/exec"
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

// createNewIssue creates a new issue in the GitLab project's issue
// tracker with default description.
func createNewIssue(git *gitlab.Client, project string) (*gitlab.Issue, error) {
	issueOpt := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(fmt.Sprintf("[%s] Your code report 📊", settings.Config.SiteTitle)),
		Description: gitlab.String("Hey!\n\nYou report is currently being generated.\n\nSit tight!"),
	}
	issue, _, err := git.Issues.CreateIssue(project, issueOpt)
	return issue, err
}

// updateIssue updates an issue posting on GitLab with a link of the generated report.
func updateIssue(git *gitlab.Client, repo *models.Repository, report *models.Report) error {
	link := fmt.Sprintf("%s/%s/report/%s", settings.Config.SiteURL, repo.RepoID, report.Commit)
	updateIssue := &gitlab.UpdateIssueOptions{
		Description: gitlab.String(fmt.Sprintf("Hey!\n\nReport has been generated on commit %s.\n\n[View report](%s)", report.Commit, link)),
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
	if err != nil && resp.StatusCode != 200 {
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
	err = models.AddReport(report)
	return report, err
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
	var workers = make(chan struct{}, settings.Config.MaxCheckWorkers)
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
				elapsed := time.Since(start)
				log.Printf("Cron: %s: Done checking project (%s)\n", path, elapsed)
			}(path, start)

			log.Printf("Cron: %s: Checking project...\n", path)
			var repo *models.Repository
			// Load Repository from database
			loadRepoMutex := make(chan bool)
			go func() {
				repo, err = models.GetRepo(path)
				if err != nil && err.Error() == "Repository does not exist" {
					repo.RepoID = path
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
			if err != nil {
				log.Panicln(err)
			}

			// Check whether report is already generated for latest commit
			latestCommit := commits[0]
			if _, ok := repo.GetReport(latestCommit.ID); ok {
				log.Panicln("Report is already latest")
			}

			// Create a GitLab issue, if doesn't exist
			createIssueMutex := make(chan bool)
			go func() {
				defer func() { createIssueMutex <- true }()
				if repo.IssueID == 0 { // Assumption: issue ID generated by GitLab is always > 0, true as of July 2019
					var issue *gitlab.Issue
					issue, err = createNewIssue(git, path)
					if err != nil {
						log.Panicf("Failed to create a new issue: %s\n", err)
					}
					repo.IssueID = issue.IID
					models.UpdateRepo(repo)
				}
			}()
			_ = repo

			// Create a temporary directory
			tempDir, err := ioutil.TempDir("", fmt.Sprintf("learning-bot-%s-%s-%s", proj.Namespace, proj.Project,
				latestCommit.ID[6:]))
			if err != nil {
				log.Panicf("Unable to create temporary directory (full disk?): %s\n", err)
			}
			defer func(path string) { // Defer the deletion of the temporary directory
				exec.Command("rm", "-rf", path).Output()
			}(tempDir)

			// Download archive as zip
			var archPath string
			archPath, err = downloadArchiveZip(git, path, latestCommit.ID, tempDir)
			if err != nil {
				log.Panicf("Cannot download zip archive: %s\n", err)
			}

			// Unzip project archive
			var newPath string
			newPath, err = unzipProjectArchive(&proj, tempDir, archPath, latestCommit.ID)
			if err != nil {
				log.Panicf("Unable to unzip project: %s\n", err)
			}

			// Run checkstyle
			var report *models.Report
			report, err = createReport(git, path, latestCommit.ID)
			if err != nil {
				log.Panicf("Unable to create report in DB: %s\n", err)
			}

			err = runCheckstyle(path, report, newPath, latestCommit.ID)
			if err != nil {
				log.Panicf("Unable to run checkstyle: %s\n", err)
			}
			reports := append(repo.Reports, report)

			// Update issue
			go func() {
				<-createIssueMutex
				close(createIssueMutex)
				err = updateIssue(git, repo, report)
				if err != nil {
					log.Panicf("Cannot update issue: %s\n", err)
				}
			}()

			// Update repo in DB
			err = models.UpdateRepositoryReports(repo, reports)
			if err != nil {
				log.Panicf("Cannot update repository reports in db: %s\n", err)
			}
		}(proj, &wg)

	}
	wg.Wait() // wait for all concurrent processes to finish
	log.Println("Cron: End of checking active repositories")
	<-checkRepoLimit
}
