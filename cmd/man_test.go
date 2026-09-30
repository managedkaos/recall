package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMan(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "man")
	if err := GenerateMan(dir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "recall.1"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{".TH", ".SH NAME", ".SH SYNOPSIS", ".SH OPTIONS", "--tag", "--metadata", "--help"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("manual missing %q", want)
		}
	}
}
