package models

import (
	"errors"
)

type Repository struct {
	RepoID             int    `xorm:"autoincr pk"`
	Owner              string `xorm:"varchar(24) notnull"`
	RepoName           string `xorm:"varchar(24) notnull"`
	IssueID            int    `xorm:"null"` // stored id of the main issue for checkstyle report
	PrivateAccessToken string `xorm:"varchar(64) null"`
}

// GetRepo returns the repository from owner/repo
func GetRepo(owner string, repo string) (Repository, error) {
	r := Repository{Owner: owner, RepoName: repo}
	has, err := engine.Get(&r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Repository does not exist")
	}
	return r, nil
}
