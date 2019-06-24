package settings

import (
	"github.com/robfig/cron"
)

var c = cron.New()

// SetupCron sets up cron for checking active Git repositories.
func SetupCron() {
	c.AddFunc(Config.CheckActiveRepoCron, checkRepositories)
}

func checkRepositories() {
	// TODO
}
