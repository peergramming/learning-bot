// Package settings manages everything which relates to configuring
// the learning bot. It manages everything related to the configuration of
// the learning bot, cron jobs, and the GitLab API client.
package settings

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"time"
)

var (
	// WorkingDir holds the current working directory of the application.
	WorkingDir string
	// ConfigPath holds the configuration path.
	ConfigPath = "config.toml"
	// Config holds the entire application configuration.
	Config Configuration
)

// Configuration represents an entire configuration file
// of the learning bot. This excludes the ActiveProjects
// configuration.
type Configuration struct {
	SiteTitle               string          `toml:"site_title"`
	SiteURL                 string          `toml:"site_url"`
	BotPrivateToken         string          `toml:"bot_private_access_token"`
	CheckstyleJarPath       string          `toml:"checkstyle_jar_path"`
	CheckstyleConfigPath    string          `toml:"checkstyle_config_path"`
	GitLabInstanceURL       string          `toml:"gitlab_instance_url"`
	DatabaseConfiguration   DBConfiguration `toml:"database_configuration"`
	LMSTitle                string          `toml:"lms_title,omitempty"`
	LMSURL                  string          `toml:"lms_url,omitempty"`
	CheckActiveRepoCron     string          `toml:"check_active_repositories_cron"`
	MaxCheckWorkers         int             `toml:"max_check_workers"`
	TimezoneName            string          `toml:"timezone"`
	Timezone                *time.Location  `toml:"-"`
	CodeSnippetIncludeLines int             `toml:"code_snippet_include_previous_lines"`
}

// DBType represents the database driver type, such as MySQL or SQLite.
type DBType int

const (
	// SQLite is a serverless and file-based SQL driver.
	SQLite = iota
	// MySQL is a standard SQL driver.
	MySQL
)

// DBConfiguration represents a database configuration, including whether
// it is a MySQL or SQLite configuration.
type DBConfiguration struct {
	Type    DBType `toml:"type,string"`
	Host    string `toml:"host,omitempty"` // For MySQL...
	Name    string `toml:"name,omitempty"`
	User    string `toml:"user,omitempty"`
	SSLMode string `toml:"ssl_mode,omitempty"`
	Path    string `toml:"path,omitempty"` // For SQLite
}

// NewConfiguration creates a new configuration struct with default
// fields prefilled.
func NewConfiguration(token string, siteURL string, instance string, checkstyleJar string,
	databaseConfig DBConfiguration) Configuration {
	return Configuration{
		SiteTitle:               "Learning Bot",
		SiteURL:                 siteURL,
		BotPrivateToken:         token,
		GitLabInstanceURL:       instance,
		CheckstyleJarPath:       checkstyleJar,
		CheckstyleConfigPath:    "./assets/checkstyle-lb.xml",
		DatabaseConfiguration:   databaseConfig,
		LMSTitle:                "Vision",
		LMSURL:                  "https://vision.hw.ac.uk",
		CheckActiveRepoCron:     "@every 1h45m",
		MaxCheckWorkers:         5,
		TimezoneName:            "Europe/London",
		CodeSnippetIncludeLines: 3,
	}
}

func init() {
	var err error
	WorkingDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory! ", err)
	}
}

// LoadConfig loads the configuration from file, then passively loading
// ActiveProjects.
func LoadConfig() {
	var err error
	if _, err = toml.DecodeFile(WorkingDir+"/"+ConfigPath, &Config); err != nil {
		log.Panicf("Failed to load the configuration file! Make sure you generate a configuration first! Error: %s", err)
	}
	Config.Timezone, err = time.LoadLocation(Config.TimezoneName)
	if err != nil {
		log.Panicf("Invalid timezone in config: %s", err)
	}

	LoadActiveProjs(false)
}
