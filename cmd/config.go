package cmd

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/peergramming/learning-bot/modules/settings"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// CmdConfig represents a command-line command
// which create a new configuration file.
var CmdConfig = cli.Command{
	Name:   "config",
	Usage:  "Create a new configuration file",
	Action: runConfig,
}

var (
	defaultInstance       = "https://gitlab.com"
	defaultCheckstylePath = "checkstyle-8.22-all.jar"
	defaultDriver         = "sqlite"
	defaultDBPath         = "./data.db"
	defaultDBHost         = "localhost:3306"
	defaultCertFile       = "server.crt"
	defaultKeyFile        = "server.key"
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

	fmt.Printf("Do you want to use TLS (HTTPS)? [n] ")
	var respTLS string
	fmt.Scanln(&respTLS)
	tlsConfig := settings.TLSServerConfiguration{Enabled: false}
	if strings.ToLower(respTLS) == "y" {
		tlsConfig.Enabled = true
		fmt.Printf("Enter the path of the certificate file: [%s] ", defaultCertFile)
		tlsConfig.CertFile = scanWithDefault(defaultCertFile)
		fmt.Printf("Enter the path of the key file: [%s] ", defaultKeyFile)
		tlsConfig.KeyFile = scanWithDefault(defaultKeyFile)
	}

	fmt.Printf("Enter your bot site URL (incl. protocol and port): ")
	var siteURL string
	fmt.Scanln(&siteURL)

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
	fmt.Println("mysql, sqlite, postgres")
	fmt.Printf("Select a database driver: [%s] ", defaultDriver)
	dbDriver := scanWithDefault(defaultDriver)

	var dbDriverType settings.DBType
	switch strings.ToLower(dbDriver) {
	case "sqlite":
		dbConfig.Type = settings.SQLite
	case "mysql":
		dbConfig.Type = settings.MySQL
	case "postgres":
		dbConfig.Type = settings.PostgreSQL
	default:
		dbConfig.Type = settings.SQLite
	}

	if dbDriverType == settings.SQLite {
		fmt.Printf("Enter a path for the SQLite file: [%s] ", defaultDBPath)
		dbConfig.Path = scanWithDefault(defaultDBPath)
	} else if dbDriverType == settings.MySQL || dbDriverType == settings.PostgreSQL {
		fmt.Printf("Enter the host of the SQL server (incl. port): [%s] ", defaultDBHost)
		dbConfig.Host = scanWithDefault(defaultDBHost)

		fmt.Printf("Enter the database name to use for the SQL server: ")
		dbConfig.Name = scanWithDefault("")

		fmt.Printf("Enter the username to use for the SQL server: ")
		dbConfig.User = scanWithDefault("")
		var defaultSSLMode string
		if dbDriverType == settings.MySQL {
			fmt.Println("Select the TLS mode to use, the following values are valid:")
			fmt.Println("true, false, skip-verify, preferred, <name>")
			defaultSSLMode = "skip-verify"
		} else if dbDriverType == settings.PostgreSQL {
			fmt.Println("Select the SSL mode to use, the following values are valid:")
			fmt.Println("disable, require, verify-ca, verify-full")
			defaultSSLMode = "require"
		}
		fmt.Printf("Enter the TLS mode to use for the SQL server: [%s] ", defaultSSLMode)
		dbConfig.SSLMode = scanWithDefault(defaultSSLMode)
		var envPassword string
		if dbDriverType == settings.MySQL {
			envPassword = "MYSQL_PASSWORD"
		} else if dbDriverType == settings.PostgreSQL {
			envPassword = "POSTGRESQL_PASSWORD"
		}
		fmt.Printf("The SQL password must be set in the '%s' environment variable\n", envPassword)
	}

	// Generate struct configuration
	config := settings.NewConfiguration(token, siteURL, instance, checkstylePath, dbConfig,
		tlsConfig)

	// Write to file
	var err error

	buf := new(bytes.Buffer)
	if err = toml.NewEncoder(buf).Encode(config); err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(settings.ConfigPath, buf.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Configuration saved at %s/%s!\n", settings.WorkingDir, settings.ConfigPath)
	fmt.Println("Review the configuration file to confirm the configuration and further customise the bot.")

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
