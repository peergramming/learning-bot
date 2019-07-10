package settings

import (
	"crypto/tls"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
)

var gitlabClient *gitlab.Client

// GetGitLabClient returns the GitLab client, with the private token and
// instance (base) URL pre-set.
func GetGitLabClient() *gitlab.Client {
	if gitlabClient == nil {
		var client *http.Client
		if Config.GitLabInsecureSkipVerify {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}
		gitlabClient = gitlab.NewClient(client, Config.BotPrivateToken)
		gitlabClient.SetBaseURL(fmt.Sprintf("%s/api/v4", Config.GitLabInstanceURL))
	}
	return gitlabClient
}
