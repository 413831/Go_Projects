package services

import (
	"fmt"

	"mcp-server/internal/domain"
	"mcp-server/internal/infrastructure/repositories"
	"mcp-server/pkg/errors"
)

// IssueServiceInterface defines the issues service interface
type IssueServiceInterface interface {
	GetIssues(req *domain.GetIssuesRequest) (*domain.GetIssuesResponse, error)
	FormatIssuesForMCP(issues []domain.Issue) []string
	ValidateGetIssuesRequest(req *domain.GetIssuesRequest) error
}

// IssueService implements business logic for issues
type IssueService struct {
	repo repositories.GitHubRepositoryInterface
}

// NewIssueService creates a new IssueService instance
func NewIssueService(repo repositories.GitHubRepositoryInterface) *IssueService {
	return &IssueService{
		repo: repo,
	}
}

// GetIssues fetches issues with validations and business logic applied
func (s *IssueService) GetIssues(req *domain.GetIssuesRequest) (*domain.GetIssuesResponse, error) {
	// Validate request
	if err := s.ValidateGetIssuesRequest(req); err != nil {
		return nil, err
	}

	// Set default state if not provided
	if req.State == "" {
		req.State = "open"
	}

	// Fetch data from repository
	response, err := s.repo.GetIssues(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// FormatIssuesForMCP formats issues for MCP output
func (s *IssueService) FormatIssuesForMCP(issues []domain.Issue) []string {
	var formatted []string
	
	for _, issue := range issues {
		// Format: #[Number] [State] Title
		formattedIssue := fmt.Sprintf("#%d [%s] %s", issue.Number, issue.State, issue.Title)
		formatted = append(formatted, formattedIssue)
	}

	return formatted
}

// ValidateGetIssuesRequest validates request parameters
func (s *IssueService) ValidateGetIssuesRequest(req *domain.GetIssuesRequest) error {
	if req.Owner == "" {
		return errors.NewValidationError("the 'owner' parameter is required")
	}
	
	if req.Repo == "" {
		return errors.NewValidationError("the 'repo' parameter is required")
	}

	// Validate that state is valid if provided
	if req.State != "" && req.State != "open" && req.State != "closed" && req.State != "all" {
		return errors.NewValidationError("the 'state' parameter must be 'open', 'closed', or 'all'")
	}

	return nil
}
