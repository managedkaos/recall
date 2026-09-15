//go:build linux

package cmd

import (
	"os"
	"time"
)

// platformTimes returns creation time for the file described by info. On Linux,
// birth (creation) time is not portably exposed via syscall.Stat_t, so created
// is always nil.
func platformTimes(info os.FileInfo) (created *time.Time) {
	return nil
}
