package cmd

import (
	"fmt"
	"os"

	"github.com/managedkaos/recall/internal/config"
)

// runInit creates the recall directory if it doesn't exist and prints a
// confirmation. When path is empty, the default directory (or RECALL_DIR) is
// used; otherwise path is used. Invoked by the root command's --init / -i flag,
// with an optional --init-path <dir> modifier.
func runInit(path string) error {
	dir := path
	if dir == "" {
		var err error
		dir, err = config.RecallDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := config.EnsureDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "recall: initialized directory at %s\n", dir)
	return nil
}
