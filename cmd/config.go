package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"os"
)

var CmdConfig = cli.Command{
	Name:   "config",
	Usage:  "Create a new configuration file",
	Action: runConfig,
}

func runConfig(clx *cli.Context) error {
	if _, err := os.Stat(settings.ConfigPath); err == nil {
		fmt.Printf("A configuration file already exists. Are you sure you want to continue and replace the current configuration? (y/n) [n] ")
		var resp string
		fmt.Scanln(&resp)
		if resp != "y" {
			fmt.Println("Exiting...")
			return nil
		}
	}

	return nil
}
