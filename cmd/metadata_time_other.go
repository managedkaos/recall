//go:build !darwin && !linux && !windows

package cmd

import (
	"os"
	"time"
)

// platformTimes is the fallback for platforms that do not expose creation or
// inode-change times through a supported syscall structure. Both return values
// are always nil.
func platformTimes(info os.FileInfo) (created, changed *time.Time) {
	return nil, nil
}
