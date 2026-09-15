//go:build !darwin && !linux && !windows

package cmd

import (
	"os"
	"time"
)

// platformTimes is the fallback for platforms that do not expose creation time
// through a supported syscall structure. The return value is always nil.
func platformTimes(info os.FileInfo) (created *time.Time) {
	return nil
}
