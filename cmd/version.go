package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Package-level variables set via -ldflags at build time.
var (
	Version string // Composite version string (reserved for future use)
	Major   string
	Minor   string
	Patch   string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of recall",
	Long:  "Print the version of the recall binary and exit.",
	Args:  cobra.NoArgs,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("recall version %s\n", formatVersion(Major, Minor, Patch))
	return nil
}

// formatVersion constructs the semantic version string from components.
// If any component is empty (ldflags not provided), returns "unknown".
func formatVersion(major, minor, patch string) string {
	if major == "" || minor == "" || patch == "" {
		return "unknown"
	}
	return major + "." + minor + "." + patch
}
