package settings

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type ActiveProjects struct {
	Projects []Project
}

type Project struct {
	Namespace string
	Project   string
}

func IsActiveProject(namespace string, project string) (bool, int) {
	for id, proj := range ActiveProjs.Projects {
		if proj.Namespace == namespace && proj.Project == project {
			return true, id
		}
	}
	return false, 0
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
