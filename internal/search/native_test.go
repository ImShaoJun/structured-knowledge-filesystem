package search

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNativeSearcherReturnsStructuredMatches(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "product-alpha", "orders"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(root, "product-alpha", "orders", "payment-retry.md"),
		[]byte("# Payment\n\nThe order enters PAYMENT_FAILED before retry.\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("PAYMENT_FAILED appears here too.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored.json"), []byte("PAYMENT_FAILED"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := NewNativeSearcher().Search(context.Background(), root, "PAYMENT_FAILED", ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2: %#v", len(results), results)
	}
	if results[0].Path != "notes.txt" || results[0].Line != 1 {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Path != "product-alpha/orders/payment-retry.md" || results[1].Line != 3 {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
	if results[1].Column == 0 {
		t.Fatal("expected a match column")
	}
}

func TestNativeSearcherRejectsPathOutsideRoot(t *testing.T) {
	root := t.TempDir()
	if _, err := NewNativeSearcher().Search(context.Background(), root, "secret", "../outside"); err == nil {
		t.Fatal("expected search path escape error")
	}
}

func TestNativeSearcherRejectsInvalidRegularExpression(t *testing.T) {
	root := t.TempDir()
	if _, err := NewNativeSearcher().Search(context.Background(), root, "[", "."); err == nil {
		t.Fatal("expected invalid regular expression error")
	}
}

func TestNativeSearcherHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := NewNativeSearcher().Search(ctx, t.TempDir(), "term", "."); err == nil {
		t.Fatal("expected canceled context error")
	}
}

func TestNewSearcherUsesNativeByDefault(t *testing.T) {
	if _, ok := NewSearcher("").(*NativeSearcher); !ok {
		t.Fatal("expected native searcher by default")
	}
	if _, ok := NewSearcher("rg").(*RipgrepSearcher); !ok {
		t.Fatal("expected ripgrep searcher when configured")
	}
}
