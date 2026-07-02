package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// RecallDir returns the absolute path to the recall directory.
// It checks RECALL_DIR env var first, then falls back to $HOME/.recall.
// Returns an error if the home directory cannot be determined.
func RecallDir() (string, error) {
	if dir := os.Getenv("RECALL_DIR"); dir != "" {
		return filepath.Clean(dir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("recall: cannot determine home directory: %w", err)
	}

	return filepath.Join(home, ".recall"), nil
}

// EnsureDir creates the recall directory (and intermediates) if it doesn't exist.
// Returns an error if the directory cannot be created or is not writable.
func EnsureDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("recall: cannot create directory %s: %w", dir, err)
	}
	return nil
}
