package checkstyle

import (
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
)

func GenerateReport(repo string, sha string) {
	// Check if sha is latest, if not, skip
	// commits/ListCommits()
	git := settings.GetGitLabClient()
	git.Commits.ListCommits(repo, &gitlab.ListCommitsOptions{})

	// Download: repositories/archive()
	// Options: Format (bz2) + SHA
}
