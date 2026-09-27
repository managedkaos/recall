package storage

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// List returns all regular (non-hidden, non-directory) filenames in
// the recall directory, sorted alphabetically.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files (those starting with '.')
		if strings.HasPrefix(name, ".") {
			continue
		}
		// Skip directories
		if entry.IsDir() {
			continue
		}
		names = append(names, name)
	}

	sort.Strings(names)
	return names, nil
}

// Read returns the content of a recall file by name.
// Returns an error if the file does not exist.
func Read(dir string, name string) ([]byte, error) {
	return os.ReadFile(FilePath(dir, name))
}

// FilePath constructs the full path to a recall file using the OS
// native path separator.
func FilePath(dir string, name string) string {
	return filepath.Join(dir, name)
}

// Exists checks whether a recall file exists.
func Exists(dir string, name string) bool {
	info, err := os.Stat(FilePath(dir, name))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Create creates a new empty file in the recall directory.
func Create(dir string, name string) error {
	f, err := os.Create(FilePath(dir, name))
	if err != nil {
		return err
	}
	return f.Close()
}

// FrontmatterScaffold is the initial content written to newly created files:
// an empty YAML front matter block followed by a blank line so the body starts
// cleanly.
const FrontmatterScaffold = "---\ntags: []\n---\n\n"

// CreateWithFrontmatter creates a new file pre-populated with an empty YAML
// front matter block. It is used when creating a file for editing so new files
// are consistent with the front matter format.
func CreateWithFrontmatter(dir string, name string) error {
	return os.WriteFile(FilePath(dir, name), []byte(FrontmatterScaffold), 0o644)
}
