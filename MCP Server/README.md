# MCP Server - GitHub Issues

An MCP (Model Context Protocol) server that provides tools to fetch issues from GitHub repositories, implemented with clean architecture and SOLID design patterns.

## 🏗️ Architecture

### Directory Structure

```
mcp-server/
├── cmd/
│   └── main.go                 # Minimal entry point
├── internal/
│   ├── config/                 # Application configuration
│   │   └── config.go
│   ├── domain/                 # Domain entities
│   │   └── models.go
│   ├── infrastructure/         # External layer (HTTP, repositories)
│   │   ├── http/
│   │   │   └── github_client.go
│   │   └── repositories/
│   │       └── github_repository.go
│   ├── application/            # Business logic
│   │   ├── services/
│   │   │   └── issue_service.go
│   │   └── tools/
│   │       └── tool_factory.go
│   └── interfaces/             # Interfaces and DI container
│       └── mcp_handlers.go
├── pkg/                        # Reusable code
│   └── errors/
│       └── errors.go
├── go.mod
├── go.sum
└── README.md
```

### Implemented Design Patterns

1. **Clean Architecture**: Clear separation of responsibilities
2. **Repository Pattern**: Data access abstraction
3. **Service Layer**: Centralized business logic
4. **Factory Pattern**: MCP tool creation
5. **Dependency Injection**: Dependency management
6. **Error Handling**: Typed and centralized errors

## 🚀 Usage

### Environment Variables

```bash
# GitHub personal access token (required)
export GITHUB_TOKEN="your_github_token"

# Optional server configuration
export MCP_SERVER_NAME="GitHubIssues"
export MCP_SERVER_VERSION="0.0.1"
export LOG_LEVEL="info"
```

### Run

```bash
# Build
go build ./cmd

# Run
./cmd.exe
```

## 📋 Available Tools

### get_issues
Fetches issues from a GitHub repository.

**Parameters:**
- `owner` (required): Repository owner
- `repo` (required): Repository name  
- `state` (optional): Issue state (`open`, `closed`, `all`)

**Example:**
```json
{
  "name": "get_issues",
  "arguments": {
    "owner": "microsoft",
    "repo": "vscode",
    "state": "open"
  }
}
```

## 🔧 Detailed Architecture

### Domain Layer (`internal/domain`)
- **Models**: Domain entities such as `Issue`, `User`, `Label`
- **Request/Response**: Structures for cross-layer communication

### Infrastructure Layer (`internal/infrastructure`)
- **HTTP Client**: GitHub API client with error handling
- **Repository**: Repository pattern implementation for data access

### Application Layer (`internal/application`)
- **Services**: Business logic, validations, and transformations
- **Tools**: Factory for creating MCP tools and handlers

### Interfaces Layer (`internal/interfaces`)
- **Container**: Dependency injection and configuration
- **MCP Handlers**: MCP server setup

### Error Handling (`pkg/errors`)
- **AppError**: Typed error structure
- **Error codes**: Predefined errors for common cases

## 🧪 Testing

The architecture makes unit testing straightforward:

```go
// Example test for IssueService
func TestIssueService_GetIssues(t *testing.T) {
    // Repository mock
    mockRepo := &MockGitHubRepository{}
    service := NewIssueService(mockRepo)
    
    // Test cases...
}
```

## 🔄 Extensibility

To add new tools:

1. Define models in `domain/models.go`
2. Implement the HTTP client in `infrastructure/http/`
3. Create a repository in `infrastructure/repositories/`
4. Add a service in `application/services/`
5. Create the tool in `application/tools/tool_factory.go`
6. Register it in `interfaces/mcp_handlers.go`

## 📝 SOLID Principles

- **S**: Each class has a single responsibility
- **O**: Open for extension, closed for modification
- **L**: Interfaces depend on abstractions
- **I**: Interfaces are specific to clients
- **D**: Dependencies are injected, not hardcoded

## 🤝 Contribution

1. Follow the established architecture
2. Add tests for new functionality
3. Document changes in the README
4. Maintain compatibility with the MCP specification
