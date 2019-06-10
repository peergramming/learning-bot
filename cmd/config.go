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
	"strings"
)

var CmdConfig = cli.Command{
	Name:   "config",
	Usage:  "Create a new configuration file",
	Action: runConfig,
}

var (
	defaultInstance       = "https://gitlab.com"
	defaultCheckstylePath = "checkstyle-8.21-all.jar"
	defaultDriver         = "sqlite"
	defaultDBPath         = "./data.db"
	defaultDBHost         = "localhost:3306"
	defaultDBSSLMode      = "skip-verify"
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
	checkstylePath := scanWithDefault(defaultCheckstylePath)

	var dbConfig settings.DBConfiguration

	fmt.Println("This program supports the following database drivers:")
	fmt.Println("mysql, sqlite")
	fmt.Printf("Select a database driver: [%s] ", defaultDriver)
	dbDriver := scanWithDefault(defaultDriver)

	var dbDriverType settings.DBType
	switch strings.ToLower(dbDriver) {
	case "sqlite":
		dbConfig.Type = settings.SQLite
	case "mysql":
		dbConfig.Type = settings.MySQL
	default:
		dbConfig.Type = settings.SQLite
	}

	if dbDriverType == settings.SQLite {
		fmt.Printf("Enter a path for the SQLite file: [%s] ", defaultDBPath)
		dbConfig.Path = scanWithDefault(defaultDBPath)
	} else if dbDriverType == settings.MySQL {
		fmt.Printf("Enter the host of the MySQL server (incl. port): [%s] ", defaultDBHost)
		dbConfig.Host = scanWithDefault(defaultDBHost)

		fmt.Printf("Enter the database name to use for the MySQL server: ")
		dbConfig.Name = scanWithDefault("")

		fmt.Printf("Enter the username to use for the MySQL server: ")
		dbConfig.User = scanWithDefault("")

		fmt.Println("Select the TLS mode to use, the following values are valid:")
		fmt.Println("true, false, skip-verify, preferred, <name>")
		fmt.Printf("Enter the TLS mode to use for the MySQL server: [%s] ", defaultDBSSLMode)
		dbConfig.SSLMode = scanWithDefault(defaultDBSSLMode)
	}

	// Generate struct configuration

	config := settings.Configuration{
		BotPrivateToken:       token,
		GitLabInstanceURL:     instance,
		CheckstyleJarPath:     checkstylePath,
		CheckstyleConfigPath:  "./assets/checkstyle-lb.xml",
		DatabaseConfiguration: dbConfig,
	}

	// Write to file

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

func scanWithDefault(defaultVal string) string {
	var temp string
	fmt.Scanln(&temp)
	if temp == "" {
		temp = defaultVal
	}
	return temp
}
