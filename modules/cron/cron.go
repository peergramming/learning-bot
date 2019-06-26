package cron

import (
	"github.com/robfig/cron"
	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"log"
)

var c = cron.New()

// SetupCron sets up cron for checking active Git repositories.
func SetupCron() {
	c.AddFunc(settings.Config.CheckActiveRepoCron, checkRepositories)
}

func checkRepositories() {
	// TODO
	settings.LoadActiveProjs(false) // NOTE: This means that if any changes
	// are applied to the list, it is lost.

	for _, proj := range settings.ActiveProjs.Projects {
		// See if project exists in DB, if not, create it.

		repo, err := models.GetRepo(proj.GetFullPath())
		_ = repo // TEMP
		if err.Error() == "Repository does not exist" {
			// Repository does not exist, let's create an issue.

		} else if err != nil {
			// Some other problem
			log.Fatalf("Failed to load repository %s: %s", proj.GetFullPath(), err)
			continue
		}
	}
}
