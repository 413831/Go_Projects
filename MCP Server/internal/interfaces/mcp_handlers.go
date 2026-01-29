package interfaces

import (
	"mcp-server/internal/application/services"
	"mcp-server/internal/application/tools"
	"mcp-server/internal/config"
	"mcp-server/internal/infrastructure/http"
	"mcp-server/internal/infrastructure/repositories"

	"github.com/mark3labs/mcp-go/server"
)

// Container holds all application dependencies
type Container struct {
	Config        *config.Config
	IssueService  services.IssueServiceInterface
	ToolFactory   *tools.ToolFactory
	GitHubRepo    repositories.GitHubRepositoryInterface
	GitHubClient  *http.GitHubClient
}

// NewContainer creates a new dependency container
func NewContainer() *Container {
	// Create HTTP client
	githubClient := http.NewGitHubClient()
	
	// Create repository
	githubRepo := repositories.NewGitHubRepository(githubClient)
	
	// Create service
	issueService := services.NewIssueService(githubRepo)
	
	// Create tool factory
	toolFactory := tools.NewToolFactory(issueService)

	return &Container{
		GitHubClient: githubClient,
		GitHubRepo:   githubRepo,
		IssueService: issueService,
		ToolFactory:  toolFactory,
	}
}

// SetupMCPServer configures the MCP server with all tools
func (c *Container) SetupMCPServer(name, version string) *server.MCPServer {
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithLogging(),
	)

	// Create and register get_issues tool
	getIssuesTool := c.ToolFactory.CreateGetIssuesTool()
	getIssuesHandler := c.ToolFactory.CreateGetIssuesHandler()
	
	mcpServer.AddTool(getIssuesTool, getIssuesHandler)

	return mcpServer
}
