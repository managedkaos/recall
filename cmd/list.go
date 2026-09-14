package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/frontmatter"
	"github.com/managedkaos/recall/internal/storage"
)

// runList lists all recall files sorted alphabetically. When tag is non-empty,
// only files carrying that tag (case-insensitive) are listed. This helper is
// invoked by the root command's --list / -l flag.
func runList(tag string) error {
	dir, err := config.RecallDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := config.EnsureDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	names, err := storage.List(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recall: cannot read directory %s: %v\n", dir, err)
		os.Exit(1)
	}

	if tag == "" {
		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	}

	// Filter by tag
	for _, name := range names {
		content, err := storage.Read(dir, name)
		if err != nil {
			continue
		}
		tags, _ := frontmatter.Parse(content)
		for _, t := range tags {
			if strings.EqualFold(t, tag) {
				fmt.Println(name)
				break
			}
		}
	}

	return nil
}
