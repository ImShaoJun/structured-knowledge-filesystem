package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config contains the filesystem root and an optional ripgrep executable path.
// When omitted, the MCP server uses its built-in search backend.
type Config struct {
	Root        string `json:"root"`
	RipgrepPath string `json:"ripgrep_path,omitempty"`
}

// Load reads a JSON configuration file, applies command-line overrides, and
// normalizes the knowledge root to an absolute path. Relative roots in a config
// file are resolved relative to that config file rather than the process cwd.
func Load(configPath, rootFlag string) (Config, error) {
	var cfg Config

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
		if cfg.Root != "" && !filepath.IsAbs(cfg.Root) {
			cfg.Root = filepath.Join(filepath.Dir(configPath), cfg.Root)
		}
	}

	if rootFlag != "" {
		cfg.Root = rootFlag
	}
	if cfg.Root == "" {
		return Config{}, errors.New("knowledge root is required: use --root or --config")
	}

	root, err := filepath.Abs(cfg.Root)
	if err != nil {
		return Config{}, fmt.Errorf("resolve knowledge root: %w", err)
	}
	cfg.Root = filepath.Clean(root)

	return cfg, nil
}
