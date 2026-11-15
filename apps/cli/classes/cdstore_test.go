package classes

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/cirrodesk/shared/sign"
	"github.com/sorglos-it/cirrodesk/shared/store"
	"github.com/sorglos-it/cirrodesk/shared/trust"
)

// placeholderStoreURL stands for the address of CirroDesk's own store until
// Thomas names it (S1, decided A on 15.11.2026). docs/cd-store-betreiben.md
// uses the same placeholder.
const placeholderStoreURL = "https://store.example.invalid/cd"

// TestFirstCDStore runs the steps of docs/cd-store-betreiben.md, call for
// call, with a test key: pack the calculator, sign the first catalog of the
// global store, name its key, check the store. If this goes red, the guide
// is wrong.
func TestFirstCDStore(t *testing.T) {
	if !trust.ValidStoreURL(placeholderStoreURL) {
		t.Fatalf("the placeholder %q is no store address", placeholderStoreURL)
	}
	site := t.TempDir()
	keyFile := filepath.Join(t.TempDir(), "cd-store.key")
	call := func(args ...string) string {
		t.Helper()
		var out, errb bytes.Buffer
		if rc := Run(args, &out, &errb); rc != 0 {
			t.Fatalf("cd-store %s: exit %d: %s", strings.Join(args, " "), rc, errb.String())
		}
		return out.String()
	}

	call("pack", filepath.Join("..", "..", "store-calculator"), filepath.Join(site, "org.sorglos.calc-1.0.0.cd"))
	created := call("catalog", site, keyFile, trust.StoreGlobal, "1", placeholderStoreURL)
	named := call("key", keyFile)
	if !strings.Contains(created, named) {
		t.Errorf("catalog said %q, key says %q: they must name the same key", created, named)
	}
	if out := call("check", "-store", trust.StoreGlobal, site); out != "store ok\n" {
		t.Fatalf("check said %q", out)
	}

	// What goes up to the web space, and nothing else: the key stays home.
	entries, err := os.ReadDir(site)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if want := "catalog.json org.sorglos.calc-1.0.0.cd " + StoreKeyFile; strings.Join(names, " ") != want {
		t.Errorf("the site holds %q, want %q", strings.Join(names, " "), want)
	}

	b, err := os.ReadFile(filepath.Join(site, CatalogFile))
	if err != nil {
		t.Fatal(err)
	}
	env, err := sign.ParseEnvelope(b, sign.MaxPayload)
	if err != nil {
		t.Fatal(err)
	}
	var c store.Catalog
	if err := json.Unmarshal(env.Payload, &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Apps) != 1 || c.Apps[0].ID != "org.sorglos.calc" {
		t.Fatalf("catalog apps %+v", c.Apps)
	}
	if u := c.Apps[0].Versions[0].URL; u != placeholderStoreURL+"/org.sorglos.calc-1.0.0.cd" {
		t.Errorf("download link %q", u)
	}
}
