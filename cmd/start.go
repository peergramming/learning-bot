package cmd

import (
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"gitlab.com/gitedulab/learning-bot/routes"
)

var CmdStart = cli.Command{
	Name:    "run",
	Aliases: []string{"start"},
	Usage:   "Start the learning bot",
	Action:  start,
}

func start(clx *cli.Context) error {
	settings.LoadConfig()
	settings.SetupCron()

	// Run macaron
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// Web routes
	m.Get("/", routes.HomepageHandler)
	m.Get("/help/:check", routes.HelpCheckHandler)

	// Project specific routes
	m.Group("/:namespace/:project", func() {
		m.Get("/report/:sha", routes.ReportPageHandler)
		m.Get("/status/:sha.json", routes.APIGetReportStatusHandler)
	})

	m.Run()
	return nil
}
