package gitlab

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

const (
	Opened = "opened"
	Closed = "closed"
)

type IssueState string

type GitLabIssue struct {
	ID          int          `json:"id"`
	State       IssueState   `json:"state"`
	Description string       `json:"description"`
	Author      GitLabAuthor `json:"author"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Title       string       `json:"title"`
	WebURL      string       `json:"web_url"`
}

type GitLabAuthor struct {
	ID       int    `json:"id"`
	WebURL   string `json:"web_url"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func GetRepoIssues(project string) []GitLabIssue {
	url := fmt.Sprintf("/projects/%s/issues", url.PathEscape(project))
	req := GetNewGitLabRequest(url)
	body := DoRequestBytes(req)

	var issues []GitLabIssue
	err := json.Unmarshal(body, &issues)
	if err != nil {
		log.Fatal(err)
	}

	return issues
}
