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

// Action flags select a mode of operation other than the default
// look-up-and-render behavior. At most one may be set per invocation.
var (
	editFlag       bool
	searchFlag     bool
	listFlag       bool
	initFlag       bool
	versionFlag    bool
	completionFlag string

	// Modifiers
	rawFlag  bool
	tagFlag  string
	initPath string
)

var rootCmd = &cobra.Command{
	Use:   "recall [filename]",
	Short: "Store, retrieve, and search markdown reference files",
	Long:  `Recall is a CLI tool that stores, retrieves, edits, lists, and searches markdown-formatted reference files from the command line.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  runRecall,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	f := rootCmd.Flags()
	// Action flags
	f.BoolVarP(&editFlag, "edit", "e", false, "edit the specified file in $EDITOR")
	f.BoolVarP(&searchFlag, "search", "s", false, "search all recall files for the given query")
	f.BoolVarP(&listFlag, "list", "l", false, "list all recall files")
	f.BoolVarP(&initFlag, "init", "i", false, "initialize the recall directory")
	f.BoolVarP(&versionFlag, "version", "v", false, "print the version of recall")
	f.StringVarP(&completionFlag, "completion", "c", "", "generate a completion script for the given shell: bash|zsh|fish|powershell")

	// Modifiers
	f.BoolVarP(&rawFlag, "raw", "r", false, "output unformatted markdown without ANSI styling")
	f.StringVar(&tagFlag, "tag", "", "filter --list by tag (case-insensitive)")
	f.StringVar(&initPath, "init-path", "", "with --init, initialize the given directory instead of the default")
}

// runRecall is the single dispatch point for the recall command. It enforces
// that at most one action flag is set, then routes to the appropriate handler.
// With no action flag, it performs the default look-up-and-render behavior.
func runRecall(cmd *cobra.Command, args []string) error {
	// Enforce at most one action flag.
	actions := 0
	for _, on := range []bool{editFlag, searchFlag, listFlag, initFlag, versionFlag, completionFlag != ""} {
		if on {
			actions++
		}
	}
	if actions > 1 {
		fmt.Fprintln(os.Stderr, "recall: only one action flag (--edit, --search, --list, --init, --version, --completion) may be used at a time")
		os.Exit(1)
	}

	// --init-path requires --init.
	if initPath != "" && !initFlag {
		fmt.Fprintln(os.Stderr, "recall: --init-path requires --init")
		os.Exit(1)
	}

	switch {
	case versionFlag:
		return runVersion(cmd, args)

	case completionFlag != "":
		return runCompletion(cmd, completionFlag)

	case initFlag:
		return runInit(initPath)

	case listFlag:
		return runList(tagFlag)

	case searchFlag:
		if len(args) == 0 {
			return cmd.Help()
		}
		return runSearch(args[0])

	case editFlag:
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "recall: --edit requires a filename")
			os.Exit(1)
		}
		return runEdit(cmd, args)

	default:
		if len(args) == 0 {
			return cmd.Help()
		}
		return renderFile(args[0])
	}
}

// renderFile reads the named recall file, strips front-matter, and writes it to
// stdout. With --raw set, the body is written unmodified; otherwise it is
// rendered with ANSI styling. A missing file exits with a non-zero code and no
// output.
func renderFile(filename string) error {
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

	output, err := renderer.Render(body)
	if err != nil {
		os.Exit(1)
	}

	fmt.Print(output)
	return nil
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
