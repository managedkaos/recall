package cmd

import (
	"os"

	"github.com/spf13/cobra/doc"
)

// GenerateMan writes the manual from the same command definition used by the CLI.
func GenerateMan(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	rootCmd.InitDefaultHelpFlag()
	return doc.GenManTree(rootCmd, &doc.GenManHeader{
		Title: "RECALL", Section: "1", Source: "Recall", Manual: "Recall Manual",
	}, dir)
}
