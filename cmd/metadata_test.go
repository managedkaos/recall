package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeTempFile creates a file with the given content in a fresh temp dir and
// returns that directory. The file is created as dir/name.
func writeTempFile(t *testing.T, name, content string) (dir string) {
	t.Helper()
	dir = t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestCollectMetadata_TaggedFile(t *testing.T) {
	dir := writeTempFile(t, "docker", "tags: docker, devops\n# Docker\nbody\n")

	md, err := collectMetadata(dir, "docker")
	if err != nil {
		t.Fatalf("collectMetadata returned error: %v", err)
	}

	if md.Name != "docker" {
		t.Errorf("expected Name 'docker', got %q", md.Name)
	}
	if !filepath.IsAbs(md.Path) {
		t.Errorf("expected absolute Path, got %q", md.Path)
	}
	wantSize := int64(len("tags: docker, devops\n# Docker\nbody\n"))
	if md.SizeBytes != wantSize {
		t.Errorf("expected SizeBytes %d, got %d", wantSize, md.SizeBytes)
	}
	if md.Modified.IsZero() {
		t.Errorf("expected non-zero Modified time")
	}
	if len(md.Tags) != 2 || md.Tags[0] != "docker" || md.Tags[1] != "devops" {
		t.Errorf("expected tags [docker devops], got %v", md.Tags)
	}
}

func TestCollectMetadata_TaglessFileHasEmptyNonNilTags(t *testing.T) {
	dir := writeTempFile(t, "notes", "# Notes\nno tags here\n")

	md, err := collectMetadata(dir, "notes")
	if err != nil {
		t.Fatalf("collectMetadata returned error: %v", err)
	}

	if md.Tags == nil {
		t.Errorf("expected non-nil Tags slice, got nil")
	}
	if len(md.Tags) != 0 {
		t.Errorf("expected empty Tags, got %v", md.Tags)
	}
}

func TestCollectMetadata_MissingFileErrors(t *testing.T) {
	dir := t.TempDir()

	if _, err := collectMetadata(dir, "does-not-exist"); err == nil {
		t.Errorf("expected error for missing file, got nil")
	}
}

func TestRenderJSON_EmptySliceEmitsEmptyArray(t *testing.T) {
	out := captureStdout(t, func() {
		renderJSON(nil)
	})
	if out != "[]\n" {
		t.Errorf("expected %q, got %q", "[]\n", out)
	}
}

func TestRenderJSON_PopulatedRoundTrips(t *testing.T) {
	dir := writeTempFile(t, "docker", "tags: docker\n# Docker\n")
	md, err := collectMetadata(dir, "docker")
	if err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		renderJSON([]FileMetadata{md})
	})

	var decoded []FileMetadata
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 element, got %d", len(decoded))
	}
	if decoded[0].Name != "docker" {
		t.Errorf("expected Name 'docker', got %q", decoded[0].Name)
	}
	if decoded[0].SizeBytes != md.SizeBytes {
		t.Errorf("expected SizeBytes %d, got %d", md.SizeBytes, decoded[0].SizeBytes)
	}
	if len(decoded[0].Tags) != 1 || decoded[0].Tags[0] != "docker" {
		t.Errorf("expected tags [docker], got %v", decoded[0].Tags)
	}
}

func TestPlatformTimes_DoesNotPanicAndIsPlausible(t *testing.T) {
	dir := writeTempFile(t, "f", "content\n")
	info, err := os.Stat(filepath.Join(dir, "f"))
	if err != nil {
		t.Fatal(err)
	}

	created := platformTimes(info)

	// Non-nil returns should be plausible (non-zero) times.
	if created != nil && created.IsZero() {
		t.Errorf("expected non-zero created time when present")
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns what was
// written.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	os.Stdout = w

	done := make(chan string)
	go func() {
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 1024)
		for {
			n, err := r.Read(tmp)
			buf = append(buf, tmp[:n]...)
			if err != nil {
				break
			}
		}
		done <- string(buf)
	}()

	fn()
	w.Close()
	os.Stdout = orig
	return <-done
}
