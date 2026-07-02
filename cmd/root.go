package cmd

import (
	"fmt"
	"os"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/frontmatter"
	"github.com/managedkaos/recall/internal/renderer"
	"github.com/managedkaos/recall/internal/storage"
	"github.com/spf13/cobra"
)

// reservedNames are subcommand names that cannot be used as filenames.
var reservedNames = map[string]bool{
	"edit":   true,
	"list":   true,
	"search": true,
	"init":   true,
}

// IsReservedName checks if a filename conflicts with a subcommand name.
func IsReservedName(name string) bool {
	return reservedNames[name]
}

var editFlag bool
var rawFlag bool

var rootCmd = &cobra.Command{
	Use:   "recall [filename]",
	Short: "Store, retrieve, and search markdown reference files",
	Long:  `Recall is a CLI tool that stores, retrieves, edits, lists, and searches markdown-formatted reference files from the command line.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  runRecall,
	// Prevent Cobra from interpreting -e on subcommands
	TraverseChildren: true,
}

func init() {
	rootCmd.Flags().BoolVarP(&editFlag, "edit", "e", false, "edit the specified file")
	rootCmd.Flags().BoolVarP(&rawFlag, "raw", "r", false, "output unformatted markdown without ANSI styling")
}

func runRecall(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	filename := args[0]

	// Mutual exclusivity check
	if rawFlag && editFlag {
		fmt.Fprintln(os.Stderr, "recall: --raw and --edit flags cannot be used together")
		os.Exit(1)
	}

	// If -e flag is set, delegate to edit logic
	if editFlag {
		return runEdit(cmd, args)
	}

	// Resolve the recall directory
	dir, err := config.RecallDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := config.EnsureDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Check if file exists; if not, exit silently with non-zero code
	if !storage.Exists(dir, filename) {
		os.Exit(1)
	}

	// Read the file
	content, err := storage.Read(dir, filename)
	if err != nil {
		os.Exit(1)
	}

	// Strip front-matter
	_, body := frontmatter.Parse(content)

	// Raw output: write body directly without rendering
	if rawFlag {
		os.Stdout.Write(body)
		return nil
	}

	// Render the markdown
	output, err := renderer.Render(body)
	if err != nil {
		os.Exit(1)
	}

	// Print to stdout
	fmt.Print(output)
	return nil
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
