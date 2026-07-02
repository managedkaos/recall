package cmd

import (
	"fmt"
	"os"

	"github.com/managedkaos/recall/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize the recall directory",
	Long: `Create the recall directory if it doesn't exist.

By default, creates ~/.recall (or the path specified by RECALL_DIR).
Optionally pass a path argument to create a specific directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var dir string
		var err error

		if len(args) == 1 {
			dir = args[0]
		} else {
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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
