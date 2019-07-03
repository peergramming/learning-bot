package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"log"
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
		start := time.Now()
		path := proj.GetFullPath()
		log.Printf("Cron: %s: Checking project...\n", path)
		var repo models.Repository

		// Load Repository from database
		repo, err = models.GetRepo(path)
		if err != nil && err.Error() == "Repository does not exist" {
			models.AddRepo(&repo)
		} else if err != nil {
			log.Fatalf("Cron: %s: Failed to load repository: %s\n", path, err)
			continue
		}

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

		log.Printf("Latest commit: %s\n", commits[0].ID)

		// Create a GitLab issue
		var issue *gitlab.Issue
		issue, err = createNewIssue(git, path)
		if err != nil {
			log.Fatalf("Cron: %s: Failed to create a new issue: %s\n",
				path, err)
			continue
		}

		repo.IssueID = issue.ID
		models.UpdateRepo(&repo)

		// TODO: Run test, and update issue

		elapsed := time.Since(start)
		log.Printf("Cron: %s: Done checking project (%s)\n", path, elapsed)
	}
	log.Println("Cron: End of checking active repositories")
}
