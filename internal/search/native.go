package search

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NativeSearcher searches supported text files in-process using Go's standard
// library. It is the default backend and does not require external binaries.
type NativeSearcher struct{}

// NewNativeSearcher creates the built-in search backend.
func NewNativeSearcher() *NativeSearcher {
	return &NativeSearcher{}
}

// NewSearcher selects the built-in backend by default. Setting ripgrepPath
// explicitly opts into the ripgrep backend for larger repositories.
func NewSearcher(ripgrepPath string) Searcher {
	if strings.TrimSpace(ripgrepPath) != "" {
		return NewRipgrepSearcher(ripgrepPath)
	}
	return NewNativeSearcher()
}

// Search searches Markdown, MDX, and plain-text files under relativePath.
// Query syntax follows Go's regular-expression implementation.
func (s *NativeSearcher) Search(ctx context.Context, root, query, relativePath string) ([]Match, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	target, err := resolveTarget(root, relativePath)
	if err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("compile search query: %w", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("stat search path: %w", err)
	}

	if !info.IsDir() {
		if !isSearchableFile(target) {
			return []Match{}, nil
		}
		return searchFile(ctx, root, target, pattern)
	}

	results := make([]Match, 0)
	err = filepath.WalkDir(target, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || !entry.Type().IsRegular() || !isSearchableFile(path) {
			return nil
		}

		matches, err := searchFile(ctx, root, path, pattern)
		if err != nil {
			return err
		}
		results = append(results, matches...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk search path: %w", err)
	}
	return results, nil
}

func searchFile(ctx context.Context, root, path string, pattern *regexp.Regexp) ([]Match, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer file.Close()

	relative, err := filepath.Rel(root, path)
	if err != nil {
		return nil, fmt.Errorf("make search path relative: %w", err)
	}

	results := make([]Match, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		lineNumber++
		line := scanner.Text()
		location := pattern.FindStringIndex(line)
		if location == nil {
			continue
		}
		results = append(results, Match{
			Path:   filepath.ToSlash(relative),
			Line:   lineNumber,
			Column: location[0] + 1,
			Text:   line,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	return results, nil
}

func isSearchableFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".mdx", ".txt":
		return true
	default:
		return false
	}
}
