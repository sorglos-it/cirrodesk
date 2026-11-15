package main

// The store side of the app: ngc-store packs the directory, a catalog signed
// with a test key passes check -store, and the package passes store.Walk,
// which the controller accepts packages with (install.Accept). The real
// image is built on a build machine with podman (CLAUDE.md); store.Walk
// only needs one to be there, so the test packs a stand-in of a few bytes.

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

func copyTo(t *testing.T, path string, b []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPackage(t *testing.T) {
	dir := t.TempDir()
	src := t.TempDir()
	for _, name := range []string{"manifest.json", "ui/index.html", "ui/app.js", "ui/app.css"} {
		b, err := os.ReadFile(filepath.FromSlash(name))
		if err != nil {
			t.Fatal(err)
		}
		copyTo(t, filepath.Join(src, filepath.FromSlash(name)), b)
	}
	copyTo(t, filepath.Join(src, "services", "amd64.tar"), []byte("stand-in image"))
	pkg := filepath.Join(dir, "org.sorglos.files-1.0.0.ngc")
	if out := run(t, "pack", src, pkg); !strings.HasPrefix(out, "org.sorglos.files 1.0.0\n") {
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
	if want := "ui/app.css ui/app.js ui/index.html services/amd64.tar"; strings.Join(names, " ") != want {
		t.Errorf("files %q, want %q", strings.Join(names, " "), want)
	}
	// Every action of the UI is in the manifest, served by the one service.
	for _, a := range []string{"list", "get", "put", "mkdir", "rename", "remove"} {
		if act, ok := m.Action(a); !ok || act.Service != "files" || act.Stream != (a == "get" || a == "put") {
			t.Errorf("action %s: %+v %v", a, act, ok)
		}
	}
	if len(m.Services) != 1 || m.Services[0].Command[0] != "/usr/lib/ngc/ngc-files" || len(m.Needs) != 0 {
		t.Errorf("service: %+v needs %v", m.Services, m.Needs)
	}
	run(t, "catalog", dir, filepath.Join(t.TempDir(), "store.key"), "test", "1")
	if out := run(t, "check", "-store", "test", dir); out != "store ok\n" {
		t.Fatalf("check said %q", out)
	}
}

// TestUI runs the node tests of ui/app.js; skipped where node is missing.
func TestUI(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not found")
	}
	out, err := exec.Command(node, "--test", filepath.Join("tests", "files.test.js")).CombinedOutput()
	if err != nil {
		t.Fatalf("node --test: %v\n%s", err, out)
	}
}
