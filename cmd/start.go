package cmd

import (
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/gitedulab/learning-bot/routes"
)

var CmdStart = cli.Command{
	Name:    "run",
	Aliases: []string{"start"},
	Usage:   "Start the learning bot",
	Action:  start,
}

func start(clx *cli.Context) error {
	// XORM initialisation

	// Run macaron
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// Web routes
	m.Get("/", routes.HomepageHandler)

	// API routes; closely resemble GitLab's API
	m.Group("/api/v1", func() {
		m.Get("/project/:id", routes.APIGetProjectHandler)
		m.Get("/project/:id/generate_report", routes.APIGenReportHandler)
		m.Get("/project/:id/reports", routes.APIGetReportsHandler)
	})

	m.Run()
	return nil
}
