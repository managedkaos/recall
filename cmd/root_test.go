package cmd_test

import (
	"encoding/json"
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
		"-X %s/cmd.Version=0.1.0 -X %s/cmd.GitBranch=test-branch -X %s/cmd.BuildDate=2026-07-21T12:00:00Z -X %s/cmd.BuildEnvironment=test",
		module, module, module, module,
	)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", binPath, ".")
	cmd.Dir = filepath.Join(getProjectRoot(t))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary with metadata: %v\n%s", err, out)
	}
	return binPath
}

// buildBinaryWithVersion compiles recall injecting a single cmd.Version ldflag
// (as GoReleaser does from the git tag).
func buildBinaryWithVersion(t *testing.T, version string) string {
	t.Helper()
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "recall")
	module := "github.com/managedkaos/recall"
	ldflags := fmt.Sprintf(
		"-X %s/cmd.Version=%s -X %s/cmd.GitBranch=test-branch -X %s/cmd.BuildDate=2026-07-21T12:00:00Z -X %s/cmd.BuildEnvironment=goreleaser",
		module, version, module, module, module,
	)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", binPath, ".")
	cmd.Dir = filepath.Join(getProjectRoot(t))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary with version: %v\n%s", err, out)
	}
	return binPath
}

// getProjectRoot returns the project root directory.
func getProjectRoot(t *testing.T) string {
	t.Helper()
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

// --- Default render path -----------------------------------------------------

func TestRender_ValidFile(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Hello, world!") {
		t.Errorf("expected rendered output containing 'Hello, world!', got %q", stdout)
	}
}

func TestRender_NonExistentFileExits1Silently(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "nonexistent")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("expected no output, got stdout %q stderr %q", stdout, stderr)
	}
}

func TestNoArgsShowsHelp(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "recall") {
		t.Errorf("expected help output containing 'recall', got: %q", stdout)
	}
}

// Former subcommand names are now valid filenames, not commands.
func TestFormerSubcommandNameTreatedAsFilename(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	for _, name := range []string{"list", "search", "edit", "init", "version"} {
		// No such file exists, so this should exit 1 silently (render path),
		// not invoke a subcommand.
		stdout, stderr, exitCode := runRecall(t, binPath, recallDir, name)
		if exitCode != 1 {
			t.Errorf("%q: expected exit 1 (treated as missing file), got %d", name, exitCode)
		}
		if stdout != "" || stderr != "" {
			t.Errorf("%q: expected no output, got stdout %q stderr %q", name, stdout, stderr)
		}
	}

	// A file literally named "list" should render.
	if err := os.WriteFile(filepath.Join(recallDir, "list"), []byte("I am a file named list\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, _, exitCode := runRecall(t, binPath, recallDir, "list")
	if exitCode != 0 {
		t.Errorf("expected exit 0 rendering file named 'list', got %d", exitCode)
	}
	if !strings.Contains(stdout, "I am a file named list") {
		t.Errorf("expected rendered content of file 'list', got %q", stdout)
	}
}

// --- Raw flag ----------------------------------------------------------------

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

// --- Edit flag ---------------------------------------------------------------

func TestEditFlag_NoArgErrors(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "--edit")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "--edit requires a filename") {
		t.Errorf("expected 'requires a filename' error, got stderr: %q", stderr)
	}
}

func TestEditFlag_CreatesAndOpensFile(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	// Use `true` as a no-op editor so the command completes without interaction.
	cmd := exec.Command(binPath, "--edit", "newnote")
	cmd.Env = append(os.Environ(), "RECALL_DIR="+recallDir, "EDITOR=true")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("expected --edit to succeed, got err %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(recallDir, "newnote")); err != nil {
		t.Errorf("expected file 'newnote' to be created, stat err: %v", err)
	}
}

// --- Search flag -------------------------------------------------------------

func TestSearchFlag_ShortAndLongEqual(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	shortOut, _, shortExit := runRecall(t, binPath, recallDir, "-s", "Hello")
	longOut, _, longExit := runRecall(t, binPath, recallDir, "--search", "Hello")

	if shortExit != 0 || longExit != 0 {
		t.Errorf("expected exit 0 for both, got -s=%d --search=%d", shortExit, longExit)
	}
	if shortOut != longOut {
		t.Errorf("expected -s and --search output to match.\n-s:       %q\n--search: %q", shortOut, longOut)
	}
	if !strings.Contains(shortOut, "hello:") {
		t.Errorf("expected a match in the 'hello' file, got: %q", shortOut)
	}
}

func TestSearchFlag_NoMatchesExitsZeroSilently(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-s", "zzznomatchzzz")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("expected no output, got stdout %q stderr %q", stdout, stderr)
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

// --- List flag ---------------------------------------------------------------

func TestListFlag_ListsFiles(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	for _, args := range [][]string{{"--list"}, {"-l"}} {
		stdout, _, exitCode := runRecall(t, binPath, recallDir, args...)
		if exitCode != 0 {
			t.Errorf("%v: expected exit code 0, got %d", args, exitCode)
		}
		expected := "hello\nplain\n"
		if stdout != expected {
			t.Errorf("%v: expected stdout %q, got %q", args, expected, stdout)
		}
	}
}

func TestListFlag_FilterByTag(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--list", "--tag", "greeting")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	expected := "hello\n"
	if stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout)
	}
}

// --- Init flag ---------------------------------------------------------------

func TestInitFlag_DefaultDir(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := filepath.Join(t.TempDir(), "created-by-init")

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--init")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "initialized directory at "+recallDir) {
		t.Errorf("expected init confirmation for %q, got: %q", recallDir, stdout)
	}
	if _, err := os.Stat(recallDir); err != nil {
		t.Errorf("expected directory %q to exist, stat err: %v", recallDir, err)
	}
}

func TestInitFlag_WithInitPath(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)
	custom := filepath.Join(t.TempDir(), "custom-init")

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--init", "--init-path", custom)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "initialized directory at "+custom) {
		t.Errorf("expected init confirmation for %q, got: %q", custom, stdout)
	}
	if _, err := os.Stat(custom); err != nil {
		t.Errorf("expected directory %q to exist, stat err: %v", custom, err)
	}
}

func TestInitPath_WithoutInitErrors(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "--init-path", "/tmp/whatever")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "--init-path requires --init") {
		t.Errorf("expected '--init-path requires --init' error, got stderr: %q", stderr)
	}
}

// --- Version flag ------------------------------------------------------------

func TestVersionFlag_ShowsMetadata(t *testing.T) {
	binPath := buildBinaryWithMetadata(t)
	recallDir := setupRecallDir(t)

	for _, args := range [][]string{{"--version"}, {"-v"}} {
		stdout, _, exitCode := runRecall(t, binPath, recallDir, args...)
		if exitCode != 0 {
			t.Errorf("%v: expected exit code 0, got %d", args, exitCode)
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
				t.Errorf("%v: expected output to contain %q, got:\n%s", args, want, stdout)
			}
		}
	}
}

func TestVersionFlag_PlainBuildIsUnknown(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)
	for _, flag := range []string{"--version", "-v"} {
		stdout, _, exitCode := runRecall(t, binPath, recallDir, flag)
		if exitCode != 0 || !strings.HasPrefix(stdout, "recall version unknown\n") {
			t.Errorf("%s: expected unknown version and exit 0, got exit %d:\n%s", flag, exitCode, stdout)
		}
	}
}

// A single injected cmd.Version (as GoReleaser sets from the tag) must be shown
// after normalization.
func TestVersionFlag_InjectedVersion(t *testing.T) {
	binPath := buildBinaryWithVersion(t, "9.9.9")
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--version")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "recall version 9.9.9") {
		t.Errorf("expected 'recall version 9.9.9', got:\n%s", stdout)
	}
}

// A tag with a leading 'v' must be normalized to a bare semantic version.
func TestVersionFlag_InjectedVersionStripsLeadingV(t *testing.T) {
	binPath := buildBinaryWithVersion(t, "v2.3.4")
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--version")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "recall version 2.3.4") {
		t.Errorf("expected 'recall version 2.3.4', got:\n%s", stdout)
	}
}

// --- Completion flag ---------------------------------------------------------

func TestCompletionFlag_SupportedShells(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-c", shell)
		if exitCode != 0 {
			t.Errorf("%s: expected exit code 0, got %d (stderr: %q)", shell, exitCode, stderr)
		}
		if len(strings.TrimSpace(stdout)) == 0 {
			t.Errorf("%s: expected a non-empty completion script", shell)
		}
	}
}

func TestCompletionFlag_UnsupportedShellErrors(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "-c", "bogus")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "unsupported shell") {
		t.Errorf("expected 'unsupported shell' error, got stderr: %q", stderr)
	}
}

// --- Mutual exclusivity ------------------------------------------------------

func TestActionFlags_MutuallyExclusive(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	pairs := [][]string{
		{"--edit", "--list"},
		{"--search", "--list"},
		{"--version", "--init"},
		{"--list", "--version"},
		{"-e", "-s"},
	}
	for _, args := range pairs {
		_, stderr, exitCode := runRecall(t, binPath, recallDir, args...)
		if exitCode != 1 {
			t.Errorf("%v: expected exit code 1, got %d", args, exitCode)
		}
		if !strings.Contains(stderr, "only one action flag") {
			t.Errorf("%v: expected 'only one action flag' error, got stderr: %q", args, stderr)
		}
	}
}

// --- Help --------------------------------------------------------------------

func TestHelp_ListsAllFlagsNoSubcommands(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "--help")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	for _, want := range []string{
		"-e, --edit", "-s, --search", "-l, --list",
		"-i, --init", "-v, --version", "-c, --completion",
		"-r, --raw", "--tag", "--init-path",
		"-m, --metadata", "-j, --json",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("expected help to list %q, got:\n%s", want, stdout)
		}
	}
	// No subcommands should be advertised.
	if strings.Contains(stdout, "Available Commands:") {
		t.Errorf("expected no subcommands section in help, got:\n%s", stdout)
	}
}

// --- Metadata flag -----------------------------------------------------------

func TestMetadataFlag_ExistingFileTextOutput(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	for _, args := range [][]string{{"-m", "hello"}, {"--metadata", "hello"}} {
		stdout, stderr, exitCode := runRecall(t, binPath, recallDir, args...)
		if exitCode != 0 {
			t.Errorf("%v: expected exit code 0, got %d (stderr: %q)", args, exitCode, stderr)
		}
		for _, want := range []string{
			"Name:     hello",
			"Path:",
			"Size:",
			"Modified:",
			"Tags:     greeting",
		} {
			if !strings.Contains(stdout, want) {
				t.Errorf("%v: expected stdout to contain %q, got:\n%s", args, want, stdout)
			}
		}
		if !strings.Contains(stdout, filepath.Join(recallDir, "hello")) {
			t.Errorf("%v: expected stdout to contain the file path, got:\n%s", args, stdout)
		}
	}
}

func TestMetadataFlag_TaglessFileShowsNone(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "-m", "plain")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Tags:     (none)") {
		t.Errorf("expected 'Tags:     (none)' for tagless file, got:\n%s", stdout)
	}
}

func TestMetadataFlag_JSONOutput(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-m", "-j", "hello")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d (stderr: %q)", exitCode, stderr)
	}

	var decoded []struct {
		Name      string   `json:"name"`
		Path      string   `json:"path"`
		SizeBytes int64    `json:"size_bytes"`
		Modified  string   `json:"modified"`
		Tags      []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
		t.Fatalf("expected valid JSON, got error %v\noutput:\n%s", err, stdout)
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 element, got %d", len(decoded))
	}
	if decoded[0].Name != "hello" {
		t.Errorf("expected name 'hello', got %q", decoded[0].Name)
	}
	if decoded[0].SizeBytes <= 0 {
		t.Errorf("expected positive size, got %d", decoded[0].SizeBytes)
	}
	if decoded[0].Modified == "" {
		t.Errorf("expected non-empty modified timestamp")
	}
	if len(decoded[0].Tags) != 1 || decoded[0].Tags[0] != "greeting" {
		t.Errorf("expected tags [greeting], got %v", decoded[0].Tags)
	}
}

func TestMetadataFlag_MissingFileErrorsExit1(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-m", "nope")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "file not found: nope") {
		t.Errorf("expected 'file not found' error, got stderr: %q", stderr)
	}
}

func TestMetadataFlag_MixedExistingAndMissingExit0(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-m", "hello", "nope")
	if exitCode != 0 {
		t.Errorf("expected exit code 0 (at least one file exists), got %d", exitCode)
	}
	if !strings.Contains(stdout, "Name:     hello") {
		t.Errorf("expected metadata for 'hello', got:\n%s", stdout)
	}
	if !strings.Contains(stderr, "file not found: nope") {
		t.Errorf("expected 'file not found: nope' on stderr, got: %q", stderr)
	}
}

func TestMetadataFlag_JSONMissingEmitsEmptyArrayExit1(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, stderr, exitCode := runRecall(t, binPath, recallDir, "-m", "-j", "nope")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if strings.TrimSpace(stdout) != "[]" {
		t.Errorf("expected stdout '[]', got %q", stdout)
	}
	if !strings.Contains(stderr, "file not found: nope") {
		t.Errorf("expected 'file not found' error, got stderr: %q", stderr)
	}
}

func TestMetadataFlag_NoArgsShowsHelp(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "-m")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "recall") {
		t.Errorf("expected help output containing 'recall', got: %q", stdout)
	}
}

func TestJSONFlag_WithoutMetadataErrors(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	_, stderr, exitCode := runRecall(t, binPath, recallDir, "-j", "hello")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "--json requires --metadata") {
		t.Errorf("expected '--json requires --metadata' error, got stderr: %q", stderr)
	}
}

func TestMetadataFlag_MutuallyExclusiveWithOtherActions(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	pairs := [][]string{
		{"-m", "-l"},
		{"--metadata", "--list", "hello"},
		{"-m", "-e", "hello"},
	}
	for _, args := range pairs {
		_, stderr, exitCode := runRecall(t, binPath, recallDir, args...)
		if exitCode != 1 {
			t.Errorf("%v: expected exit code 1, got %d", args, exitCode)
		}
		if !strings.Contains(stderr, "only one action flag") {
			t.Errorf("%v: expected 'only one action flag' error, got stderr: %q", args, stderr)
		}
	}
}

func TestMetadataFlag_DoesNotMutateFiles(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	before, err := os.ReadFile(filepath.Join(recallDir, "hello"))
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(filepath.Join(recallDir, "hello"))
	if err != nil {
		t.Fatal(err)
	}

	if _, _, exitCode := runRecall(t, binPath, recallDir, "-m", "-j", "hello"); exitCode != 0 {
		t.Fatalf("expected exit 0, got %d", exitCode)
	}

	after, err := os.ReadFile(filepath.Join(recallDir, "hello"))
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(filepath.Join(recallDir, "hello"))
	if err != nil {
		t.Fatal(err)
	}

	if string(before) != string(after) {
		t.Errorf("file content changed after metadata read")
	}
	if !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		t.Errorf("file modtime changed after metadata read")
	}
}
