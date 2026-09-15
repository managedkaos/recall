package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestCompleteRecallFiles_Unit exercises the unexported completeRecallFiles
// helper directly, driving it via RECALL_DIR so it does not require a built
// binary. t.Setenv restores the previous value automatically at test end.
func TestCompleteRecallFiles_Unit(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "alpine", "beta"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Hidden files must be excluded.
	if err := os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RECALL_DIR", dir)

	tests := []struct {
		name       string
		toComplete string
		want       []string
	}{
		{"empty prefix lists all", "", []string{"alpha", "alpine", "beta"}},
		{"prefix alp", "alp", []string{"alpha", "alpine"}},
		{"prefix beta", "bet", []string{"beta"}},
		{"no match", "zzz", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, directive := completeRecallFiles(tc.toComplete)
			if directive != cobra.ShellCompDirectiveNoFileComp {
				t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("index %d: expected %q, got %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}

func TestCompleteRecallFiles_MissingDir(t *testing.T) {
	t.Setenv("RECALL_DIR", filepath.Join(t.TempDir(), "nope"))

	got, directive := completeRecallFiles("")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}
	if len(got) != 0 {
		t.Errorf("expected no candidates for missing dir, got %v", got)
	}
}
