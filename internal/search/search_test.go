package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchContent_BasicMatch(t *testing.T) {
	content := []byte("Hello World\nfoo bar\nHello again")
	results := SearchContent("test", content, "hello")

	if len(results) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(results))
	}
	if results[0].LineNum != 1 {
		t.Errorf("expected line 1, got %d", results[0].LineNum)
	}
	if results[0].Line != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", results[0].Line)
	}
	if results[1].LineNum != 3 {
		t.Errorf("expected line 3, got %d", results[1].LineNum)
	}
	if results[1].Line != "Hello again" {
		t.Errorf("expected 'Hello again', got %q", results[1].Line)
	}
}

func TestSearchContent_CaseInsensitive(t *testing.T) {
	content := []byte("GO is great\ngo routines\nGo channels")
	results := SearchContent("test", content, "GO")

	if len(results) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(results))
	}
}

func TestSearchContent_NoMatch(t *testing.T) {
	content := []byte("Hello World\nfoo bar")
	results := SearchContent("test", content, "xyz")

	if len(results) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(results))
	}
}

func TestSearchContent_EmptyQuery(t *testing.T) {
	content := []byte("Hello World\nfoo bar")
	results := SearchContent("test", content, "")

	if results != nil {
		t.Fatalf("expected nil results for empty query, got %v", results)
	}
}

func TestSearchContent_EmptyContent(t *testing.T) {
	results := SearchContent("test", []byte(""), "hello")

	if len(results) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(results))
	}
}

func TestSearchContent_FilenamePreserved(t *testing.T) {
	content := []byte("match here")
	results := SearchContent("myfile", content, "match")

	if len(results) != 1 {
		t.Fatalf("expected 1 match, got %d", len(results))
	}
	if results[0].Filename != "myfile" {
		t.Errorf("expected filename 'myfile', got %q", results[0].Filename)
	}
}

func TestSearchContent_AscendingLineOrder(t *testing.T) {
	content := []byte("alpha hit\nno dice\nbeta hit\nno dice\ngamma hit")
	results := SearchContent("test", content, "hit")

	if len(results) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].LineNum <= results[i-1].LineNum {
			t.Errorf("results not in ascending order: line %d after line %d", results[i].LineNum, results[i-1].LineNum)
		}
	}
}

func TestSearch_MultipleFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	writeFile(t, dir, "alpha", "hello world\nsecond line")
	writeFile(t, dir, "beta", "no match here\nHELLO there")

	results, err := Search(dir, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 file results, got %d", len(results))
	}

	// Files should be in alphabetical order (from storage.List)
	if results[0].Filename != "alpha" {
		t.Errorf("expected first file 'alpha', got %q", results[0].Filename)
	}
	if results[1].Filename != "beta" {
		t.Errorf("expected second file 'beta', got %q", results[1].Filename)
	}

	// Check matches
	if len(results[0].Matches) != 1 {
		t.Errorf("expected 1 match in alpha, got %d", len(results[0].Matches))
	}
	if len(results[1].Matches) != 1 {
		t.Errorf("expected 1 match in beta, got %d", len(results[1].Matches))
	}
}

func TestSearch_NoMatches(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "alpha", "hello world")

	results, err := Search(dir, "xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 file results, got %d", len(results))
	}
}

func TestSearch_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	results, err := Search(dir, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 file results, got %d", len(results))
	}
}

func TestSearch_SkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".hidden", "hello match")
	writeFile(t, dir, "visible", "hello match")

	results, err := Search(dir, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 file result, got %d", len(results))
	}
	if results[0].Filename != "visible" {
		t.Errorf("expected filename 'visible', got %q", results[0].Filename)
	}
}

func TestSearch_InvalidDirectory(t *testing.T) {
	_, err := Search("/nonexistent/path", "hello")
	if err == nil {
		t.Fatal("expected error for invalid directory")
	}
}

func TestSearch_GroupedByFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "fileA", "apple pie\napple sauce\nbanana split")
	writeFile(t, dir, "fileB", "apple cider")

	results, err := Search(dir, "apple")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 file results, got %d", len(results))
	}

	// fileA should have 2 matches
	if results[0].Filename != "fileA" {
		t.Errorf("expected first file 'fileA', got %q", results[0].Filename)
	}
	if len(results[0].Matches) != 2 {
		t.Errorf("expected 2 matches in fileA, got %d", len(results[0].Matches))
	}

	// fileB should have 1 match
	if results[1].Filename != "fileB" {
		t.Errorf("expected second file 'fileB', got %q", results[1].Filename)
	}
	if len(results[1].Matches) != 1 {
		t.Errorf("expected 1 match in fileB, got %d", len(results[1].Matches))
	}
}

// writeFile is a helper to create a test file.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file %s: %v", name, err)
	}
}
