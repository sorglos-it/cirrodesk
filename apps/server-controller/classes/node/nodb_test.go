package node

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sorglos-it/cirrodesk/apps/server-controller/classes/state"
	"github.com/sorglos-it/cirrodesk/shared/config"
)

// A2: on a paired server the store is there, so a file that is gone while its
// directory stands is a lost one. The controller does not start over it, says
// why, and opens nothing - a new, empty store would report no apps at all.
func TestAControllerWithoutItsDatabaseDoesNotStart(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("a controller database cannot be named on windows: config.DB.Validate wants a leading %q - measured on linux instead", "/")
	}
	path := filepath.Join(t.TempDir(), "controller.json")
	cfg := writeController(t, path)
	cfg.DB = filepath.Join(t.TempDir(), "controller.db")
	if err := config.SaveController(path, cfg); err != nil {
		t.Fatal(err)
	}
	opened := false
	open := func(*config.Controller) (*state.Store, error) {
		opened = true
		return nil, os.ErrInvalid
	}
	var out, errOut strings.Builder
	if code := runWith("1.2.3", path, open, &out, &errOut); code != 1 {
		t.Fatalf("runWith = %d, want 1", code)
	}
	if opened {
		t.Fatal("the store was opened, and that would have made a new, empty one")
	}
	if !strings.Contains(errOut.String(), msgNoDB) {
		t.Fatalf("stderr = %q, want %q", errOut.String(), msgNoDB)
	}
	if strings.Contains(errOut.String(), cfg.DB) {
		t.Fatalf("stderr names the database file: %q", errOut.String())
	}

	// INTERMEDIATE: the same server with the file in place reaches the store,
	// so what is measured above is the missing file.
	if err := os.WriteFile(cfg.DB, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	errOut.Reset()
	if code := runWith("1.2.3", path, open, &out, &errOut); code != 1 || !opened {
		t.Fatalf("with the file: code %d, opened %v, want the store reached", code, opened)
	}
	if strings.Contains(errOut.String(), msgNoDB) {
		t.Fatalf("stderr = %q names a missing database although it is there", errOut.String())
	}
}

// pair is where the store is born: afterwards the file is there and opens, and
// a second pairing keeps what is in it.
func TestPairCreatesTheDatabase(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("a controller database cannot be named on windows: config.DB.Validate wants a leading %q - measured on linux instead", "/")
	}
	api := newAPI(t)
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "controller.json")
	cfg := writeController(t, cfgPath)
	createState = func(c *config.Controller) error { // the real one, put back
		st, err := openState(c)
		if err != nil {
			return err
		}
		return st.Close()
	}
	cfg.DB = filepath.Join(dir, "state.db")
	if err := config.SaveController(cfgPath, cfg); err != nil {
		t.Fatal(err)
	}
	var out, errOut strings.Builder
	if code := pair(cfgPath, filepath.Join(dir, "anchor.json"), nil, typed(api.code.String()), &out, &errOut); code != 0 {
		t.Fatalf("pair = %d, stderr = %q", code, errOut.String())
	}
	if lostDB(cfg.DB) {
		t.Fatal("pair left no database")
	}
	st, err := openState(cfg)
	if err != nil {
		t.Fatalf("the store pair made does not open: %v", err)
	}
	st.Close()
}
