package main

import (
	"github.com/peergramming/learning-bot/cmd"
	"log"
	"os"

	"github.com/urfave/cli"
)

// VERSION specifies the version of learning-bot
var VERSION = "0.3.0"

func main() {
	app := cli.NewApp()
	app.Name = "learning-bot"
	app.Usage = "a GitLab bot for providing advice from code repair tools."
	app.Version = VERSION
	app.Commands = []cli.Command{
		cmd.CmdStart,
		cmd.CmdConfig,
		cmd.CmdManage,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
