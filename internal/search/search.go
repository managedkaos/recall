package search

import (
	"bytes"
	"strings"

	"github.com/mjenkins/recall/internal/storage"
)

// Result represents a single search match.
type Result struct {
	Filename string
	LineNum  int    // 1-based
	Line     string // Full content of the matching line
}

// FileResults groups results by file.
type FileResults struct {
	Filename string
	Matches  []Result
}

// Search performs a case-insensitive substring scan of all files in the
// given directory. Returns results grouped by file, with each file's
// matches in ascending line-number order.
func Search(dir string, query string) ([]FileResults, error) {
	names, err := storage.List(dir)
	if err != nil {
		return nil, err
	}

	var results []FileResults
	for _, name := range names {
		content, err := storage.Read(dir, name)
		if err != nil {
			continue
		}
		matches := SearchContent(name, content, query)
		if len(matches) > 0 {
			results = append(results, FileResults{
				Filename: name,
				Matches:  matches,
			})
		}
	}
	return results, nil
}

// SearchContent performs a case-insensitive search within a single file's
// content. Returns matches in ascending line-number order.
func SearchContent(filename string, content []byte, query string) []Result {
	if query == "" {
		return nil
	}

	lowerQuery := strings.ToLower(query)
	lines := bytes.Split(content, []byte("\n"))

	var results []Result
	for i, line := range lines {
		if strings.Contains(strings.ToLower(string(line)), lowerQuery) {
			results = append(results, Result{
				Filename: filename,
				LineNum:  i + 1, // 1-based line numbers
				Line:     string(line),
			})
		}
	}
	return results
}
