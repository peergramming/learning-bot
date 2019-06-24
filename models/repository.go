package models

import (
	"errors"
)

// Repository represents a specific GitLab project.
// It keeps track of the issue tracker post and all
// the reports of that specific repository.
type Repository struct {
	ProjectID string `xorm:"varchar(64) pk"`
	IssueID   int    `xorm:"null"` // stored id of the main issue for checkstyle report
	Reports   []Report
}

// GetRepo returns the repository from the project ID.
// It returns the repository and error (if exists).
func GetRepo(id string) (Repository, error) {
	r := Repository{ProjectID: id}
	has, err := engine.Get(&r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Repository does not exist")
	}
	return r, nil
}
