//go:build windows

package cmd

import (
	"os"
	"syscall"
	"time"
)

// platformTimes returns the creation time for the file described by info. On
// Windows, creation time is available via syscall.Win32FileAttributeData. The
// return value may be nil if the underlying Sys() value is not the expected
// type.
func platformTimes(info os.FileInfo) (created *time.Time) {
	d, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok || d == nil {
		return nil
	}

	b := time.Unix(0, d.CreationTime.Nanoseconds())
	return &b
}
