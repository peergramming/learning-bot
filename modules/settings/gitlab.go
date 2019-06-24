package settings

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

var gitlabClient *gitlab.Client

func GetGitLabClient() *gitlab.Client {
	if gitlabClient == nil {
		gitlabClient = gitlab.NewClient(nil, Config.BotPrivateToken)
		gitlabClient.SetBaseURL(fmt.Sprintf("%s/api/v4", Config.GitLabInstanceURL))
	}
	return gitlabClient
}
