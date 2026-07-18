package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResolvesRootRelativeToConfig(t *testing.T) {
	configDir := t.TempDir()
	root := filepath.Join(configDir, "knowledge")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(configDir, "config.json")
	contents := `{"root":"./knowledge","ripgrep_path":"custom-rg"}`
	if err := os.WriteFile(configPath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, "")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Root != filepath.Clean(root) {
		t.Fatalf("root = %q, want %q", cfg.Root, root)
	}
	if cfg.RipgrepPath != "custom-rg" {
		t.Fatalf("ripgrep path = %q, want custom-rg", cfg.RipgrepPath)
	}
}

func TestLoadRequiresRoot(t *testing.T) {
	if _, err := Load("", ""); err == nil {
		t.Fatal("expected missing root error")
	}
}
