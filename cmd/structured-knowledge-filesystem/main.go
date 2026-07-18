package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ImShaoJun/structured-knowledge-filesystem/internal/config"
	"github.com/ImShaoJun/structured-knowledge-filesystem/internal/knowledge"
	"github.com/ImShaoJun/structured-knowledge-filesystem/internal/search"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// main loads the repository configuration and exposes the read-only MCP tools
// over stdio so that desktop clients can start this process on demand.
func main() {
	log.SetOutput(os.Stderr)

	configPath := flag.String("config", "", "path to a JSON configuration file")
	rootFlag := flag.String("root", "", "knowledge repository root directory")
	flag.Parse()

	cfg, err := config.Load(*configPath, *rootFlag)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := knowledge.NewRepository(cfg.Root)
	if err != nil {
		log.Fatal(err)
	}

	searcher := search.NewRipgrepSearcher(cfg.RipgrepPath)
	s := server.NewMCPServer(
		"Structured Knowledge Filesystem",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(listDirectoryTool(), listDirectoryHandler(repo))
	s.AddTool(readFileTool(), readFileHandler(repo))
	s.AddTool(searchTool(), searchHandler(repo, searcher))

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

// listDirectoryTool describes the directory navigation tool exposed to MCP
// clients. The repository implementation enforces the configured root.
func listDirectoryTool() mcp.Tool {
	return mcp.NewTool(
		"list_directory",
		mcp.WithDescription("List files and directories in the knowledge repository. Use this before reading an unfamiliar path."),
		mcp.WithString("path", mcp.Description("Repository-relative directory path. Use . for the root directory.")),
	)
}

// readFileTool describes the read-only file access tool.
func readFileTool() mcp.Tool {
	return mcp.NewTool(
		"read_file",
		mcp.WithDescription("Read a text file from the knowledge repository."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Repository-relative file path.")),
	)
}

// searchTool describes the ripgrep-backed text search tool.
func searchTool() mcp.Tool {
	return mcp.NewTool(
		"search",
		mcp.WithDescription("Search the knowledge repository with ripgrep and return matching file paths, line numbers, and snippets."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Text or regular expression to search for.")),
		mcp.WithString("path", mcp.Description("Repository-relative directory or file path. Use . when omitted.")),
	)
}

// listDirectoryHandler maps an MCP request to a repository directory listing.
func listDirectoryHandler(repo *knowledge.Repository) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := stringArgument(request, "path")
		entries, err := repo.List(ctx, path)
		if err != nil {
			return toolError(err), nil
		}
		return jsonResult(entries)
	}
}

// readFileHandler reads one repository-relative text file.
func readFileHandler(repo *knowledge.Repository) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := stringArgument(request, "path")
		if path == "" {
			return toolError(errors.New("path is required")), nil
		}
		content, err := repo.Read(ctx, path)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(content), nil
	}
}

// searchHandler searches only within the configured knowledge repository.
func searchHandler(repo *knowledge.Repository, searcher search.Searcher) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := stringArgument(request, "query")
		if strings.TrimSpace(query) == "" {
			return toolError(errors.New("query is required")), nil
		}

		path := stringArgument(request, "path")
		results, err := searcher.Search(ctx, repo.Root(), query, path)
		if err != nil {
			return toolError(err), nil
		}
		return jsonResult(results)
	}
}

// stringArgument returns a trimmed string argument without allowing a malformed
// MCP value to panic the handler.
func stringArgument(request mcp.CallToolRequest, name string) string {
	value, _ := request.GetArguments()[name].(string)
	return strings.TrimSpace(value)
}

// jsonResult serializes structured tool output in a stable, human-readable form.
func jsonResult(value any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(data)), nil
}

// toolError returns a tool-level error so the MCP client can present the
// validation or filesystem failure without terminating the server.
func toolError(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("%v", err))
}
