package search

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Searcher abstracts text search so the MCP layer can use either the built-in
// backend or the optional ripgrep process.
type Searcher interface {
	Search(ctx context.Context, root, query, relativePath string) ([]Match, error)
}

// Match is a single text match returned to an MCP client.
type Match struct {
	Path   string `json:"path"`
	Line   int    `json:"line"`
	Column int    `json:"column,omitempty"`
	Text   string `json:"text"`
}

// RipgrepSearcher invokes ripgrep with machine-readable JSON output.
type RipgrepSearcher struct {
	executable string
}

// NewRipgrepSearcher creates a searcher using executable when provided, or
// "rg" so the binary can rely on the user's PATH by default.
func NewRipgrepSearcher(executable string) *RipgrepSearcher {
	if executable == "" {
		executable = "rg"
	}
	return &RipgrepSearcher{executable: executable}
}

// Search searches Markdown, MDX, and plain-text files under relativePath.
// Ripgrep's exit code 1 means "no matches" and is therefore treated as success.
func (s *RipgrepSearcher) Search(ctx context.Context, root, query, relativePath string) ([]Match, error) {
	target, err := resolveTarget(root, relativePath)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(
		ctx,
		s.executable,
		"--json",
		"--line-number",
		"--column",
		"--color", "never",
		"--no-heading",
		"--glob", "*.md",
		"--glob", "*.mdx",
		"--glob", "*.txt",
		query,
		target,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			message := strings.TrimSpace(stderr.String())
			if message == "" {
				message = err.Error()
			}
			return nil, fmt.Errorf("ripgrep search failed: %s", message)
		}
	}

	return parseResults(root, &stdout)
}

// resolveTarget keeps the search scope inside root before spawning ripgrep.
// This mirrors the repository path checks used by list and read operations.
func resolveTarget(root, relativePath string) (string, error) {
	if relativePath == "" || relativePath == "." {
		return root, nil
	}
	if filepath.IsAbs(relativePath) {
		return "", errors.New("search path must be repository-relative")
	}

	target := filepath.Clean(filepath.Join(root, relativePath))
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("search path escapes the knowledge root")
	}
	return target, nil
}

type rgMessage struct {
	Type string `json:"type"`
	Data struct {
		Path struct {
			Text string `json:"text"`
		} `json:"path"`
		Lines struct {
			Text string `json:"text"`
		} `json:"lines"`
		LineNumber int `json:"line_number"`
		Submatches []struct {
			Start int `json:"start"`
		} `json:"submatches"`
	} `json:"data"`
}

// parseResults converts ripgrep's JSON event stream into stable repository-
// relative matches and ignores non-match events such as summaries.
func parseResults(root string, output *bytes.Buffer) ([]Match, error) {
	results := make([]Match, 0)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var message rgMessage
		if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
			return nil, fmt.Errorf("parse ripgrep output: %w", err)
		}
		if message.Type != "match" {
			continue
		}

		path, err := filepath.Rel(root, message.Data.Path.Text)
		if err != nil {
			path = message.Data.Path.Text
		}
		match := Match{
			Path: filepath.ToSlash(path),
			Line: message.Data.LineNumber,
			Text: strings.TrimRight(message.Data.Lines.Text, "\r\n"),
		}
		if len(message.Data.Submatches) > 0 {
			match.Column = message.Data.Submatches[0].Start + 1
		}
		results = append(results, match)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read ripgrep output: %w", err)
	}
	return results, nil
}
