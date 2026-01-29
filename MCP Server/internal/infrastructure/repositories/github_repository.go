package repositories

import (
	"mcp-server/internal/domain"
	"mcp-server/internal/infrastructure/http"
)

// GitHubRepositoryInterface defines the repository interface
type GitHubRepositoryInterface interface {
	GetIssues(req *domain.GetIssuesRequest) (*domain.GetIssuesResponse, error)
}

// GitHubRepository implements the Repository pattern for GitHub
type GitHubRepository struct {
	client *http.GitHubClient
}

// NewGitHubRepository creates a new GitHubRepository instance
func NewGitHubRepository(client *http.GitHubClient) *GitHubRepository {
	return &GitHubRepository{
		client: client,
	}
}

// GetIssues fetches issues using the HTTP client
func (r *GitHubRepository) GetIssues(req *domain.GetIssuesRequest) (*domain.GetIssuesResponse, error) {
	return r.client.GetIssues(req)
}
