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

type Searcher interface {
	Search(ctx context.Context, root, query, relativePath string) ([]Match, error)
}

type Match struct {
	Path   string `json:"path"`
	Line   int    `json:"line"`
	Column int    `json:"column,omitempty"`
	Text   string `json:"text"`
}

type RipgrepSearcher struct {
	executable string
}

func NewRipgrepSearcher(executable string) *RipgrepSearcher {
	if executable == "" {
		executable = "rg"
	}
	return &RipgrepSearcher{executable: executable}
}

func (s *RipgrepSearcher) Search(ctx context.Context, root, query, relativePath string) ([]Match, error) {
	target := root
	if relativePath != "" && relativePath != "." {
		target = filepath.Join(root, relativePath)
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
