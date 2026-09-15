//go:build darwin

package cmd

import (
	"os"
	"syscall"
	"time"
)

// platformTimes returns the creation (birth) time for the file described by
// info, when the platform exposes it. On macOS it is available via
// syscall.Stat_t. The return value may be nil if the underlying Sys() value is
// not the expected type.
func platformTimes(info os.FileInfo) (created *time.Time) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return nil
	}

	b := time.Unix(st.Birthtimespec.Sec, st.Birthtimespec.Nsec)
	return &b
}
