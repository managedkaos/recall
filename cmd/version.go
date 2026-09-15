package cmd

import (
	"fmt"

	"github.com/managedkaos/recall/internal/buildinfo"
	"github.com/spf13/cobra"
)

// Package-level variables set via -ldflags at build time.
var (
	Version string // Full semantic version, injected from the git tag by GoReleaser.
	Major   string
	Minor   string
	Patch   string

	GitBranch        string
	BuildEnvironment string
	BuildDate        string
)

// runVersion prints the build/version metadata. Invoked by the root command's
// --version / -v flag.
func runVersion(cmd *cobra.Command, args []string) error {
	meta := buildinfo.Collect(Version, Major, Minor, Patch, GitBranch, BuildEnvironment, BuildDate)
	fmt.Print(meta.String())
	return nil
}

// formatVersion constructs the semantic version string from components.
// If any component is empty (ldflags not provided), returns "unknown".
func formatVersion(major, minor, patch string) string {
	return buildinfo.FormatVersion(major, minor, patch)
}
