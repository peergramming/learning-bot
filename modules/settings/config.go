package settings

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	WorkingDir      string
	ConfigPath      = "config.toml"
	ActiveProjsPath = "active-projects.toml"
	ActiveProjs     ActiveProjects
	Config          Configuration
)

type Configuration struct {
	SiteTitle             string          `toml:"site_title"`
	BotPrivateToken       string          `toml:"bot_private_access_token"`
	CheckstyleJarPath     string          `toml:"checkstyle_jar_path"`
	CheckstyleConfigPath  string          `toml:"checkstyle_config_path"`
	GitLabInstanceURL     string          `toml:"gitlab_instance_url"`
	DatabaseConfiguration DBConfiguration `toml:"database_configuration"`
	LMSTitle              string          `toml:"lms_title,omitempty"`
	LMSURL                string          `toml:"lms_url,omitempty"`
	CheckActiveRepoCron   string          `toml:"check_active_repositories_cron"`
}

type DBType int

const (
	SQLite = iota
	MySQL
)

type DBConfiguration struct {
	Type    DBType `toml:"type,string"`
	Host    string `toml:"host,omitempty"` // For MySQL...
	Name    string `toml:"name,omitempty"`
	User    string `toml:"user,omitempty"`
	SSLMode string `toml:"ssl_mode,omitempty"`
	Path    string `toml:"path,omitempty"` // For SQLite
}

func NewConfiguration(token string, instance string, checkstyleJar string,
	databaseConfig DBConfiguration) Configuration {
	return Configuration{
		SiteTitle:             "Learning Bot",
		BotPrivateToken:       token,
		GitLabInstanceURL:     instance,
		CheckstyleJarPath:     checkstyleJar,
		CheckstyleConfigPath:  "./assets/checkstyle-lb.xml",
		DatabaseConfiguration: databaseConfig,
		LMSTitle:              "Vision",
		LMSURL:                "https://vision.hw.ac.uk",
		CheckActiveRepoCron:   "@every 1h45m",
	}
}

func init() {
	var err error
	WorkingDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory! ", err)
	}
}

func LoadConfig() {
	var err error
	if _, err = toml.DecodeFile(WorkingDir+"/"+ConfigPath, &Config); err != nil {
		log.Panicf("Failed to load the configuration file! Make sure you generate a configuration first! Error: %s", err)
	}

	LoadActiveProjs(false)
}
