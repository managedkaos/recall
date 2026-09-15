package cmd

import (
	"strings"

	"github.com/managedkaos/recall/internal/config"
	"github.com/managedkaos/recall/internal/storage"
	"github.com/spf13/cobra"
)

// completeRecallFiles returns the stored recall file names that begin with
// toComplete, along with a directive telling the shell not to fall back to
// filesystem path completion. It resolves the recall directory the same way the
// runtime does (honoring RECALL_DIR) and degrades gracefully: if the directory
// cannot be determined or read, it returns no candidates rather than erroring.
func completeRecallFiles(toComplete string) ([]string, cobra.ShellCompDirective) {
	dir, err := config.RecallDir()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	names, err := storage.List(dir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var matches []string
	for _, name := range names {
		if strings.HasPrefix(name, toComplete) {
			matches = append(matches, name)
		}
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}

// completeArgs is the root command's ValidArgsFunction. It provides dynamic
// completion for positional arguments based on the active action flag:
//
//   - --list, --init, --version: no positional argument is used, so no
//     candidates are offered.
//   - --completion: the positional slot is not used (the shell is a flag
//     value), so no candidates are offered here.
//   - --search: the argument is a free-text query, not a filename, so no
//     candidates are offered.
//   - default, --edit, --raw, --metadata: complete stored recall file names.
//
// In every case ShellCompDirectiveNoFileComp is returned so the shell does not
// fall back to filesystem path completion.
//
// It inspects the package-level action-flag variables (listFlag, initFlag,
// etc.), which Cobra has already parsed from the command line by the time
// completion runs, to decide which candidates are appropriate.
func completeArgs(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch {
	case listFlag, initFlag, versionFlag, completionFlag != "", searchFlag:
		return nil, cobra.ShellCompDirectiveNoFileComp
	default:
		return completeRecallFiles(toComplete)
	}
}
