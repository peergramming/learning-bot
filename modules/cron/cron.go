package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/models/checkstyle"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"io/ioutil"
	"log"
	"os/exec"
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

func updateIssue(git *gitlab.Client, repo *models.Repository, report *models.Report) error {
	link := fmt.Sprintf("%s/%s/report/%s", settings.Config.SiteURL, repo.RepoID, report.Commit)
	updateIssue := &gitlab.UpdateIssueOptions{
		Description: gitlab.String(fmt.Sprintf("Hey!\n\nReport has been generated on commit %s.\n\n[View report](%s)", report.Commit, link)),
		StateEvent:  gitlab.String("reopen"),
	}
	fmt.Println(repo.RepoID, repo.IssueID)
	_, _, err := git.Issues.UpdateIssue(repo.RepoID, repo.IssueID, updateIssue)
	return err
}

func getRepoArchive(git *gitlab.Client, project string, sha string) ([]byte, error) {
	archiveOpt := &gitlab.ArchiveOptions{
		Format: gitlab.String("zip"),
		SHA:    gitlab.String(sha),
	}
	archive, _, err := git.Repositories.Archive(project, archiveOpt)
	return archive, err
}

// checkRepositories checks active git repositories, cron job.
func checkRepositoriesCron() {
	log.Println("Cron: Starting to check active repositories")
	// NOTE: This means that if any changes
	// are applied to the list, it is lost.
	settings.LoadActiveProjs(false)
	git := settings.GetGitLabClient()
	var err error

	for _, proj := range settings.ActiveProjs.Projects { // TODO concurrent checking
		path := proj.GetFullPath()

		// Benchmarking
		start := time.Now()
		defer func(path string, start time.Time) {
			elapsed := time.Since(start)
			log.Printf("Cron: %s: Done checking project (%s)\n", path, elapsed)
		}(path, start)

		log.Printf("Cron: %s: Checking project...\n", path)
		var repo *models.Repository

		// Load Repository from database
		repo, err = models.GetRepo(path)
		if err != nil && err.Error() == "Repository does not exist" {
			repo.RepoID = path
			models.AddRepo(repo)
		} else if err != nil {
			log.Printf("Cron: %s: Failed to load repository: %s\n", path, err)
			continue
		}

		fmt.Printf("Cron: %s: Has %d reports\n", path, len(repo.Reports))

		// Load project's commits
		var commits []*gitlab.Commit
		commits, resp, err := git.Commits.ListCommits(path, &gitlab.ListCommitsOptions{})
		if resp.StatusCode != 200 {
			log.Printf("Cron: %s: Cannot access commits, response code %s\n", path, resp.Status)
			continue
		} else if len(commits) == 0 {
			log.Printf("Cron: %s: Project has no commits, cannot proceed...\n", path)
			continue
		}

		// Check whether report is already generated for latest commit
		latestCommit := commits[0]
		if _, ok := repo.GetReport(latestCommit.ID); ok {
			log.Printf("Cron: %s: Report is already latest\n", path)
			continue
		}

		// Create a GitLab issue, if doesn't exist
		if repo.IssueID == 0 {
			var issue *gitlab.Issue
			issue, err = createNewIssue(git, path)
			if err != nil {
				log.Printf("Cron: %s: Failed to create a new issue: %s\n",
					path, err)
				continue
			}
			repo.IssueID = issue.IID
			models.UpdateRepo(repo)
		}

		// Download archive as zip
		var arch []byte
		arch, err = getRepoArchive(git, path, latestCommit.ID)
		dir, err := ioutil.TempDir("", fmt.Sprintf("learning-bot-%s-%s-%s", proj.Namespace, proj.Project, latestCommit.ID[6:]))
		if err != nil {
			log.Printf("Cron: %s: Cannot create a temporary directory, is disk space full?\n", path)
			continue
		}
		archPath := fmt.Sprintf("%s/archive.zip", dir)
		ioutil.WriteFile(archPath, arch, 0644)
		var o []byte

		// Unzip project archive
		o, err = exec.Command("unzip", "-d", dir, archPath).Output()
		if err != nil {
			log.Printf("Cron: %s: Failed to extract project archive: %s %s\n", path, err, o)
			continue
		}
		newPath := fmt.Sprintf("%s/%s-%s-%s", dir, proj.Project, latestCommit.ID, latestCommit.ID)
		log.Printf("Cron: %s: Downloaded: %s\n", path, newPath)

		// TODO: Run checkstyle test, parse, and update issue

		// Run checkstyle
		o, err = exec.Command("java", "-jar", settings.Config.CheckstyleJarPath, "-c", settings.Config.CheckstyleConfigPath, newPath).Output()
		if err != nil {
			log.Printf("Cron: %s: Cannot run checkstyle on project: %s, %s\n", path, err, o)
			continue
		}
		report := checkstyle.GenerateReport(string(o), latestCommit.ID, newPath)
		report.RepositoryID = path
		reports := append(repo.Reports, &report)

		err = updateIssue(git, repo, &report)
		if err != nil {
			log.Printf("Cron: %s: Cannot update issue: %s\n", path, err)
			continue
		}
		err = models.UpdateRepositoryReports(repo, reports)
		if err != nil {
			log.Printf("Cron: %s: Cannot update repository reports in db: %s\n", path, err)
			continue
		}

	}
	log.Println("Cron: End of checking active repositories")
}
