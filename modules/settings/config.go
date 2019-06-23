package settings

import (
	"github.com/BurntSushi/toml"
	"github.com/xanzy/go-gitlab"
	"log"
	"fmt"
	"os"
)

var (
	ConfigPath string
	Config     Configuration
)

type Configuration struct {
	SiteTitle             string
	BotPrivateToken       string
	CheckstyleJarPath     string
	CheckstyleConfigPath  string
	GitLabInstanceURL     string
	DatabaseConfiguration DBConfiguration
	LMSTitle              string
	LMSURL                string
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

var gitlabClient *gitlab.Client

func GetGitLabClient() *gitlab.Client {
	if gitlabClient == nil {
		gitlabClient = gitlab.NewClient(nil, Config.BotPrivateToken)
		gitlabClient.SetBaseURL(fmt.Sprintf("%s/api/v4", Config.GitLabInstanceURL))
	}
	return gitlabClient
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
