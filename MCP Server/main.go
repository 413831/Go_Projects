package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Issue struct {
	Number  int    `json:"numbre"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
}

func main() {
	s := server.NewMCPServer(
		"GitHubIssues",
		"0.0.1",
		server.WithLogging())

	tool := mcp.NewTool("get_issues",
		mcp.WithDescription("Obtiene las issues abiertas de un repositorio"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Propietario del repo (org o user)")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Nombre del repositorio")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var (
			owner string = ""
			repo  string = ""
		)

		args, ok := req.Params.Arguments.(map[string]interface)
		if !ok {
			return nil, errors.New("Cannot read request arguments")
		}

		if o, ok := args["owner"].(string); ok {
			owner = o
		}
		if r, ok := args["repo"].(string); ok {
			repo = r
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

		resp, err := http.DefaultClient.Do(request)
		if err != nil{
			return mcp.NewToolResultErrorFromErr("Falló petición a Github", err), nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK{
			errResult := jsonDecodeError(resp)

			return errResult, nil // Refactor
		}

		var issues []Issue
		if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
			return mcp.NewToolResultErrorFromErr("Error decodificando JSON de issues", err), nil
		}


	})
}

func jsonDecodeError(resp *http.Response) (*mcp.CallToolResult) {
	var apiErr struct {
		Message string `json:message`
	}

	if decErr := json.NewDecoder(resp.Body).Decode(&apiErr); decErr != nil {
		return mcp.NewToolResultError("Error leyendo JSON de error")
	}
	return mcp.NewToolResultError(apiErr.Message)
}
