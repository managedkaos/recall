package renderer

import (
	"strings"
	"testing"
)

func TestRender_EmptyContent(t *testing.T) {
	out, err := Render([]byte{})
	if err != nil {
		t.Fatalf("unexpected error for empty content: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string for empty content, got %q", out)
	}
}

func TestRender_NilContent(t *testing.T) {
	out, err := Render(nil)
	if err != nil {
		t.Fatalf("unexpected error for nil content: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string for nil content, got %q", out)
	}
}

func TestRender_SimpleHeading(t *testing.T) {
	input := []byte("# Hello World\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Hello World") {
		t.Errorf("expected output to contain 'Hello World', got %q", out)
	}
}

func TestRender_BoldText(t *testing.T) {
	input := []byte("This is **bold** text.\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "bold") {
		t.Errorf("expected output to contain 'bold', got %q", out)
	}
}

func TestRender_CodeBlock(t *testing.T) {
	input := []byte("```go\nfmt.Println(\"hello\")\n```\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Println") {
		t.Errorf("expected output to contain 'Println', got %q", out)
	}
}

func TestRender_UnorderedList(t *testing.T) {
	input := []byte("- item one\n- item two\n- item three\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "item one") {
		t.Errorf("expected output to contain 'item one', got %q", out)
	}
	if !strings.Contains(out, "item two") {
		t.Errorf("expected output to contain 'item two', got %q", out)
	}
}

func TestRender_OrderedList(t *testing.T) {
	input := []byte("1. first\n2. second\n3. third\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "first") {
		t.Errorf("expected output to contain 'first', got %q", out)
	}
	if !strings.Contains(out, "second") {
		t.Errorf("expected output to contain 'second', got %q", out)
	}
}

func TestRender_MultipleElements(t *testing.T) {
	input := []byte("# Title\n\nSome paragraph with **bold** and *italic*.\n\n- list item\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Title") {
		t.Errorf("expected output to contain 'Title', got %q", out)
	}
	if !strings.Contains(out, "bold") {
		t.Errorf("expected output to contain 'bold', got %q", out)
	}
	if !strings.Contains(out, "list item") {
		t.Errorf("expected output to contain 'list item', got %q", out)
	}
}

func TestRender_ProducesOutput(t *testing.T) {
	// Even plain text should produce some output
	input := []byte("Just a simple paragraph.\n")
	out, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(strings.TrimSpace(out)) == 0 {
		t.Error("expected non-empty output for valid markdown content")
	}
}
