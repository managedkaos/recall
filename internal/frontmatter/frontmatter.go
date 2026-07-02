package frontmatter

import (
	"bytes"
	"strings"
)

const tagsPrefix = "tags:"

// Parse extracts tags from the first line of content if it matches the
// "tags: tag1, tag2, tag3" format. Returns the parsed tags (trimmed,
// non-empty) and the remaining content with the front-matter line removed.
// If no front-matter is found, returns nil tags and the original content.
func Parse(content []byte) (tags []string, body []byte) {
	if len(content) == 0 {
		return nil, content
	}

	// Find the first line
	idx := bytes.IndexByte(content, '\n')
	var firstLine string
	if idx == -1 {
		// Single line, no newline
		firstLine = string(content)
	} else {
		firstLine = string(content[:idx])
	}

	// Check if the first line starts with the tags prefix (case-sensitive)
	if !strings.HasPrefix(firstLine, tagsPrefix) {
		return nil, content
	}

	// Parse the tag values from the first line
	tagsPart := firstLine[len(tagsPrefix):]
	tags = ParseTagLine(tagsPart)

	// Compute the body (everything after the first line)
	if idx == -1 {
		// The entire content was the tags line with no newline
		body = []byte{}
	} else {
		body = content[idx+1:]
	}

	return tags, body
}

// ParseTagLine parses a single tags line string into a slice of trimmed,
// non-empty tag values. Handles consecutive commas, trailing commas,
// and whitespace around values.
func ParseTagLine(line string) []string {
	parts := strings.Split(line, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
