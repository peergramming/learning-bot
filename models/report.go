package models

import (
	"errors"
)

// ReportStatus represents the generation status of a report.
type ReportStatus int

const (
	// InProgress is when the report generation is in-progress.
	InProgress = iota
	// Complete is when the report generation is complete.
	Complete
	// Failed is when the report generation failed.
	Failed
)

// Report represents a report of a specific commit of a project.
type Report struct {
	ReportID            int64        `xorm:"pk autoincr"`
	RepositoryID        string       `xorm:"varchar(64) notnull"`
	Commit              string       `xorm:"index varchar(40) notnull"`
	Status              ReportStatus `xorm:"notnull"`
	RawCheckstyleOutput string       `xorm:"mediumtext"`
	CreatedUnix         int64        `xorm:"created"`
	Issues              []*Issue     `xorm:"-"`
}

// Issue represents a single issue, which is usually a part of
// a Report along with other Issues.
type Issue struct {
	IssueID       int    `xorm:"autoincr pk"`
	ReportID      int64  `xorm:"notnull"`
	CheckName     string `xorm:"varchar(16) notnull"`
	FilePath      string `xorm:"varchar(32) notnull"`
	LineNumber    int    `xorm:"null"`
	ColumnNumber  int    `xorm:"null"`
	Description   string `xorm:"varchar(128) notnull"`
	SourceSnippet string `xorm:"text null"`
}

func getReportsByRepoID(id string) (reports []*Report, err error) {
	return reports, engine.Where("repository_id = ?", id).Find(&reports)
}

// AddReport adds a new report to the database. Returns an error
// if fails.
func AddReport(r *Report) (err error) {
	if r == nil {
		return errors.New("Report is nil")
	}
	_, err = engine.Insert(r)
	return err
}

// UpdateReport updates an existing report in the database. Returns
// an error if fails.
func UpdateReport(r *Report) (err error) {
	if r == nil {
		return errors.New("Report is nil")
	}
	_, err = engine.Update(r)
	return err
}

// HasReport returns whether a specific report is in the database.
func HasReport(repoID string, commit string) bool {
	has, _ := engine.Get(&Report{RepositoryID: repoID, Commit: commit})
	return has
}

// UpdateIssues updates or inserts the issues of a report in the database.
func UpdateIssues(r *Report) (err error) {
	sess := engine.NewSession() // transaction
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Where("report_id = ?", r.ReportID).Delete(new(Issue)); err != nil {
		return err
	}
	for _, issue := range r.Issues {
		issue.ReportID = r.ReportID
	}
	if _, err = sess.Insert(r.Issues); err != nil {
		return err
	}
	return sess.Commit()
}

// LoadIssues loads all issues in a report.
func (r *Report) LoadIssues() (err error) {
	return engine.Where("report_id = ?", r.ReportID).Find(&r.Issues)
}
