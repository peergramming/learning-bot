package models

import (
	"github.com/go-xorm/xorm"
)

// Report represents a report of a specific commit of a project.
type Report struct {
	ReportID            int64    `xorm:"pk autoincr"`
	RepositoryID        string   `xorm:"varchar(64) notnull"`
	Commit              string   `xorm:"index varchar(40) notnull"`
	RawCheckstyleOutput string   `xorm:"mediumtext"`
	CreatedUnix         int64    `xorm:"created"`
	Issues              []*Issue `xorm:"-"`
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

func UpdateRepositoryReports(repo *Repository, reports []*Report) (err error) {
	sess := engine.NewSession() // transaction
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Where("repository_id = ?", repo.RepoID).Delete(new(Report)); err != nil {
		return err
	}

	if _, err = sess.Insert(reports); err != nil {
		return err
	}

	for _, report := range reports {
		err = report.updateIssues(sess)
		if err != nil {
			return err
		}
	}

	return sess.Commit()
}

func (r *Report) updateIssues(sess *xorm.Session) (err error) {
	if _, err = sess.Where("report_id = ?", r.ReportID).Delete(new(Issue)); err != nil {
		return err
	}
	for _, issue := range r.Issues {
		issue.ReportID = r.ReportID
	}
	if _, err = sess.Insert(r.Issues); err != nil {
		return err
	}
	return nil
}

func (r *Report) LoadIssues() (err error) {
	return engine.Where("report_id = ?", r.ReportID).Find(&r.Issues)
}
