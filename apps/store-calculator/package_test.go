// Package storecalculator holds no Go code: the app is ui/ and manifest.json.
// This test proves the directory is a store app - ngc-store packs it, signs a
// catalog with a test key, checks the store, and the package passes
// store.Walk, which the controller accepts packages with (install.Accept).
package storecalculator

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/apps/cli/classes"
	"github.com/sorglos-it/nextgencloud/shared/store"
)

func run(t *testing.T, args ...string) string {
	t.Helper()
	var out, errb bytes.Buffer
	if rc := classes.Run(args, &out, &errb); rc != 0 {
		t.Fatalf("ngc-store %s: exit %d: %s", strings.Join(args, " "), rc, errb.String())
	}
	return out.String()
}

func TestPackage(t *testing.T) {
	dir := t.TempDir()
	pkg := filepath.Join(dir, "org.sorglos.calc-1.0.0.ngc")
	if out := run(t, "pack", ".", pkg); !strings.HasPrefix(out, "org.sorglos.calc 1.0.0\n") {
		t.Fatalf("pack said %q", out)
	}

	f, err := os.Open(pkg)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var names []string
	m, err := store.Walk(f, func(_ *store.Manifest, p store.PackFile, _ io.Reader) error {
		names = append(names, p.Path)
		return nil
	})
	if err != nil {
		t.Fatalf("store.Walk refuses the package: %v", err)
	}
	if m.UI == nil || m.UI.Entry != "index.html" || len(m.Services) != 0 {
		t.Errorf("manifest: %+v", m)
	}
	if want := "ui/app.css ui/app.js ui/index.html"; strings.Join(names, " ") != want {
		t.Errorf("files %q, want %q", strings.Join(names, " "), want)
	}

	run(t, "catalog", dir, filepath.Join(t.TempDir(), "store.key"), "test", "1")
	if out := run(t, "check", "-store", "test", dir); out != "store ok\n" {
		t.Fatalf("check said %q", out)
	}
}

// TestParser runs the node tests of ui/app.js, fuzz included. Like the other
// tests that need a tool, it skips where node is missing.
func TestParser(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not found")
	}
	out, err := exec.Command(node, "--test", filepath.Join("tests", "calc.test.js")).CombinedOutput()
	if err != nil {
		t.Fatalf("node --test: %v\n%s", err, out)
	}
}
