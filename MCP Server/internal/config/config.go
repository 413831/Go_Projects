package config

import (
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	ServerName    string
	ServerVersion string
	GitHubToken   string
	LogLevel      string
}

// NewConfig creates a new configuration
func NewConfig() *Config {
	return &Config{
		ServerName:    getEnvOrDefault("MCP_SERVER_NAME", "GitHubIssues"),
		ServerVersion: getEnvOrDefault("MCP_SERVER_VERSION", "0.0.1"),
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	return nil
}

// getEnvOrDefault returns an environment variable or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
