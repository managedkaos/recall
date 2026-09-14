package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/storage"
	"github.com/spf13/cobra"
)

// runEdit opens a recall file in $EDITOR, creating it first if it doesn't
// exist. Invoked by the root command's --edit / -e flag.
func runEdit(cmd *cobra.Command, args []string) error {
	filename := args[0]

	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Fprintln(os.Stderr, "recall: $EDITOR environment variable is not set")
		os.Exit(1)
	}

	dir, err := config.RecallDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := config.EnsureDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !storage.Exists(dir, filename) {
		if err := storage.Create(dir, filename); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("recall: cannot create file %s: %v", filename, err))
			os.Exit(1)
		}
	}

	filePath := storage.FilePath(dir, filename)

	editorCmd := exec.Command(editor, filePath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, fmt.Sprintf("recall: failed to run editor: %v", err))
		os.Exit(1)
	}

	return nil
}
