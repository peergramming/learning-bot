package models

import (
	"errors"
)

// Repository represents a specific GitLab project.
// It keeps track of the issue tracker post and all
// the reports of that specific repository.
type Repository struct {
	RepoID  string    `xorm:"varchar(64) pk"`
	IssueID int       `xorm:"null"` // stored id of the main issue for checkstyle report
	Reports []*Report `xorm:"-"`
}

// GetRepo returns the repository from the project ID.
// It returns the repository and error (if exists).
func GetRepo(id string) (*Repository, error) {
	r := new(Repository)
	has, err := engine.ID(id).Get(r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Repository does not exist")
	}
	err = r.getReports()
	if err != nil {
		return r, err
	}
	return r, nil
}

// AddRepo adds a new Repository to the database.
// It returns the columns affected and error (if exists).
func AddRepo(r *Repository) (err error) {
	_, err = engine.Insert(r)
	return err
}

// UpdateRepo updates a Repository in the database.
// It returns the columns affected and error (if exists).
func UpdateRepo(r *Repository) (err error) {
	_, err = engine.Update(r)
	return err
}

func (r *Repository) getReports() (err error) {
	if r.Reports != nil {
		return nil
	}

	r.Reports, err = getReportsByRepoID(r.RepoID)
	return err
}

// GetReport returns the report and whether it exists or not
// based on the commit SHA.
func (r *Repository) GetReport(sha string) (*Report, bool) {
	for _, rep := range r.Reports {
		if rep.Commit == sha {
			return rep, true
		}
	}
	return nil, false
}
