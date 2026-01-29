package tools

import (
	"context"
	"fmt"

	"mcp-server/internal/application/services"
	"mcp-server/internal/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

// ToolFactory creates MCP tools using the Factory pattern
type ToolFactory struct {
	issueService services.IssueServiceInterface
}

// NewToolFactory creates a new ToolFactory instance
func NewToolFactory(issueService services.IssueServiceInterface) *ToolFactory {
	return &ToolFactory{
		issueService: issueService,
	}
}

// CreateGetIssuesTool creates the tool for fetching GitHub issues
func (f *ToolFactory) CreateGetIssuesTool() mcp.Tool {
	return mcp.NewTool("get_issues",
		mcp.WithDescription("Fetches issues from a GitHub repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner (organization or user)")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("state", mcp.Description("Issue state: open, closed, all (default: open)")),
	)
}

// CreateGetIssuesHandler creates the handler for the get_issues tool
func (f *ToolFactory) CreateGetIssuesHandler() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract and validate arguments
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Unable to read request arguments"), nil
		}

		// Build request
		request := &domain.GetIssuesRequest{
			Owner: getStringArg(args, "owner"),
			Repo:  getStringArg(args, "repo"),
			State: getStringArg(args, "state"),
		}

		// Execute business logic
		response, err := f.issueService.GetIssues(request)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error fetching issues", err), nil
		}

		// Format response for MCP
		formattedIssues := f.issueService.FormatIssuesForMCP(response.Issues)
		
		var contents []mcp.Content
		for _, issue := range formattedIssues {
			contents = append(contents, mcp.NewTextContent(issue))
		}

		// Add summary information
		summary := mcp.NewTextContent(fmt.Sprintf("\nFound %d issues in %s/%s", response.Count, request.Owner, request.Repo))
		contents = append(contents, summary)

		return &mcp.CallToolResult{Content: contents}, nil
	}
}

// getStringArg retrieves a string argument from the arguments map
func getStringArg(args map[string]interface{}, key string) string {
	if value, ok := args[key].(string); ok {
		return value
	}
	return ""
}
