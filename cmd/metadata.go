package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/frontmatter"
	"github.com/managedkaos/recall/internal/storage"
)

// timeFormat is the timestamp layout used for human-readable text output.
// JSON output relies on time.Time's default RFC 3339 marshaling.
const timeFormat = time.RFC3339

// FileMetadata is the reported metadata for a single recall file. Time fields
// are pointers so that platform-unavailable times are omitted from output
// rather than shown as zero values.
type FileMetadata struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	SizeBytes int64      `json:"size_bytes"`
	Modified  time.Time  `json:"modified"`
	Created   *time.Time `json:"created,omitempty"`
	Changed   *time.Time `json:"changed,omitempty"`
	Tags      []string   `json:"tags"`
}

// collectMetadata gathers the reportable metadata for a single recall file:
// its absolute path, size, timestamps, and front-matter tags. It returns an
// error if the file cannot be stat'd or read.
func collectMetadata(dir, name string) (FileMetadata, error) {
	path := storage.FilePath(dir, name)

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}

	info, err := os.Stat(path)
	if err != nil {
		return FileMetadata{}, err
	}

	created, changed := platformTimes(info)

	content, err := storage.Read(dir, name)
	if err != nil {
		return FileMetadata{}, err
	}

	tags, _ := frontmatter.Parse(content)
	if tags == nil {
		tags = []string{}
	}

	return FileMetadata{
		Name:      name,
		Path:      abs,
		SizeBytes: info.Size(),
		Modified:  info.ModTime(),
		Created:   created,
		Changed:   changed,
		Tags:      tags,
	}, nil
}

// renderText writes a human-readable, labeled block for each item to stdout.
// Unavailable timestamps are omitted; files without tags show "(none)".
// Blocks are separated by a blank line.
func renderText(items []FileMetadata) {
	for i, md := range items {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("Name:     %s\n", md.Name)
		fmt.Printf("Path:     %s\n", md.Path)
		fmt.Printf("Size:     %d bytes\n", md.SizeBytes)
		fmt.Printf("Modified: %s\n", md.Modified.Format(timeFormat))
		if md.Created != nil {
			fmt.Printf("Created:  %s\n", md.Created.Format(timeFormat))
		}
		if md.Changed != nil {
			fmt.Printf("Changed:  %s\n", md.Changed.Format(timeFormat))
		}
		if len(md.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(md.Tags, ", "))
		} else {
			fmt.Printf("Tags:     (none)\n")
		}
	}
}

// renderJSON writes items to stdout as an indented JSON array. A nil slice is
// normalized so that empty results encode as "[]".
func renderJSON(items []FileMetadata) {
	if items == nil {
		items = []FileMetadata{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(items)
}

// runMetadata reports metadata for each named file that exists in the recall
// directory. Missing or unreadable files produce a stderr error but do not stop
// processing of the remaining names. Output is JSON when asJSON is true,
// otherwise human-readable text. The process exits with code 1 only when no
// file could be collected; otherwise it returns nil.
func runMetadata(names []string, asJSON bool) error {
	dir, err := config.RecallDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := config.EnsureDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var collected []FileMetadata
	for _, name := range names {
		if !storage.Exists(dir, name) {
			fmt.Fprintf(os.Stderr, "recall: file not found: %s\n", name)
			continue
		}
		md, err := collectMetadata(dir, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "recall: cannot read metadata for %s: %v\n", name, err)
			continue
		}
		collected = append(collected, md)
	}

	if asJSON {
		renderJSON(collected)
	} else {
		renderText(collected)
	}

	if len(collected) == 0 {
		os.Exit(1)
	}

	return nil
}
