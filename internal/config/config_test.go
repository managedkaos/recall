package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/managedkaos/recall/internal/config"
	"pgregory.net/rapid"
)

func TestPlaceholder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Placeholder to ensure rapid dependency is retained in go.mod.
		_ = rapid.Int().Draw(t, "n")
	})
}

func TestRecallDir_EnvVarSet(t *testing.T) {
	t.Setenv("RECALL_DIR", "/custom/path")

	dir, err := config.RecallDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "/custom/path" {
		t.Errorf("expected /custom/path, got %s", dir)
	}
}

func TestRecallDir_EnvVarEmpty(t *testing.T) {
	t.Setenv("RECALL_DIR", "")

	dir, err := config.RecallDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".recall")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestRecallDir_DefaultFallback(t *testing.T) {
	t.Setenv("RECALL_DIR", "")

	dir, err := config.RecallDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".recall")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestRecallDir_EnvVarCleaned(t *testing.T) {
	t.Setenv("RECALL_DIR", "/custom/path/../other")

	dir, err := config.RecallDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Clean("/custom/path/../other")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestEnsureDir_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "a", "b", "c")

	err := config.EnsureDir(target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory")
	}
}

func TestEnsureDir_ExistingDirectory(t *testing.T) {
	tmp := t.TempDir()

	err := config.EnsureDir(tmp)
	if err != nil {
		t.Fatalf("unexpected error for existing dir: %v", err)
	}
}

func TestEnsureDir_InvalidPath(t *testing.T) {
	// Try to create a directory under a file (should fail)
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "afile")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(filePath, "subdir")
	err := config.EnsureDir(target)
	if err == nil {
		t.Fatal("expected error when creating directory under a file")
	}
}
