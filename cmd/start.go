package cmd

import (
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
	"html/template"
	"log"
	"net/http"

	"gitlab.com/gitedulab/learning-bot/models"
	"gitlab.com/gitedulab/learning-bot/modules/cron"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"gitlab.com/gitedulab/learning-bot/modules/utils"
	"gitlab.com/gitedulab/learning-bot/routes"
)

// CmdStart represents a command-line command
// which starts the bot.
var CmdStart = cli.Command{
	Name:    "run",
	Aliases: []string{"start", "web"},
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "no-web", Usage: "Does not run the web server"},
	},
	Usage:  "Start the learning bot",
	Action: start,
}

func start(clx *cli.Context) (err error) {
	if err = utils.SystemPackagesCheck(); err != nil {
		panic(err)
	}
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()
	cron.SetupCron()

	if clx.Bool("no-web") {
		log.Println("Running cron-only without web server.")
		select {}
	} else {
		// Run macaron
		m := macaron.Classic()
		funcMap := []template.FuncMap{map[string]interface{}{
			"Spacify": utils.Spacify,
		}}

		m.Use(macaron.Renderer(macaron.RenderOptions{
			Funcs: funcMap,
		}))

		// Web routes
		m.Get("/", routes.HomepageHandler)
		m.Get("/help/:check", routes.HelpCheckHandler)

		// Project specific routes
		m.Group("/:namespace/:project", func() {
			m.Get("/report/:sha", routes.ReportPageHandler)
			m.Get("/status/:sha.json", routes.APIGetReportStatusHandler)
		})

		log.Printf("Starting web server on port %s\n", settings.Config.SitePort)
		log.Fatal(http.ListenAndServe(settings.Config.SitePort, m))
	}
	return nil
}
