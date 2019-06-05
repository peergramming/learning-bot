package cmd

import (
	"github.com/urfave/cli"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	macaron "gopkg.in/macaron.v1"
)

var CmdStart = cli.Command{
	Name:    "run",
	Aliases: []string{"start"},
	Usage:   "Start the learning bot",
	Action:  start,
}

func start(clx *cli.Context) error {
	// Load configuration
	settings.LoadConfig()

	// XORM initialisation

	// Run macaron
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	m.Run()
	return nil
}
