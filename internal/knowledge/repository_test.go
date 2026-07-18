package knowledge

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryListAndRead(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "product", "module"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "product", "README.md"), []byte("Product docs"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "product", "module", "feature.md"), []byte("Feature details"), 0o644); err != nil {
		t.Fatal(err)
	}

	repo, err := NewRepository(root)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := repo.List(context.Background(), ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name != "product" || entries[0].Type != "directory" {
		t.Fatalf("unexpected root entries: %#v", entries)
	}

	content, err := repo.Read(context.Background(), "product/module/feature.md")
	if err != nil {
		t.Fatal(err)
	}
	if content != "Feature details" {
		t.Fatalf("content = %q, want Feature details", content)
	}
}

func TestRepositoryRejectsPathOutsideRoot(t *testing.T) {
	repo, err := NewRepository(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repo.Read(context.Background(), "../outside.md"); err == nil {
		t.Fatal("expected path escape error")
	}
}

func TestRepositoryHonorsCanceledContext(t *testing.T) {
	repo, err := NewRepository(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := repo.List(ctx, "."); err == nil {
		t.Fatal("expected canceled context error")
	}
}
