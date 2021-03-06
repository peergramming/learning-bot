// Package settings manages everything which relates to configuring
// the learning bot. It manages everything related to the configuration of
// the learning bot, cron jobs, and the GitLab API client.
package settings

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strings"
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
	SiteTitle                string                   `toml:"site_title"`
	SiteURL                  string                   `toml:"site_url"`
	SitePort                 string                   `toml:"-"`
	BotPrivateToken          string                   `toml:"bot_private_access_token"`
	CheckstyleJarPath        string                   `toml:"checkstyle_jar_path"`
	CheckstyleConfigPath     string                   `toml:"checkstyle_config_path"`
	GitLabInstanceURL        string                   `toml:"gitlab_instance_url"`
	GitLabInsecureSkipVerify bool                     `toml:"gitlab_insecure_skip_verify"`
	DatabaseConfiguration    DBConfiguration          `toml:"database"`
	LMSTitle                 string                   `toml:"lms_title,omitempty"`
	LMSURL                   string                   `toml:"lms_url,omitempty"`
	CheckActiveRepoCron      string                   `toml:"check_active_repositories_cron"`
	TimezoneName             string                   `toml:"timezone"`
	Timezone                 *time.Location           `toml:"-"`
	CodeSnippetIncludeLines  int                      `toml:"code_snippet_include_previous_lines"`
	TLSConfiguration         TLSServerConfiguration   `toml:"tls_server_configuration"`
	Limits                   LimitsConfiguration      `toml:"limits"`
	GitLabCustomisation      GitLabIssueCustomisation `toml:"gitlab_issue"`
	Survey                   SurveyConfiguration      `toml:"survey"`
}

// TLSServerConfiguration represents the server TLS
// configuration to enable HTTPS.
type TLSServerConfiguration struct {
	Enabled  bool   `toml:"enabled"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

// LimitsConfiguration represents limits for report items and concurrency.
type LimitsConfiguration struct {
	MaxCheckWorkers          int `toml:"max_check_workers"`
	MaxIssuesPerReport       int `toml:"max_issues_per_report"`
	MaxIssuePerTypePerReport int `toml:"max_issues_per_type_per_report"`
}

// GitLabIssueCustomisation represents the customisation configuration
// for generated GitLab issues.
type GitLabIssueCustomisation struct {
	IssueTitle     string `toml:"title"`
	GeneratingBody string `toml:"generating_body"`
	CompleteBody   string `toml:"complete_body"`
}

// SurveyConfiguration represents the configuration and customisation
// of the surveying function.
type SurveyConfiguration struct {
	ShowSurvey bool   `toml:"show_survey"`
	Title      string `toml:"title"`
	Message    string `toml:"message"`
	SurveyURL  string `toml:"ext_url"`
}

// DBType represents the database driver type, such as MySQL or SQLite.
type DBType int

const (
	// SQLite is a serverless and file-based SQL driver.
	SQLite = iota
	// MySQL is a standard SQL driver.
	MySQL
	PostgreSQL
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
	databaseConfig DBConfiguration, tlsConfig TLSServerConfiguration) Configuration {
	return Configuration{
		SiteTitle:                "Learning Bot",
		SiteURL:                  siteURL,
		BotPrivateToken:          token,
		GitLabInstanceURL:        instance,
		GitLabInsecureSkipVerify: false,
		CheckstyleJarPath:        checkstyleJar,
		CheckstyleConfigPath:     "./assets/checkstyle-lb.xml",
		DatabaseConfiguration:    databaseConfig,
		LMSTitle:                 "Vision",
		LMSURL:                   "https://vision.hw.ac.uk",
		CheckActiveRepoCron:      "@every 1h45m",
		TimezoneName:             "Europe/London",
		CodeSnippetIncludeLines:  3,
		TLSConfiguration:         tlsConfig,
		Limits: LimitsConfiguration{
			MaxCheckWorkers:          5,
			MaxIssuesPerReport:       -1,
			MaxIssuePerTypePerReport: -1,
		},
		GitLabCustomisation: GitLabIssueCustomisation{
			IssueTitle: "[$site_title] Your code report 📊",
			GeneratingBody: `Hey @$username!

Your report is currently being generating for $commit.

Sit tight!  
You can view the progress here:  
[View Report]($report_link)`,
			CompleteBody: `Hey @$username!

The report is generated for $commit, and you can view it in the link below!

[View report]($report_link)`,
		},
		Survey: SurveyConfiguration{
			ShowSurvey: false,
			Title:      "Survey",
			Message:    "We are conducting a study on effectiveness of code repair tools on programming. Please take a minute to fill out our survey.",
			SurveyURL:  "https://example.com/forms/form-id?user=$username",
		},
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
	url := strings.Split(Config.SiteURL, ":")
	if len(url) > 1 {
		Config.SitePort = fmt.Sprintf(":%s", url[2])
	} else {
		Config.SitePort = ":4000"
	}

	LoadActiveProjs(false)
}
