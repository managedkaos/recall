package buildinfo

import (
	"runtime/debug"
	"strings"
	"testing"
)

func TestFormatVersion_AllPopulated(t *testing.T) {
	got := FormatVersion("0", "1", "0")
	if got != "0.1.0" {
		t.Errorf("expected '0.1.0', got %q", got)
	}
}

func TestFormatVersion_LargeNumbers(t *testing.T) {
	got := FormatVersion("12", "34", "56")
	if got != "12.34.56" {
		t.Errorf("expected '12.34.56', got %q", got)
	}
}

func TestFormatVersion_EmptyMajor(t *testing.T) {
	got := FormatVersion("", "1", "0")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_EmptyMinor(t *testing.T) {
	got := FormatVersion("0", "", "0")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_EmptyPatch(t *testing.T) {
	got := FormatVersion("0", "1", "")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_AllEmpty(t *testing.T) {
	got := FormatVersion("", "", "")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestMergeMetadata_UsesBuildInfoSettings(t *testing.T) {
	info := &debug.BuildInfo{
		GoVersion: "go1.26.4",
		Main: debug.Module{
			Path: "github.com/managedkaos/recall",
		},
		Settings: []debug.BuildSetting{
			{Key: "GOOS", Value: "linux"},
			{Key: "GOARCH", Value: "amd64"},
			{Key: "vcs.revision", Value: "abc1234"},
			{Key: "vcs.time", Value: "2026-07-21T20:55:00Z"},
			{Key: "vcs.modified", Value: "false"},
			{Key: "CGO_ENABLED", Value: "0"},
		},
	}

	got := mergeMetadata("0", "1", "0", "main", "GitHub Actions", "2026-01-01T00:00:00Z", info)

	if got.Version != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %q", got.Version)
	}
	if got.GoVersion != "go1.26.4" {
		t.Errorf("expected GoVersion go1.26.4, got %q", got.GoVersion)
	}
	if got.Platform != "linux/amd64" {
		t.Errorf("expected platform linux/amd64, got %q", got.Platform)
	}
	if got.Built != "2026-07-21T20:55:00Z" {
		t.Errorf("expected Built from vcs.time, got %q", got.Built)
	}
	if got.Environment != "GitHub Actions" {
		t.Errorf("expected Environment GitHub Actions, got %q", got.Environment)
	}
	if got.Branch != "main" {
		t.Errorf("expected Branch main, got %q", got.Branch)
	}
	if got.Commit != "abc1234" {
		t.Errorf("expected Commit abc1234, got %q", got.Commit)
	}
	if got.Modified != "false" {
		t.Errorf("expected Modified false, got %q", got.Modified)
	}
	if got.Module != "github.com/managedkaos/recall" {
		t.Errorf("expected module path, got %q", got.Module)
	}
	if got.CGOEnabled != "0" {
		t.Errorf("expected CGO_ENABLED 0, got %q", got.CGOEnabled)
	}
}

func TestMergeMetadata_FallsBackToLdflags(t *testing.T) {
	got := mergeMetadata("0", "1", "0", "feature/test", "local (Darwin)", "2026-07-21T12:00:00Z", nil)

	if got.Built != "2026-07-21T12:00:00Z" {
		t.Errorf("expected Built from ldflags, got %q", got.Built)
	}
	if got.Branch != "feature/test" {
		t.Errorf("expected Branch from ldflags, got %q", got.Branch)
	}
	if got.Environment != "local (Darwin)" {
		t.Errorf("expected Environment from ldflags, got %q", got.Environment)
	}
	if got.Commit != "unknown" {
		t.Errorf("expected Commit unknown, got %q", got.Commit)
	}
	if got.Modified != "unknown" {
		t.Errorf("expected Modified unknown, got %q", got.Modified)
	}
}

func TestMergeMetadata_MissingLdflagsAreUnknown(t *testing.T) {
	got := mergeMetadata("", "", "", "", "", "", nil)

	if got.Version != "unknown" {
		t.Errorf("expected version unknown, got %q", got.Version)
	}
	if got.Environment != "unknown" {
		t.Errorf("expected Environment unknown, got %q", got.Environment)
	}
	if got.Branch != "unknown" {
		t.Errorf("expected Branch unknown, got %q", got.Branch)
	}
	if got.Built != "unknown" {
		t.Errorf("expected Built unknown, got %q", got.Built)
	}
}

func TestMetadataString_ContainsExpectedLines(t *testing.T) {
	meta := Metadata{
		Version:     "0.1.0",
		GoVersion:   "go1.26.4",
		Platform:    "darwin/arm64",
		Built:       "2026-07-21T20:55:00Z",
		Environment: "local (Darwin)",
		Branch:      "main",
		Commit:      "abc1234",
		Modified:    "false",
		Module:      "github.com/managedkaos/recall",
		CGOEnabled:  "false",
	}

	out := meta.String()
	for _, want := range []string{
		"recall version 0.1.0",
		"Go version:     go1.26.4",
		"Platform:       darwin/arm64",
		"Environment:    local (Darwin)",
		"Module:         github.com/managedkaos/recall",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
