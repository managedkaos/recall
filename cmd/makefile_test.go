package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the real Makefile in isolated repositories without compiling a binary
// for each Git state. CLI tests separately verify the injected version display.
func TestMakeVersion(t *testing.T) {
	for _, tool := range []string{"git", "make"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is required: %v", tool, err)
		}
	}
	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join(getProjectRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Makefile"), data, 0600); err != nil {
		t.Fatal(err)
	}
	// Keep ambient overrides and Git configuration from influencing fixtures.
	var env []string
	for _, entry := range os.Environ() {
		key := strings.SplitN(entry, "=", 2)[0]
		if key == "VERSION" || key == "MAKEFLAGS" || key == "MFLAGS" || key == "MAKEOVERRIDES" || strings.HasPrefix(key, "GIT_") {
			continue
		}
		env = append(env, entry)
	}
	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull)
	runGit := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir, cmd.Env = dir, env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	check := func(name, want string, extraEnv []string, args ...string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command("make", append([]string{"--no-print-directory", "-n", "build"}, args...)...)
			cmd.Dir = dir
			cmd.Env = append(append([]string{}, env...), extraEnv...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("make: %v\n%s", err, out)
			}
			if !strings.Contains(string(out), "/cmd.Version="+want+"'") {
				t.Fatalf("expected injected version %q, got:\n%s", want, out)
			}
		})
	}
	check("outside repository", "unknown", nil)
	check("environment override without Git metadata", "env-version", []string{"VERSION=env-version"})
	check("command override wins", "cli-version", []string{"VERSION=env-version"}, "VERSION=cli-version")
	check("empty override", "", nil, "VERSION=")
	runGit("init")
	runGit("config", "user.name", "Version Test")
	runGit("config", "user.email", "version-test@example.invalid")
	runGit("add", "Makefile")
	runGit("-c", "commit.gpgsign=false", "commit", "-m", "initial")
	check("no tags", "unknown", nil)
	runGit("tag", "unrelated")
	check("no matching tags", "unknown", nil)
	runGit("tag", "v1.2.3")
	check("lightweight tag", "v1.2.3", nil)
	if err := os.WriteFile(filepath.Join(dir, "untracked"), []byte("ignored for local suffix"), 0600); err != nil {
		t.Fatal(err)
	}
	check("untracked file stays clean", "v1.2.3", nil)
	if err := os.WriteFile(filepath.Join(dir, "Makefile"), append(data, []byte("\n# fixture change\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	check("unstaged change", "v1.2.3-local", nil)
	runGit("add", "Makefile")
	check("staged change", "v1.2.3-local", nil)
	runGit("-c", "commit.gpgsign=false", "commit", "-m", "next")
	hash := runGit("rev-parse", "--short", "HEAD")
	check("after tag", "v1.2.3-1-g"+hash, nil)
	runGit("-c", "tag.gpgsign=false", "tag", "-a", "V2.0.0-rc.1", "-m", "prerelease")
	check("annotated uppercase prerelease", "V2.0.0-rc.1", nil)
	runGit("-c", "commit.gpgsign=false", "commit", "--allow-empty", "-m", "date version")
	runGit("tag", "2026.09.15-1")
	check("numeric date tag", "2026.09.15-1", nil)
}
