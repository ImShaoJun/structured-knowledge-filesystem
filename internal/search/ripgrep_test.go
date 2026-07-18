package search

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRipgrepSearcherReturnsStructuredMatches(t *testing.T) {
	rgPath, err := exec.LookPath("rg")
	if err != nil {
		t.Skip("ripgrep is not installed")
	}

	root := t.TempDir()
	file := filepath.Join(root, "product-alpha", "orders", "payment-retry.md")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Payment\n\nThe order enters PAYMENT_FAILED before retry.\n"
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored.json"), []byte("PAYMENT_FAILED"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := NewRipgrepSearcher(rgPath).Search(context.Background(), root, "PAYMENT_FAILED", ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1: %#v", len(results), results)
	}
	if results[0].Path != "product-alpha/orders/payment-retry.md" {
		t.Fatalf("path = %q", results[0].Path)
	}
	if results[0].Line != 3 {
		t.Fatalf("line = %d, want 3", results[0].Line)
	}
	if results[0].Column == 0 {
		t.Fatal("expected a match column")
	}
}

func TestRipgrepSearcherReturnsEmptyForNoMatch(t *testing.T) {
	rgPath, err := exec.LookPath("rg")
	if err != nil {
		t.Skip("ripgrep is not installed")
	}

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("nothing relevant\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := NewRipgrepSearcher(rgPath).Search(context.Background(), root, "not-present", ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("got results for missing query: %#v", results)
	}
}
