package settings

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

var gitlabClient *gitlab.Client

// GetGitLabClient returns the GitLab client, with the private token and
// instance (base) URL pre-set.
func GetGitLabClient() *gitlab.Client {
	if gitlabClient == nil {
		gitlabClient = gitlab.NewClient(nil, Config.BotPrivateToken)
		gitlabClient.SetBaseURL(fmt.Sprintf("%s/api/v4", Config.GitLabInstanceURL))
	}
	return gitlabClient
}
