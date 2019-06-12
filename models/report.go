package models

type Report struct {
	ReportID   int    `xorm:"autoincr pk"`
	User       string `xorm:"varchar(32) notnull"`
	Repository string `xorm:"varchar(32) notnull"`
	Commit     string `xorm:"varchar(40) notnull"`
	Issues     []Issue
}

type Issue struct {
	IssueID       int    `xorm:"autoincr pk"`
	CheckName     string `xorm:"varchar(16) notnull"`
	FilePath      string `xorm:"varchar(32) notnull"`
	LineNumber    int    `xorm:"notnull"`
	Description   string `xorm:"varchar(128) notnull"`
	SourceSnippit string `xorm:"text null"`
}
