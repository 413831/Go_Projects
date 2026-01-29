package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"mcp-server/internal/domain"
	"mcp-server/pkg/errors"
)

// GitHubClient is an HTTP client for the GitHub API
type GitHubClient struct {
	client  *http.Client
	baseURL string
	token   string
}

// NewGitHubClient creates a new GitHubClient instance
func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		client:  http.DefaultClient,
		baseURL: "https://api.github.com",
		token:   os.Getenv("GITHUB_TOKEN"),
	}
}

// GetIssues fetches issues from a repository
func (c *GitHubClient) GetIssues(req *domain.GetIssuesRequest) (*domain.GetIssuesResponse, error) {
	if c.token == "" {
		return nil, errors.NewUnauthorizedError()
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues", c.baseURL, req.Owner, req.Repo)
	if req.State != "" {
		url += fmt.Sprintf("?state=%s", req.State)
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewNetworkError(fmt.Sprintf("creating request: %v", err))
	}

	httpReq.Header.Set("Authorization", "token "+c.token)
	httpReq.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, errors.NewNetworkError(fmt.Sprintf("executing request: %v", err))
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return c.parseIssuesResponse(resp.Body)
	case http.StatusUnauthorized:
		return nil, errors.NewUnauthorizedError()
	case http.StatusNotFound:
		return nil, errors.NewNotFoundError(fmt.Sprintf("repository %s/%s", req.Owner, req.Repo))
	default:
		return nil, c.handleAPIError(resp)
	}
}

// parseIssuesResponse parses the issues response
func (c *GitHubClient) parseIssuesResponse(body io.ReadCloser) (*domain.GetIssuesResponse, error) {
	var issues []domain.Issue
	if err := json.NewDecoder(body).Decode(&issues); err != nil {
		return nil, errors.NewJSONDecodingError(fmt.Sprintf("decoding issues: %v", err))
	}

	return &domain.GetIssuesResponse{
		Issues: issues,
		Count:  len(issues),
	}, nil
}

// handleAPIError handles GitHub API errors
func (c *GitHubClient) handleAPIError(resp *http.Response) error {
	var apiErr struct {
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"errors,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return errors.NewGitHubAPIError(fmt.Sprintf("error %d: unknown error", resp.StatusCode))
	}

	details := apiErr.Message
	if len(apiErr.Errors) > 0 {
		details += fmt.Sprintf(" (%v)", apiErr.Errors)
	}

	return errors.NewGitHubAPIError(fmt.Sprintf("error %d: %s", resp.StatusCode, details))
}
