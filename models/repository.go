package models

import (
	"errors"
)

type Repository struct {
	ProjectID          string `xorm:"varchar(64) pk"`
	IssueID            int    `xorm:"null"` // stored id of the main issue for checkstyle report
}

// GetRepo returns the repository from owner/repo
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
