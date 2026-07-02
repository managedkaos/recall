package renderer

import (
	"github.com/charmbracelet/glamour"
)

// Render takes markdown content (with front-matter already stripped)
// and returns terminal-formatted output using Glamour.
// Uses the "auto" style which adapts to the terminal's color capabilities.
func Render(content []byte) (string, error) {
	if len(content) == 0 {
		return "", nil
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return "", err
	}

	out, err := r.Render(string(content))
	if err != nil {
		return "", err
	}

	return out, nil
}
