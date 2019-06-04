package models

import (
	"time"
)

const (
	Opened = "opened"
	Closed = "closed"
)

type GitLabIssue struct {
	ID          int          `json:"id"`
	State       string       `json:"state"`
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
