package main

import (
	"gitlab.com/gitedulab/learning-bot/cmd"
	"log"
	"os"

	"github.com/urfave/cli"
)

var VERSION = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "learning-bot"
	app.Usage = "a GitLab bot for providing advice from code repair tools."
	app.Version = VERSION
	app.Commands = []cli.Command{
		cmd.CmdStart,
		cmd.CmdConfig,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
