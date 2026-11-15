//go:build linux

package node

import (
	"errors"
	"os"
	"syscall"
)

// lockRun holds an exclusive lock beside the database of this controller for
// the life of this process, and answers errRunning while another holds it.
// The kernel lets go of the lock when the process ends, however it ends, so no
// file left behind by a crash can lock a server out.
func lockRun(db string) (func(), error) {
	f, err := os.OpenFile(db+".lock", os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, errRunning
		}
		return nil, err
	}
	return func() { f.Close() }, nil
}
