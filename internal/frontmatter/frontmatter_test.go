package frontmatter

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestMetaUnmarshal verifies the data model unmarshals both inline and block
// list forms of the tags field.
func TestMetaUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want []string
	}{
		{
			name: "inline list",
			yaml: "tags: [go, testing]",
			want: []string{"go", "testing"},
		},
		{
			name: "block list",
			yaml: "tags:\n  - go\n  - testing",
			want: []string{"go", "testing"},
		},
		{
			name: "no tags key",
			yaml: "other: value",
			want: nil,
		},
		{
			name: "empty inline list",
			yaml: "tags: []",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m meta
			if err := yaml.Unmarshal([]byte(tt.yaml), &m); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			// Compare via normalizeTags to treat empty/nil uniformly; the model
			// itself may produce an empty non-nil slice for `tags: []`.
			if !slicesEqual(normalizeTags(m.Tags), tt.want) {
				t.Errorf("Tags = %v, want %v", m.Tags, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		wantTags []string
		wantBody string
	}{
		{
			name:     "inline list front matter",
			content:  []byte("---\ntags: [go, testing]\n---\n# Hello\nWorld"),
			wantTags: []string{"go", "testing"},
			wantBody: "# Hello\nWorld",
		},
		{
			name:     "block list front matter",
			content:  []byte("---\ntags:\n  - go\n  - testing\n---\n# Hello"),
			wantTags: []string{"go", "testing"},
			wantBody: "# Hello",
		},
		{
			name:     "no front matter",
			content:  []byte("# Hello\nWorld"),
			wantTags: nil,
			wantBody: "# Hello\nWorld",
		},
		{
			name:     "empty content",
			content:  []byte{},
			wantTags: nil,
			wantBody: "",
		},
		{
			name:     "nil content",
			content:  nil,
			wantTags: nil,
			wantBody: "",
		},
		{
			name:     "empty front matter block",
			content:  []byte("---\n---\n# Body"),
			wantTags: nil,
			wantBody: "# Body",
		},
		{
			name:     "front matter with no tags key",
			content:  []byte("---\ntitle: Something\n---\nBody"),
			wantTags: nil,
			wantBody: "Body",
		},
		{
			name:     "empty tags list",
			content:  []byte("---\ntags: []\n---\nBody"),
			wantTags: nil,
			wantBody: "Body",
		},
		{
			name:     "unterminated block returns original content",
			content:  []byte("---\ntags: [go]\n# never closed"),
			wantTags: nil,
			wantBody: "---\ntags: [go]\n# never closed",
		},
		{
			name:     "lone delimiter no newline",
			content:  []byte("---"),
			wantTags: nil,
			wantBody: "---",
		},
		{
			name:     "body preserves multiple lines",
			content:  []byte("---\ntags: [go]\n---\nline1\nline2\nline3"),
			wantTags: []string{"go"},
			wantBody: "line1\nline2\nline3",
		},
		{
			name:     "closing delimiter at end no trailing body",
			content:  []byte("---\ntags: [go]\n---\n"),
			wantTags: []string{"go"},
			wantBody: "",
		},
		{
			name:     "closing delimiter at end no newline",
			content:  []byte("---\ntags: [go]\n---"),
			wantTags: []string{"go"},
			wantBody: "",
		},
		{
			name:     "clean break: old leading tags line is body",
			content:  []byte("tags: go, testing\n# Hello"),
			wantTags: nil,
			wantBody: "tags: go, testing\n# Hello",
		},
		{
			name:     "CRLF line endings",
			content:  []byte("---\r\ntags: [go, testing]\r\n---\r\n# Hello\r\nWorld"),
			wantTags: []string{"go", "testing"},
			wantBody: "# Hello\r\nWorld",
		},
		{
			name:     "tags with surrounding whitespace",
			content:  []byte("---\ntags: [\" go \", \"testing \"]\n---\nBody"),
			wantTags: []string{"go", "testing"},
			wantBody: "Body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTags, gotBody := Parse(tt.content)
			if !slicesEqual(gotTags, tt.wantTags) {
				t.Errorf("Parse() tags = %v, want %v", gotTags, tt.wantTags)
			}
			if string(gotBody) != tt.wantBody {
				t.Errorf("Parse() body = %q, want %q", string(gotBody), tt.wantBody)
			}
		})
	}
}

// slicesEqual compares two string slices for equality.
// nil and empty slices are treated differently: nil != []string{}.
func slicesEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
