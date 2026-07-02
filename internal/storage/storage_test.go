package storage_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mjenkins/recall/internal/storage"
)

func TestList_RegularFiles(t *testing.T) {
	tmp := t.TempDir()

	// Create some regular files
	for _, name := range []string{"banana", "apple", "cherry"} {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := storage.List(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"apple", "banana", "cherry"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(got), got)
	}
	for i, name := range expected {
		if got[i] != name {
			t.Errorf("index %d: expected %q, got %q", i, name, got[i])
		}
	}
}

func TestList_ExcludesHiddenFiles(t *testing.T) {
	tmp := t.TempDir()

	os.WriteFile(filepath.Join(tmp, ".hidden"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(tmp, "visible"), []byte("x"), 0o644)

	got, err := storage.List(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 || got[0] != "visible" {
		t.Errorf("expected [visible], got %v", got)
	}
}

func TestList_ExcludesDirectories(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(filepath.Join(tmp, "subdir"), 0o755)
	os.WriteFile(filepath.Join(tmp, "file"), []byte("x"), 0o644)

	got, err := storage.List(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 || got[0] != "file" {
		t.Errorf("expected [file], got %v", got)
	}
}

func TestList_EmptyDirectory(t *testing.T) {
	tmp := t.TempDir()

	got, err := storage.List(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected empty list, got %v", got)
	}
}

func TestList_NonExistentDirectory(t *testing.T) {
	_, err := storage.List("/nonexistent/path/xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}

func TestList_SortedAlphabetically(t *testing.T) {
	tmp := t.TempDir()

	names := []string{"zebra", "mango", "alpha", "beta"}
	for _, name := range names {
		os.WriteFile(filepath.Join(tmp, name), []byte("x"), 0o644)
	}

	got, err := storage.List(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"alpha", "beta", "mango", "zebra"}
	for i, name := range expected {
		if got[i] != name {
			t.Errorf("index %d: expected %q, got %q", i, name, got[i])
		}
	}
}

func TestRead_ExistingFile(t *testing.T) {
	tmp := t.TempDir()
	content := []byte("hello world\n")
	os.WriteFile(filepath.Join(tmp, "test"), content, 0o644)

	got, err := storage.Read(tmp, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("expected %q, got %q", content, got)
	}
}

func TestRead_NonExistentFile(t *testing.T) {
	tmp := t.TempDir()

	_, err := storage.Read(tmp, "missing")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestFilePath_UsesNativeSeparator(t *testing.T) {
	got := storage.FilePath("/home/user/.recall", "notes")

	var expected string
	if runtime.GOOS == "windows" {
		expected = "\\home\\user\\.recall\\notes"
	} else {
		expected = "/home/user/.recall/notes"
	}
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFilePath_JoinsCorrectly(t *testing.T) {
	got := storage.FilePath("/a/b", "c")
	expected := filepath.Join("/a/b", "c")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExists_FileExists(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "present"), []byte("x"), 0o644)

	if !storage.Exists(tmp, "present") {
		t.Error("expected Exists to return true for existing file")
	}
}

func TestExists_FileDoesNotExist(t *testing.T) {
	tmp := t.TempDir()

	if storage.Exists(tmp, "absent") {
		t.Error("expected Exists to return false for nonexistent file")
	}
}

func TestExists_DirectoryNotFile(t *testing.T) {
	tmp := t.TempDir()
	os.Mkdir(filepath.Join(tmp, "subdir"), 0o755)

	if storage.Exists(tmp, "subdir") {
		t.Error("expected Exists to return false for a directory")
	}
}

func TestCreate_NewFile(t *testing.T) {
	tmp := t.TempDir()

	err := storage.Create(tmp, "newfile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was created
	info, err := os.Stat(filepath.Join(tmp, "newfile"))
	if err != nil {
		t.Fatalf("file was not created: %v", err)
	}
	if info.IsDir() {
		t.Error("expected a file, not a directory")
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}

func TestCreate_InvalidPath(t *testing.T) {
	err := storage.Create("/nonexistent/path/xyz", "file")
	if err == nil {
		t.Fatal("expected error for invalid directory")
	}
}
