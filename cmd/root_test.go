package cmd_test

import (
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
