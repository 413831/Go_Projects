package main

import (
	"log"

	"mcp-server/internal/config"
	"mcp-server/internal/interfaces"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create dependency container
	container := interfaces.NewContainer()

	// Configure MCP server
	mcpServer := container.SetupMCPServer(cfg.ServerName, cfg.ServerVersion)

	// Start server
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Error starting MCP server: %v", err)
	}
}
