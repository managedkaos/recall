package cmd

import (
	"fmt"

	"github.com/managedkaos/recall/internal/buildinfo"
	"github.com/spf13/cobra"
)

// Package-level variables set via -ldflags at build time.
var (
	Version string // Composite version string (reserved for future use)
	Major   string
	Minor   string
	Patch   string
	GitBranch         string
	BuildEnvironment  string
	BuildDate         string
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
	meta := buildinfo.Collect(Major, Minor, Patch, GitBranch, BuildEnvironment, BuildDate)
	fmt.Print(meta.String())
	return nil
}

// formatVersion constructs the semantic version string from components.
// If any component is empty (ldflags not provided), returns "unknown".
func formatVersion(major, minor, patch string) string {
	return buildinfo.FormatVersion(major, minor, patch)
}
