package settings

import (
	"github.com/BurntSushi/toml"
	"log"
)

var (
	ConfigPath = "./config.toml"
	Config     Configuration
)

type Configuration struct {
	BotPrivateToken      string
	CheckstyleJarPath    string
	CheckstyleConfigPath string
	GitLabInstanceURL    string
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

func LoadConfig() {
	if _, err := toml.DecodeFile(ConfigPath, &Config); err != nil {
		log.Println("Failed to load the configuration file!")
		log.Fatal(err)
	}

}
