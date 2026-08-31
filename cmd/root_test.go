package cmd_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the recall binary into a temp directory and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "recall")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(getProjectRoot(t))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

// buildBinaryWithMetadata compiles recall with ldflags for version metadata tests.
func buildBinaryWithMetadata(t *testing.T) string {
	t.Helper()
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "recall")
	module := "github.com/managedkaos/recall"
	ldflags := fmt.Sprintf(
		"-X %s/cmd.Major=0 -X %s/cmd.Minor=1 -X %s/cmd.Patch=0 -X %s/cmd.GitBranch=test-branch -X %s/cmd.BuildDate=2026-07-21T12:00:00Z -X %s/cmd.BuildEnvironment=test",
		module, module, module, module, module, module,
	)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", binPath, ".")
	cmd.Dir = filepath.Join(getProjectRoot(t))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary with metadata: %v\n%s", err, out)
	}
	return binPath
}

// getProjectRoot returns the project root directory.
func getProjectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from the test file's location to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}
	// We're in cmd/, go up one level
	return filepath.Dir(dir)
}

// setupRecallDir creates a temp recall directory with test files.
func setupRecallDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// File with front-matter
	err := os.WriteFile(filepath.Join(dir, "hello"), []byte("tags: greeting\nHello, world!\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// File without front-matter
	err = os.WriteFile(filepath.Join(dir, "plain"), []byte("Just plain content.\nSecond line.\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

// runRecall executes the recall binary with given args and env, returning stdout, stderr, and exit code.
func runRecall(t *testing.T, binPath, recallDir string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "RECALL_DIR="+recallDir)

	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running recall: %v", err)
		}
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode
}

func TestRawFlag_ValidFile(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--raw", "hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	expected := "Hello, world!\n"
	if stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout)
	}
}

func TestRawFlag_ShorthandR(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "-r", "hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	expected := "Hello, world!\n"
	if stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout)
	}
}

func TestRawFlag_FileWithoutFrontmatter(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--raw", "plain")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	expected := "Just plain content.\nSecond line.\n"
	if stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout)
	}
}

func TestRawFlag_WithEditFlagProducesError(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "--raw", "--edit", "hello")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "--raw and --edit") {
		t.Errorf("expected error about mutually exclusive flags, got stderr: %q", stderr)
	}
}

func TestRawFlag_WithoutFilenameShowsHelp(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--raw")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	// Help text should mention the usage
	if !strings.Contains(stdout, "recall") {
		t.Errorf("expected help output containing 'recall', got: %q", stdout)
	}
}

func TestRawFlag_NonExistentFileExits1Silently(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "--raw", "nonexistent")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no stdout, got: %q", stdout)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got: %q", stderr)
	}
}

func TestRawFlag_RejectedOnEditSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "edit", "--raw", "hello")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --raw used on edit subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestRawFlag_RejectedOnListSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "list", "--raw")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --raw used on list subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestRawFlag_RejectedOnSearchSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "search", "--raw", "hello")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --raw used on search subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestRawFlag_RejectedOnInitSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "init", "--raw")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --raw used on init subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestRawFlag_ShorthandRejectedOnSubcommands(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	// Test -r shorthand on list subcommand
	_, stderr, exitCode := runRecall(t, binPath, recallDir, "list", "-r")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when -r used on list subcommand")
	}
	if !strings.Contains(stderr, "unknown shorthand flag") {
		t.Errorf("expected 'unknown shorthand flag' error, got stderr: %q", stderr)
	}
}

func TestLsAliasListsFiles(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "ls")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	expected := "hello\nplain\n"
	if stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout)
	}
}

func TestVersionCommandShowsMetadata(t *testing.T) {
	binPath := buildBinaryWithMetadata(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "version")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	for _, want := range []string{
		"recall version 0.1.0",
		"Go version:",
		"Platform:",
		"Environment:    test",
		"Branch:         test-branch",
		"Module:",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, stdout)
		}
	}
}

func TestSearchFlag_MatchesSubcommandOutput(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	subStdout, _, subExit := runRecall(t, binPath, recallDir, "search", "Hello")
	flagStdout, _, flagExit := runRecall(t, binPath, recallDir, "-s", "Hello")

	if subExit != 0 || flagExit != 0 {
		t.Errorf("expected exit code 0 for both, got subcommand=%d flag=%d", subExit, flagExit)
	}
	if flagStdout != subStdout {
		t.Errorf("expected -s output to equal search subcommand output.\nsubcommand: %q\nflag:       %q", subStdout, flagStdout)
	}
	// Sanity: the search should have actually matched the hello file.
	if !strings.Contains(flagStdout, "hello:") {
		t.Errorf("expected a match in the 'hello' file, got: %q", flagStdout)
	}
}

func TestSearchFlag_LongFormMatchesSubcommandOutput(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	subStdout, _, _ := runRecall(t, binPath, recallDir, "search", "Hello")
	flagStdout, _, exitCode := runRecall(t, binPath, recallDir, "--search", "Hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if flagStdout != subStdout {
		t.Errorf("expected --search output to equal search subcommand output.\nsubcommand: %q\nflag:       %q", subStdout, flagStdout)
	}
}

func TestSearchFlag_NoMatchesExitsZeroSilently(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-s", "zzznomatchzzz")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no stdout, got: %q", stdout)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got: %q", stderr)
	}
}

func TestSearchFlag_WithoutQueryShowsHelp(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "-s")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "recall") {
		t.Errorf("expected help output containing 'recall', got: %q", stdout)
	}
}

func TestSearchFlag_WithEditFlagProducesError(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "-s", "--edit", "hello")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "--search and --edit") {
		t.Errorf("expected error about mutually exclusive flags, got stderr: %q", stderr)
	}
}

func TestSearchFlag_WithRawFlagIsAllowed(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-s", "-r", "Hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when combining -s and -r, got %d (stderr: %q)", exitCode, stderr)
	}
	if !strings.Contains(stdout, "hello:") {
		t.Errorf("expected a match in the 'hello' file, got: %q", stdout)
	}
}

func TestSearchFlag_ShownInHelp(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--help")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "--search") || !strings.Contains(stdout, "-s,") {
		t.Errorf("expected help to list the --search / -s flag, got: %q", stdout)
	}
}

func TestSearchFlag_RejectedOnEditSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "edit", "--search", "hello")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --search used on edit subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestSearchFlag_RejectedOnListSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "list", "--search")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --search used on list subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestSearchFlag_RejectedOnSearchSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "search", "--search", "hello")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --search used on search subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestSearchFlag_RejectedOnInitSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "init", "--search")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when --search used on init subcommand")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' error, got stderr: %q", stderr)
	}
}

func TestSearchFlag_ShorthandRejectedOnListSubcommand(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "list", "-s")

	if exitCode == 0 {
		t.Error("expected non-zero exit code when -s used on list subcommand")
	}
	if !strings.Contains(stderr, "unknown shorthand flag") {
		t.Errorf("expected 'unknown shorthand flag' error, got stderr: %q", stderr)
	}
}
