package dbcli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fake(t *testing.T, fail string) *[][]string {
	t.Helper()
	var calls [][]string
	old := run
	run = func(_ context.Context, _, _ io.Writer, args ...string) error {
		calls = append(calls, args)
		if fail != "" && args[0] == fail {
			return errors.New("exit status 1")
		}
		return nil
	}
	t.Cleanup(func() { run = old })
	return &calls
}

func TestCheckDir(t *testing.T) {
	for _, p := range []string{"", "rel/dir", "/", "/a/../b", "/a/", "/a//b", "/a\nb", "C:\\x", "/" + strings.Repeat("a", MaxPath)} {
		if _, err := checkDir(p); err == nil {
			t.Errorf("checkDir(%q) accepted", p)
		}
	}
	if _, err := checkDir("/var/backups/cd/2026-11-14"); err != nil {
		t.Error(err)
	}
}

func TestMainUsage(t *testing.T) {
	var out, errb bytes.Buffer
	for _, a := range [][]string{nil, {"backup"}, {"dump", "/x"}, {"backup", "/x", "y"}, {"backup", "x"}} {
		if c := Main("1", a, &out, &errb); c != 2 {
			t.Errorf("%q: code %d", a, c)
		}
	}
	if c := Main("1.2", []string{"version"}, &out, &errb); c != 0 || !strings.Contains(out.String(), "cd-db 1.2") {
		t.Errorf("version: %d %q", c, out.String())
	}
}

func TestBackupArgs(t *testing.T) {
	calls := fake(t, "")
	dir := filepath.ToSlash(filepath.Join(t.TempDir(), "neu"))
	if err := Backup(context.Background(), dir, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"--backup", "--galera-info", "--target-dir=" + dir}, {"--prepare", "--target-dir=" + dir}}
	if !reflect.DeepEqual(*calls, want) {
		t.Errorf("calls %q", *calls)
	}
}

func TestBackupRefusesExisting(t *testing.T) {
	calls := fake(t, "")
	if err := Backup(context.Background(), t.TempDir(), io.Discard, io.Discard); err == nil || len(*calls) != 0 {
		t.Errorf("err %v, calls %q", err, *calls)
	}
}

func TestBackupFailureStopsBeforePrepare(t *testing.T) {
	calls := fake(t, "--backup")
	dir := filepath.Join(t.TempDir(), "neu")
	if err := Backup(context.Background(), dir, io.Discard, io.Discard); err == nil || len(*calls) != 1 {
		t.Errorf("err %v, calls %q", err, *calls)
	}
}

func TestRestore(t *testing.T) {
	calls := fake(t, "")
	dir := t.TempDir()
	if err := Restore(context.Background(), dir, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	if want := [][]string{{"--copy-back", "--target-dir=" + dir}}; !reflect.DeepEqual(*calls, want) {
		t.Errorf("calls %q", *calls)
	}
	f := filepath.Join(dir, "f")
	os.WriteFile(f, nil, 0o600)
	for _, p := range []string{filepath.Join(dir, "fehlt"), f} {
		if err := Restore(context.Background(), p, io.Discard, io.Discard); err == nil {
			t.Errorf("%s accepted", p)
		}
	}
	fake(t, "--copy-back")
	if err := Restore(context.Background(), dir, io.Discard, io.Discard); err == nil {
		t.Error("copy-back failure swallowed")
	}
}
