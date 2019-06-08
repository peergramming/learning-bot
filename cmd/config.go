package cmd

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"io/ioutil"
	"log"
	"os"
)

var CmdConfig = cli.Command{
	Name:   "config",
	Usage:  "Create a new configuration file",
	Action: runConfig,
}

var (
	defaultInstance       = "https://gitlab.com"
	defaultCheckstylePath = "checkstyle-8.21-all.jar"
)

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

	fmt.Printf("Enter your GitLab instance URL (incl. protocol scheme): [%s] ", defaultInstance)
	var instance string
	fmt.Scanln(&instance)
	if instance == "" {
		instance = defaultInstance
	}

	fmt.Println("You have to generate a GitLab personal access token with at least the following scopes:")
	fmt.Println("\tapi, read_user, read_repository, write_repository")
	fmt.Println("A token can be generated at: https://gitlab.com/profile/personal_access_tokens")
	fmt.Printf("Enter your GitLab personal access token: ")
	var token string
	fmt.Scanln(&token)

	// TODO validate the instance and token before continuing...

	// Select checkstyle jar and config path
	fmt.Println("This program requires checkstyle to generate reports")
	fmt.Println("Checkstyle can be downloaded from https://github.com/checkstyle/checkstyle/releases")
	fmt.Printf("Enter the path of the checkstyle jar file: [%s] ", defaultCheckstylePath)
	var checkstylePath string
	fmt.Scanln(&checkstylePath)
	if checkstylePath == "" {
		checkstylePath = defaultCheckstylePath
	}

	config := settings.Configuration{
		BotPrivateToken:      token,
		GitLabInstanceURL:    instance,
		CheckstyleJarPath:    checkstylePath,
		CheckstyleConfigPath: "./assets/checkstyle-lb.xml",
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile("config.toml", buf.Bytes(), 0644)
	if err2 != nil {
		log.Fatal(err2)
	}

	return nil
}
