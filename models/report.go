package models

// Report represents a report of a specific commit of a project.
type Report struct {
	ReportID            int64  `xorm:"pk autoincr"`
	Commit              string `xorm:"index varchar(40) notnull"`
	RawCheckstyleOutput string `xorm:"mediumtext"`
	CreatedUnix         int64  `xorm:"created"`
	Issues              []Issue
}

// Issue represents a single issue, which is usually a part of
// a Report along with other Issues.
type Issue struct {
	IssueID       int    `xorm:"autoincr pk"`
	CheckName     string `xorm:"varchar(16) notnull"`
	FilePath      string `xorm:"varchar(32) notnull"`
	LineNumber    int    `xorm:"null"`
	ColumnNumber  int    `xorm:"null"`
	Description   string `xorm:"varchar(128) notnull"`
	SourceSnippet string `xorm:"text null"`
}
