package cmd

import (
	"fmt"
	"os"

	"github.com/mjenkins/recall/internal/config"
	"github.com/mjenkins/recall/internal/search"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search all recall files for a string",
	Long:  `Search all recall files in the recall directory for a case-insensitive substring match.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "recall: search requires a query argument")
			os.Exit(1)
		}

		query := args[0]

		dir, err := config.RecallDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := config.EnsureDir(dir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		results, err := search.Search(dir, query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "recall: search failed: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			return nil
		}

		for i, fileResult := range results {
			if i > 0 {
				fmt.Fprintf(os.Stdout, "----------\n")
			}
			for _, match := range fileResult.Matches {
				fmt.Fprintf(os.Stdout, "%s:%d:%s\n", match.Filename, match.LineNum, match.Line)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
