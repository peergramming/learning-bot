package models

type Repository struct {
	RepoID   int
	Owner    string
	RepoName string
	IssueID  int // stored id of the main issue for checkstyle report
}
