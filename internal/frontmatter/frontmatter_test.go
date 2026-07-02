package frontmatter

import (
	"testing"
)

func TestParseTagLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "simple tags",
			line: " go, testing, concurrency",
			want: []string{"go", "testing", "concurrency"},
		},
		{
			name: "consecutive commas",
			line: " go,, testing",
			want: []string{"go", "testing"},
		},
		{
			name: "trailing comma",
			line: " go, testing,",
			want: []string{"go", "testing"},
		},
		{
			name: "leading comma",
			line: ", go, testing",
			want: []string{"go", "testing"},
		},
		{
			name: "whitespace-only values",
			line: " go,   , testing",
			want: []string{"go", "testing"},
		},
		{
			name: "empty string",
			line: "",
			want: nil,
		},
		{
			name: "only whitespace",
			line: "   ",
			want: nil,
		},
		{
			name: "only commas",
			line: ",,,",
			want: nil,
		},
		{
			name: "single tag with spaces",
			line: "  golang  ",
			want: []string{"golang"},
		},
		{
			name: "tags with extra internal spaces",
			line: " my tag , another tag ",
			want: []string{"my tag", "another tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTagLine(tt.line)
			if !slicesEqual(got, tt.want) {
				t.Errorf("ParseTagLine(%q) = %v, want %v", tt.line, got, tt.want)
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
			name:     "standard front-matter",
			content:  []byte("tags: go, testing\n# Hello\nWorld"),
			wantTags: []string{"go", "testing"},
			wantBody: "# Hello\nWorld",
		},
		{
			name:     "no front-matter",
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
			name:     "tags line only no newline",
			content:  []byte("tags: go, testing"),
			wantTags: []string{"go", "testing"},
			wantBody: "",
		},
		{
			name:     "tags line only with newline",
			content:  []byte("tags: go, testing\n"),
			wantTags: []string{"go", "testing"},
			wantBody: "",
		},
		{
			name:     "case-sensitive prefix - Tags does not match",
			content:  []byte("Tags: go, testing\n# Hello"),
			wantTags: nil,
			wantBody: "Tags: go, testing\n# Hello",
		},
		{
			name:     "case-sensitive prefix - TAGS does not match",
			content:  []byte("TAGS: go, testing\n# Hello"),
			wantTags: nil,
			wantBody: "TAGS: go, testing\n# Hello",
		},
		{
			name:     "consecutive commas in front-matter",
			content:  []byte("tags: go,, testing\nBody"),
			wantTags: []string{"go", "testing"},
			wantBody: "Body",
		},
		{
			name:     "trailing comma in front-matter",
			content:  []byte("tags: go, testing,\nBody"),
			wantTags: []string{"go", "testing"},
			wantBody: "Body",
		},
		{
			name:     "whitespace-only tag values",
			content:  []byte("tags: go,   , testing\nBody"),
			wantTags: []string{"go", "testing"},
			wantBody: "Body",
		},
		{
			name:     "tags prefix without space",
			content:  []byte("tags:go, testing\nBody"),
			wantTags: []string{"go", "testing"},
			wantBody: "Body",
		},
		{
			name:     "line starting with tags but not prefix",
			content:  []byte("tagster: something\nBody"),
			wantTags: nil,
			wantBody: "tagster: something\nBody",
		},
		{
			name:     "body preserves multiple lines",
			content:  []byte("tags: go\nline1\nline2\nline3"),
			wantTags: []string{"go"},
			wantBody: "line1\nline2\nline3",
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
