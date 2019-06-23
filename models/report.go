package models

import (
	"errors"
)

type Report struct {
	ProjectID  string `xorm:"varchar(64) pk"`
	Commit     string `xorm:"varchar(40) notnull"`
	Issues     []Issue
}

type Issue struct {
	IssueID       int    `xorm:"autoincr pk"`
	CheckName     string `xorm:"varchar(16) notnull"`
	FilePath      string `xorm:"varchar(32) notnull"`
	LineNumber    int    `xorm:"notnull"`
	Description   string `xorm:"varchar(128) notnull"`
	SourceSnippet string `xorm:"text null"`
}

func GetReport(id string, commit string) (Report, error) {
	r := Report{ProjectID: id, Commit: commit}

	has, err := engine.Get(&r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Report does not exist")
	}
	return r, nil
}
