//go:build darwin

package cmd

import (
	"os"
	"syscall"
	"time"
)

// platformTimes returns the creation (birth) and inode-change (ctime) times for
// the file described by info, when the platform exposes them. On macOS both are
// available via syscall.Stat_t. Either return value may be nil if the
// underlying Sys() value is not the expected type.
func platformTimes(info os.FileInfo) (created, changed *time.Time) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return nil, nil
	}

	b := time.Unix(st.Birthtimespec.Sec, st.Birthtimespec.Nsec)
	c := time.Unix(st.Ctimespec.Sec, st.Ctimespec.Nsec)
	return &b, &c
}
