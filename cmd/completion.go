package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// runCompletion writes a shell completion script for the named shell to stdout.
// Supported shells: bash, zsh, fish, powershell. An unsupported or empty shell
// value prints an error to stderr and exits with a non-zero code. Invoked by
// the root command's --completion / -c <shell> flag.
func runCompletion(cmd *cobra.Command, shell string) error {
	root := cmd.Root()

	switch strings.ToLower(shell) {
	case "bash":
		return root.GenBashCompletionV2(os.Stdout, true)
	case "zsh":
		return root.GenZshCompletion(os.Stdout)
	case "fish":
		return root.GenFishCompletion(os.Stdout, true)
	case "powershell":
		return root.GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "recall: unsupported shell %q (supported: bash, zsh, fish, powershell)\n", shell)
		os.Exit(1)
		return nil
	}
}
