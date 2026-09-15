package cmd

import (
	"fmt"

	"github.com/managedkaos/recall/internal/buildinfo"
	"github.com/spf13/cobra"
)

// Package-level variables set via -ldflags at build time.
var (
	Version string // Injected by Make or GoReleaser at build time.

	GitBranch        string
	BuildEnvironment string
	BuildDate        string
)

// runVersion prints the build/version metadata. Invoked by the root command's
// --version / -v flag.
func runVersion(cmd *cobra.Command, args []string) error {
	meta := buildinfo.Collect(Version, GitBranch, BuildEnvironment, BuildDate)
	fmt.Print(meta.String())
	return nil
}
