package settings

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	ConfigPath string
	Config     Configuration
)

type Configuration struct {
	BotPrivateToken       string
	CheckstyleJarPath     string
	CheckstyleConfigPath  string
	GitLabInstanceURL     string
	DatabaseConfiguration DBConfiguration
}

type DBType int

const (
	SQLite = iota
	MySQL
)

type DBConfiguration struct {
	Type    DBType
	Host    string
	Name    string
	User    string
	SSLMode string
	Path    string // For SQLite
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory! ", err)
	}

	ConfigPath = wd + "/config.toml"
	if _, err2 := toml.DecodeFile(ConfigPath, &Config); err2 != nil {
		log.Fatal("Failed to load the configuration file! Make sure you generate a configuration first!\n", err2)
	}
}
