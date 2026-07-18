package knowledge

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Repository struct {
	root string
}

type Entry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size,omitempty"`
}

func NewRepository(root string) (*Repository, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat knowledge root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("knowledge root is not a directory: %s", root)
	}
	return &Repository{root: filepath.Clean(root)}, nil
}

func (r *Repository) Root() string {
	return r.root
}

func (r *Repository) List(ctx context.Context, relativePath string) ([]Entry, error) {
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}
	directory, err := r.resolve(relativePath)
	if err != nil {
		return nil, err
	}
	items, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("list %q: %w", relativePath, err)
	}

	entries := make([]Entry, 0, len(items))
	for _, item := range items {
		if err := ctxErr(ctx); err != nil {
			return nil, err
		}
		path := filepath.Join(directory, item.Name())
		relative, err := filepath.Rel(r.root, path)
		if err != nil {
			return nil, err
		}
		entry := Entry{Name: item.Name(), Path: filepath.ToSlash(relative)}
		if item.IsDir() {
			entry.Type = "directory"
		} else {
			entry.Type = "file"
			if info, statErr := item.Info(); statErr == nil {
				entry.Size = info.Size()
			}
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type == "directory"
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})
	return entries, nil
}

func (r *Repository) Read(ctx context.Context, relativePath string) (string, error) {
	if err := ctxErr(ctx); err != nil {
		return "", err
	}
	path, err := r.resolve(relativePath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat %q: %w", relativePath, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory: %s", relativePath)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", relativePath, err)
	}
	return string(data), nil
}

func (r *Repository) resolve(relativePath string) (string, error) {
	if relativePath == "" || relativePath == "." {
		return r.root, nil
	}
	if filepath.IsAbs(relativePath) {
		return "", errors.New("path must be repository-relative")
	}
	path := filepath.Clean(filepath.Join(r.root, relativePath))
	relative, err := filepath.Rel(r.root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes the knowledge root")
	}
	return path, nil
}

func ctxErr(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
