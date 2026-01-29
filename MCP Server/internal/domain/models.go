package domain

import "time"

// Issue represents a GitHub issue
type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
	Labels    []Label   `json:"labels"`
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

// Label represents a GitHub label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// GetIssuesRequest defines parameters for fetching issues
type GetIssuesRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	State string `json:"state,omitempty"` // open, closed, all
}

// GetIssuesResponse contains the issues response
type GetIssuesResponse struct {
	Issues []Issue `json:"issues"`
	Count  int     `json:"count"`
}
