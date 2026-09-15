package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// completionNames parses the stdout of Cobra's hidden __complete command into
// the list of candidate names, dropping the trailing ":<directive>" line and
// any per-candidate "\tdescription" suffix.
func completionNames(stdout string) []string {
	var names []string
	for _, line := range strings.Split(strings.TrimRight(stdout, "\n"), "\n") {
		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// The final line is the directive marker, e.g. ":4".
		if strings.HasPrefix(trimmed, ":") {
			continue
		}
		// Strip an optional tab-separated description.
		if i := strings.IndexByte(line, '\t'); i >= 0 {
			line = line[:i]
		}
		names = append(names, strings.TrimSpace(line))
	}
	return names
}

// assertNoFileComp fails the test unless the __complete output carries the
// ShellCompDirectiveNoFileComp marker (directive value 4).
func assertNoFileComp(t *testing.T, stdout string) {
	t.Helper()
	if !strings.Contains(stdout, ":4") {
		t.Errorf("expected NoFileComp directive (:4) in output, got:\n%s", stdout)
	}
}

func TestComplete_EmptyPrefixListsAllFiles(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	got := completionNames(stdout)
	want := map[string]bool{"hello": true, "plain": true}
	if len(got) != len(want) {
		t.Fatalf("expected %d candidates, got %d: %v", len(want), len(got), got)
	}
	for _, name := range got {
		if !want[name] {
			t.Errorf("unexpected candidate %q in %v", name, got)
		}
	}
	// NoFileComp directive is value 4.
	if !strings.Contains(stdout, ":4") {
		t.Errorf("expected NoFileComp directive (:4) in output, got:\n%s", stdout)
	}
}

func TestComplete_PrefixFilters(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "he")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	got := completionNames(stdout)
	if len(got) != 1 || got[0] != "hello" {
		t.Errorf("expected only 'hello' for prefix 'he', got %v", got)
	}
}

func TestComplete_NoMatchReturnsNoNames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "zzz")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if got := completionNames(stdout); len(got) != 0 {
		t.Errorf("expected no candidates for prefix 'zzz', got %v", got)
	}
}

func TestComplete_EditFlagCompletesFilenames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "-e", "he")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if got := completionNames(stdout); len(got) != 1 || got[0] != "hello" {
		t.Errorf("expected 'hello' under -e, got %v", got)
	}
}

func TestComplete_MetadataFlagCompletesFilenames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "-m", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	got := completionNames(stdout)
	if len(got) != 2 {
		t.Errorf("expected 2 candidates under -m, got %v", got)
	}
}

func TestComplete_ActionFlagsSuppressFilenames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	for _, flag := range []string{"-l", "-i", "-v"} {
		stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", flag, "")
		if exitCode != 0 {
			t.Fatalf("%s: expected exit code 0, got %d", flag, exitCode)
		}
		if got := completionNames(stdout); len(got) != 0 {
			t.Errorf("%s: expected no filename candidates, got %v", flag, got)
		}
		assertNoFileComp(t, stdout)
	}
}

// A filename that begins with '-' should still be completable. Cobra requires
// callers to pass such positional values after "--"; the completion machinery
// forwards the raw toComplete to the ValidArgsFunction, so prefix matching must
// still work.
func TestComplete_DashPrefixedFilename(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	if err := os.WriteFile(filepath.Join(recallDir, "-dashy"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "--", "-da")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	got := completionNames(stdout)
	if len(got) != 1 || got[0] != "-dashy" {
		t.Errorf("expected only '-dashy' for prefix '-da', got %v", got)
	}
	assertNoFileComp(t, stdout)
}

func TestComplete_SearchQueryNotCompletedAsFilename(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "-s", "he")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if got := completionNames(stdout); len(got) != 0 {
		t.Errorf("expected no filename candidates for search query, got %v", got)
	}
}

func TestComplete_CompletionFlagOffersShellNames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "-c", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	got := completionNames(stdout)
	want := map[string]bool{"bash": true, "zsh": true, "fish": true, "powershell": true}
	if len(got) != len(want) {
		t.Fatalf("expected %d shells, got %d: %v", len(want), len(got), got)
	}
	for _, name := range got {
		if !want[name] {
			t.Errorf("unexpected shell candidate %q in %v", name, got)
		}
	}
	assertNoFileComp(t, stdout)
}

func TestComplete_NonexistentDirDegradesGracefully(t *testing.T) {
	binPath := buildBinary(t)
	// Point RECALL_DIR at a path that does not exist.
	recallDir := filepath.Join(t.TempDir(), "does-not-exist")

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0 even for missing dir, got %d", exitCode)
	}
	if got := completionNames(stdout); len(got) != 0 {
		t.Errorf("expected no candidates for missing dir, got %v", got)
	}
}

func TestComplete_EmptyDirReturnsNoNames(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := t.TempDir() // exists but empty

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if got := completionNames(stdout); len(got) != 0 {
		t.Errorf("expected no candidates for empty dir, got %v", got)
	}
}

func TestComplete_HiddenFilesExcluded(t *testing.T) {
	binPath := buildBinary(t)
	recallDir := setupRecallDir(t)

	// A hidden file must not appear as a completion candidate.
	if err := os.WriteFile(filepath.Join(recallDir, ".secret"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, exitCode := runRecall(t, binPath, recallDir, "__complete", "")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	for _, name := range completionNames(stdout) {
		if strings.HasPrefix(name, ".") {
			t.Errorf("hidden file %q should not be completed", name)
		}
	}
}
