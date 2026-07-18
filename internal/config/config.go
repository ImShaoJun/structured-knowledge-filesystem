package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Root        string `json:"root"`
	RipgrepPath string `json:"ripgrep_path,omitempty"`
}

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
