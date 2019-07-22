package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddGetRepo(t *testing.T) {
	assert := assert.New(t)
	SetupTestEngine()
	repo := &Repository{RepoID: "some/repo", IssueID: 5}
	AddRepo(repo)
	report := &Report{
		RepositoryID:        repo.RepoID,
		Commit:              "cd86dcf4737a47ae6a9909e43a9a01d7aee71fa0",
		Status:              Complete,
		RawCheckstyleOutput: "Some output",
	}
	if HasReport(repo.RepoID, report.Commit) {
		UpdateReport(report)
	} else {
		AddReport(report)
	}
	// Getting
	repo2, err := GetRepo("some/repo")
	assert.Nil(err)
	assert.Equal(repo2.RepoID, "some/repo", "The repo ID should be correct")
	assert.Equal(len(repo2.Reports), 1, "The number of reports must be correct")
	assert.Equal(repo2.Reports[0].RepositoryID, repo2.RepoID, "The report's repository ID must match the repository ID")
	assert.Equal(repo2.Reports[0].Status, ReportStatus(Complete))
}

func TestNotHasRepo(t *testing.T) {
	SetupTestEngine()
	assert.False(t, HasReport("some/nonexistent", "aaabbbccc"), "Report should not exist")
}
