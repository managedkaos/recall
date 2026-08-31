package cmd

import (
	"fmt"
	"os"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/search"
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

		return runSearch(args[0])
	},
}

// runSearch performs a case-insensitive search across all recall files for
// query and prints the results to stdout in the standard
// "filename:linenumber:linecontent" format, grouping matches by file and
// separating file groups with a line containing "----------". Empty results
// produce no output. Directory-resolution and search errors are printed to
// stderr and cause a non-zero exit. This helper is shared by the "search"
// subcommand and the root command's --search / -s flag.
func runSearch(query string) error {
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
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
