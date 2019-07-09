package settings

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

var (
	// ActiveProjsPath is the path of the active projects list.
	ActiveProjsPath = "active-projects.toml"
	// ActiveProjs holds the list of all active projects.
	ActiveProjs ActiveProjects
)

// ActiveProjects represents the active projects list file,
// containing a list of all active GitLab projects, which are
// checked at an interval.
type ActiveProjects struct {
	Projects []Project `toml:"projects,omitempty"`
}

// Project represents a GitLab project URL.
type Project struct {
	Namespace string `toml:"namespace"`
	Project   string `toml:"project"`
}

// GetFullPath returns the full path (namespace / project) of
// the project.
func (p *Project) GetFullPath() string {
	return fmt.Sprintf("%s/%s", p.Namespace, p.Project)
}

// IsActiveProject returns whether a project exists in the active projects
// list.
// It returns whether it exists as a boolean, and the element number in the
// array of ActiveProjects.Projects.
func IsActiveProject(namespace string, project string) (bool, int) {
	for id, proj := range ActiveProjs.Projects {
		if proj.Namespace == namespace && proj.Project == project {
			return true, id
		}
	}
	return false, 0
}

// LoadActiveProjs loads the active projects configuration from file.
// quiet determines whether to fail quitely.
func LoadActiveProjs(quiet bool) {
	var err error
	if _, err = toml.DecodeFile(WorkingDir+"/"+ActiveProjsPath, &ActiveProjs); err != nil && !quiet {
		log.Printf("Cannot load active projects file! Error: %s", err)
		log.Printf("It is safe to ignore this error if you haven't created active projects file yet.")
	}
}

// SaveActiveProjs saves the list of active projects back to its
// configuration file.
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
