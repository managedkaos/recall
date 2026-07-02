package cmd

import (
	"testing"
)

func TestFormatVersion_AllPopulated(t *testing.T) {
	got := formatVersion("0", "1", "0")
	if got != "0.1.0" {
		t.Errorf("expected '0.1.0', got %q", got)
	}
}

func TestFormatVersion_LargeNumbers(t *testing.T) {
	got := formatVersion("12", "34", "56")
	if got != "12.34.56" {
		t.Errorf("expected '12.34.56', got %q", got)
	}
}

func TestFormatVersion_EmptyMajor(t *testing.T) {
	got := formatVersion("", "1", "0")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_EmptyMinor(t *testing.T) {
	got := formatVersion("0", "", "0")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_EmptyPatch(t *testing.T) {
	got := formatVersion("0", "1", "")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestFormatVersion_AllEmpty(t *testing.T) {
	got := formatVersion("", "", "")
	if got != "unknown" {
		t.Errorf("expected 'unknown', got %q", got)
	}
}

func TestVersionReservedName(t *testing.T) {
	if !IsReservedName("version") {
		t.Error("expected 'version' to be a reserved name")
	}
}
