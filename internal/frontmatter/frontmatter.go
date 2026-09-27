package frontmatter

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// delimiter is the line that opens and closes a YAML front matter block.
const delimiter = "---"

// meta is the data model for the supported front matter fields. Additional
// fields can be added here as the format grows.
type meta struct {
	Tags []string `yaml:"tags"`
}

// Parse extracts tags from a YAML front matter block at the top of the content
// and returns the parsed tags along with the body (everything after the closing
// delimiter).
//
// A front matter block must begin on the very first line with a line that is
// exactly "---" and end with a subsequent line that is exactly "---". The text
// between the delimiters is parsed as YAML. Both inline (`tags: [a, b]`) and
// block (`tags:\n  - a`) list forms are supported.
//
// If no valid front matter block is present — including the case of an
// unterminated opening delimiter — Parse returns nil tags and the original
// content unchanged. Empty tags or a missing `tags` key yield nil tags.
func Parse(content []byte) (tags []string, body []byte) {
	if len(content) == 0 {
		return nil, content
	}

	// The opening delimiter must be the very first line. Determine the end of
	// the first line, tolerating a trailing carriage return (CRLF).
	firstLineEnd := bytes.IndexByte(content, '\n')
	var firstLine string
	if firstLineEnd == -1 {
		firstLine = string(content)
	} else {
		firstLine = string(content[:firstLineEnd])
	}
	if strings.TrimRight(firstLine, "\r") != delimiter {
		return nil, content
	}

	// A lone delimiter with no newline (and thus no closing delimiter) is not a
	// valid block.
	if firstLineEnd == -1 {
		return nil, content
	}

	// Scan subsequent lines for the closing delimiter.
	rest := content[firstLineEnd+1:]
	offset := 0
	for {
		lineEnd := bytes.IndexByte(rest[offset:], '\n')
		var line []byte
		var nextOffset int
		if lineEnd == -1 {
			line = rest[offset:]
			nextOffset = len(rest)
		} else {
			line = rest[offset : offset+lineEnd]
			nextOffset = offset + lineEnd + 1
		}

		if strings.TrimRight(string(line), "\r") == delimiter {
			// Found the closing delimiter. YAML is everything before this line;
			// body is everything after it.
			yamlBytes := rest[:offset]
			body = rest[nextOffset:]

			var m meta
			if err := yaml.Unmarshal(yamlBytes, &m); err != nil {
				// Malformed YAML: treat as no usable tags but still strip the
				// block so the body renders cleanly.
				return nil, body
			}
			return normalizeTags(m.Tags), body
		}

		if lineEnd == -1 {
			// Reached the end without a closing delimiter: not a valid block.
			return nil, content
		}
		offset = nextOffset
	}
}

// normalizeTags trims whitespace from each tag and drops empty entries.
// It returns nil when no non-empty tags remain.
func normalizeTags(tags []string) []string {
	var result []string
	for _, t := range tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
