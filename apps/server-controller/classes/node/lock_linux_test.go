//go:build linux

package node

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/apps/server-controller/classes/state"
	"github.com/sorglos-it/nextgencloud/shared/config"
)

// E32: a second ngc-controller on the same database does not start, and it
// says so before it opens anything. Once the first lets go, a start goes on.
func TestASecondControllerOnOneDatabaseDoesNotStart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "controller.json")
	cfg := writeController(t, path)
	cfg.DB = filepath.Join(t.TempDir(), "controller.db")
	// pair left the store there (A2); what is in it this case never reads.
	if err := os.WriteFile(cfg.DB, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := config.SaveController(path, cfg); err != nil {
		t.Fatal(err)
	}
	unlock, err := lockRun(cfg.DB)
	if err != nil {
		t.Fatalf("the first lock: %v", err)
	}
	if _, err := lockRun(cfg.DB); !errors.Is(err, errRunning) {
		t.Fatalf("a second lock while the first is held: %v, want errRunning", err)
	}
	opened := false
	open := func(*config.Controller) (*state.Store, error) {
		opened = true
		return nil, errors.New("the case stops here")
	}
	var out, errOut strings.Builder
	if code := runWith("1.2.3", path, open, &out, &errOut); code != 1 {
		t.Fatalf("runWith = %d, want 1", code)
	}
	if !strings.Contains(errOut.String(), msgRunning) {
		t.Errorf("stderr = %q, want %q", errOut.String(), msgRunning)
	}
	if opened {
		t.Error("the second controller opened the database")
	}

	// Let go, and a start reaches the store again.
	unlock()
	errOut.Reset()
	if code := runWith("1.2.3", path, open, &out, &errOut); code != 1 || !opened {
		t.Fatalf("after the first let go: code %d, opened %v, want the store reached", code, opened)
	}
	if strings.Contains(errOut.String(), msgRunning) {
		t.Errorf("stderr = %q names a running controller after it let go", errOut.String())
	}
}
