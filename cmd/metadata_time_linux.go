//go:build linux

package cmd

import (
	"os"
	"syscall"
	"time"
)

// platformTimes returns the inode-change (ctime) time for the file described by
// info. On Linux, birth (creation) time is not portably exposed via
// syscall.Stat_t, so created is always nil. changed may be nil if the
// underlying Sys() value is not the expected type.
func platformTimes(info os.FileInfo) (created, changed *time.Time) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return nil, nil
	}

	c := time.Unix(st.Ctim.Sec, st.Ctim.Nsec)
	return nil, &c
}
