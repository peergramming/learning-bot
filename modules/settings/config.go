package settings

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
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

type ActiveProjects struct {
	Projects []Project
}

type Project struct {
	Namespace string
	Project   string
}

func GetGitLabClient() *gitlab.Client {
	if gitlabClient == nil {
		gitlabClient = gitlab.NewClient(nil, Config.BotPrivateToken)
		gitlabClient.SetBaseURL(fmt.Sprintf("%s/api/v4", Config.GitLabInstanceURL))
	}
	return gitlabClient
}

func IsActiveProject(namespace string, project string) (bool, int) {
	for id, proj := range ActiveProjs.Projects {
		if proj.Namespace == namespace && proj.Project == project {
			return true, id
		}
	}
	return false, 0
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

func LoadActiveProjs(quiet bool) {
	var err error
	if _, err = toml.DecodeFile(WorkingDir+"/"+ActiveProjsPath, &ActiveProjs); err != nil && !quiet {
		log.Printf("Cannot load active projects file! Error: %s", err)
		log.Printf("It is safe to ignore this error if you haven't created active projects file yet.")
	}
}

func SaveActiveProjs() {
	var err error

	buf := new(bytes.Buffer)
	if err = toml.NewEncoder(buf).Encode(ActiveProjs); err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(ActiveProjsPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
