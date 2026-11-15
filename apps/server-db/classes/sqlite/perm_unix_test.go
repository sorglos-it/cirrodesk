//go:build unix

package sqlite

import (
	"os"
	"syscall"
	"testing"
)

// Under the usual umask 0022 the VFS made every file 0644. The file of an
// older start that is still 0644 is closed as well.
func TestFilesAreOwnerOnly(t *testing.T) {
	old := syscall.Umask(0o022)
	defer syscall.Umask(old)
	cfg := testConfig(t, "state.db")
	if err := os.WriteFile(cfg.Path+"-shm", nil, 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Open(cfg, "test", 1790000000)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for _, p := range []string{cfg.Path, cfg.Path + "-wal", cfg.Path + "-shm"} {
		st, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if m := st.Mode().Perm(); m != 0o600 {
			t.Errorf("%s: mode %o, want 600", p, m)
		}
	}
}
