package classes

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/shared/store"
)

// TestPackWebPassesWalkUI packs the real apps/web the way a release is
// packed and reads it back the way the controller does: store.WalkUI, and
// the release check of ngc-store. The tests of apps/web stay out.
func TestPackWebPassesWalkUI(t *testing.T) {
	out := filepath.Join(t.TempDir(), "ngc-web-all.tar")
	if err := PackWeb(filepath.Join("..", "..", "web"), out); err != nil {
		t.Fatalf("PackWeb: %v", err)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var names []string
	if err := store.WalkUI(f, func(pf store.PackFile, _ io.Reader) error {
		names = append(names, pf.Path)
		return nil
	}); err != nil {
		t.Fatalf("store.WalkUI refuses what ngc-store web wrote: %v", err)
	}
	for _, n := range names {
		if strings.Contains(n, "tests/") || strings.HasSuffix(n, "testlink.js") {
			t.Errorf("%s is in the release", n)
		}
	}
	if len(names) == 0 {
		t.Fatal("no file was packed")
	}
	if err := CheckContents(out, store.ReleaseWeb); err != nil {
		t.Errorf("CheckContents: %v", err)
	}
}

func TestPackWebRefusesALink(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "public", "index.php"), []byte("<?php"))
	if err := os.Symlink("/etc/passwd", filepath.Join(dir, "public", "x")); err != nil {
		t.Skipf("no symlinks here: %v", err)
	}
	out := filepath.Join(t.TempDir(), "w.tar")
	if err := PackWeb(dir, out); err == nil {
		t.Fatal("a symlink was packed")
	}
	if _, err := os.Stat(out); err == nil {
		t.Error("a refused run left a release file")
	}
}
