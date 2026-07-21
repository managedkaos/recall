package buildinfo

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Metadata holds build and version information for the recall binary.
type Metadata struct {
	Version     string
	GoVersion   string
	Platform    string
	Built       string
	Environment string
	Branch      string
	Commit      string
	Modified    string
	Module      string
	CGOEnabled  string
}

// FormatVersion constructs the semantic version string from components.
// If any component is empty (ldflags not provided), returns "unknown".
func FormatVersion(major, minor, patch string) string {
	if major == "" || minor == "" || patch == "" {
		return "unknown"
	}
	return major + "." + minor + "." + patch
}

// Collect gathers build metadata from embedded build info and ldflags fallbacks.
func Collect(major, minor, patch, gitBranch, buildEnv, buildDate string) Metadata {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return mergeMetadata(major, minor, patch, gitBranch, buildEnv, buildDate, nil)
	}
	return mergeMetadata(major, minor, patch, gitBranch, buildEnv, buildDate, info)
}

func mergeMetadata(major, minor, patch, gitBranch, buildEnv, buildDate string, info *debug.BuildInfo) Metadata {
	goos := settingValue(info, "GOOS")
	goarch := settingValue(info, "GOARCH")
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	platform := goos + "/" + goarch

	modified := settingValue(info, "vcs.modified")
	if modified == "" {
		modified = "unknown"
	}

	cgoEnabled := settingValue(info, "CGO_ENABLED")
	if cgoEnabled == "" {
		cgoEnabled = "unknown"
	}

	goVersion := "unknown"
	module := "unknown"
	if info != nil {
		if info.GoVersion != "" {
			goVersion = info.GoVersion
		}
		if info.Main.Path != "" {
			module = info.Main.Path
		}
	}

	return Metadata{
		Version:     FormatVersion(major, minor, patch),
		GoVersion:   goVersion,
		Platform:    platform,
		Built:       firstNonEmpty(settingValue(info, "vcs.time"), buildDate),
		Environment: valueOrUnknown(buildEnv),
		Branch:      valueOrUnknown(gitBranch),
		Commit:      valueOrUnknown(settingValue(info, "vcs.revision")),
		Modified:    modified,
		Module:      module,
		CGOEnabled:  cgoEnabled,
	}
}

// String returns the multi-line version output.
func (m Metadata) String() string {
	return fmt.Sprintf(
		"recall version %s\n\n"+
			"Go version:     %s\n"+
			"Platform:       %s\n"+
			"Built:          %s\n"+
			"Environment:    %s\n\n"+
			"Branch:         %s\n"+
			"Commit:         %s\n"+
			"Modified:       %s\n\n"+
			"Module:         %s\n"+
			"CGO enabled:    %s\n",
		m.Version,
		m.GoVersion,
		m.Platform,
		m.Built,
		m.Environment,
		m.Branch,
		m.Commit,
		m.Modified,
		m.Module,
		m.CGOEnabled,
	)
}

func settingValue(info *debug.BuildInfo, key string) string {
	if info == nil {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == key {
			return s.Value
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return "unknown"
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}
