package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"log"
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

// checkRepositories checks active git repositories, cron job.
func checkRepositoriesCron() {
	log.Println("Cron: Starting to check active repositories")
	// NOTE: This means that if any changes
	// are applied to the list, it is lost.
	settings.LoadActiveProjs(false)
	git := settings.GetGitLabClient()

	for _, proj := range settings.ActiveProjs.Projects { // TODO concurrent checking
		path := proj.GetFullPath()
		log.Printf("Cron: Checking project: %s\n", path)

		repo, err := models.GetRepo(path)
		if err != nil && err.Error() == "Repository does not exist" {
			// Repository does not exist, let's create an issue.
			var issue *gitlab.Issue
			issue, err = createNewIssue(git, path) // TODO: mechanism if a repo is deleted and recreated
			if err != nil {
				log.Fatalf("Cron: Failed to create a new issue for repository %s: %s",
					path, err)
				continue
			}
			repo.IssueID = issue.ID
			models.AddRepo(&repo)
		} else if err != nil {
			log.Fatalf("Cron: Failed to load repository %s: %s", path, err)
			continue
		}

		// TODO: Run test, and update issue
	}
	log.Println("Cron: End of checking active repositories")
}
