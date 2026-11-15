package classes

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/cirrodesk/shared/service"
)

const user = "AAAAAAAAAAAAAAAAAAAAAA" // 16 zero bytes, base64url: a trust.ValidID

type result struct {
	code   int
	out    []byte
	answer Answer
}

func call(t *testing.T, data, usr, action, params string, stdin []byte) result {
	t.Helper()
	env := map[string]string{service.EnvUser: usr, service.EnvAction: action, service.EnvParams: params}
	var out, errb bytes.Buffer
	code := Run(func(k string) string { return env[k] }, data, bytes.NewReader(stdin), &out, &errb)
	r := result{code: code, out: out.Bytes()}
	if action != "get" && action != "put" && code == 0 {
		if err := json.Unmarshal(out.Bytes(), &r.answer); err != nil {
			t.Fatalf("%s %s: answer %q: %v", action, params, out.String(), err)
		}
	}
	return r
}

func p(path string) string {
	b, _ := json.Marshal(Params{Path: path})
	return string(b)
}

func TestFiles(t *testing.T) {
	u := user
	data := t.TempDir()

	if r := call(t, data, u, "list", "", nil); r.answer.Error != "" || len(r.answer.Entries) != 0 {
		t.Fatalf("empty list: %+v", r)
	}
	if r := call(t, data, u, "mkdir", p("Fotos 2026"), nil); r.answer.Error != "" {
		t.Fatalf("mkdir: %+v", r.answer)
	}
	if r := call(t, data, u, "mkdir", p("Fotos 2026"), nil); r.answer.Error != ErrExists {
		t.Fatalf("mkdir twice: %+v", r.answer)
	}
	body := []byte("Hallo, Welt\n")
	if r := call(t, data, u, "put", p("Fotos 2026/gruß.txt"), body); r.code != 0 {
		t.Fatalf("put: %d", r.code)
	}
	if r := call(t, data, u, "get", p("Fotos 2026/gruß.txt"), nil); r.code != 0 || !bytes.Equal(r.out, body) {
		t.Fatalf("get: %d %q", r.code, r.out)
	}
	r := call(t, data, u, "list", p("Fotos 2026"), nil)
	if len(r.answer.Entries) != 1 || r.answer.Entries[0].Name != "gruß.txt" || r.answer.Entries[0].Size != int64(len(body)) {
		t.Fatalf("list: %+v", r.answer)
	}
	if r := call(t, data, u, "remove", p("Fotos 2026"), nil); r.answer.Error != ErrNotEmpty {
		t.Fatalf("remove full folder: %+v", r.answer)
	}
	ren, _ := json.Marshal(Params{Path: "Fotos 2026/gruß.txt", To: "gruß.txt"})
	if r := call(t, data, u, "rename", string(ren), nil); r.answer.Error != "" {
		t.Fatalf("rename: %+v", r.answer)
	}
	if r := call(t, data, u, "remove", p("Fotos 2026"), nil); r.answer.Error != "" {
		t.Fatalf("remove empty folder: %+v", r.answer)
	}
	if r := call(t, data, u, "get", p("fehlt.txt"), nil); r.code != 1 {
		t.Fatalf("get of a missing file: %d", r.code)
	}
	r = call(t, data, u, "list", "", nil)
	if len(r.answer.Entries) != 1 || r.answer.Entries[0].Name != "gruß.txt" {
		t.Fatalf("list after: %+v", r.answer)
	}
	// Folders first, then names.
	call(t, data, u, "mkdir", p("z"), nil)
	if r := call(t, data, u, "list", "", nil); r.answer.Entries[0].Name != "z" {
		t.Fatalf("order: %+v", r.answer)
	}
	// The files lie below users/<user>, nowhere else.
	if _, err := os.Stat(filepath.Join(data, "users", u, "gruß.txt")); err != nil {
		t.Fatalf("not in the user's folder: %v", err)
	}
}

func TestRefused(t *testing.T) {
	u := user
	data := t.TempDir()
	cases := []struct{ user, action, params, want string }{
		{"", "list", "", ErrBadCall},
		{"../x", "list", "", ErrBadCall},
		{u, "format", "", ErrBadCall},
		{u, "list", `{"path":"a","extra":1}`, ErrBadCall},
		{u, "list", `{"Path":"a"}`, ErrBadCall},
		{u, "list", `{"path":"a","path":"b"}`, ErrBadCall},
		{u, "list", `[]`, ErrBadCall},
		{u, "list", p(".."), ErrBadPath},
		{u, "list", p("a/../.."), ErrBadPath},
		{u, "list", p("/etc"), ErrBadPath},
		{u, "list", p("a\\b"), ErrBadPath},
		{u, "list", p("a//b"), ErrBadPath},
		{u, "list", p(" a"), ErrBadPath},
		{u, "list", p("a\nb"), ErrBadPath},
		{u, "list", p(strings.Repeat("a/", 40) + "a"), ErrBadPath},
		{u, "mkdir", p(""), ErrBadPath},
		{u, "remove", p(""), ErrBadPath},
		{u, "list", p("fehlt"), ErrNotFound},
	}
	for _, c := range cases {
		r := call(t, data, c.user, c.action, c.params, nil)
		if r.code != 0 || r.answer.Error != c.want {
			t.Errorf("%q %s %s: %d %+v, want %s", c.user, c.action, c.params, r.code, r.answer, c.want)
		}
	}
}

// A link that points out of the user's folder is neither listed nor followed.
func TestLinkOut(t *testing.T) {
	u := user
	data := t.TempDir()
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "geheim"), []byte("x"), 0o600)
	home := filepath.Join(data, "users", u)
	os.MkdirAll(home, 0o700)
	if err := os.Symlink(outside, filepath.Join(home, "raus")); err != nil {
		t.Skipf("no symlinks here: %v", err)
	}
	if r := call(t, data, u, "list", "", nil); len(r.answer.Entries) != 0 {
		t.Errorf("the link is listed: %+v", r.answer)
	}
	if r := call(t, data, u, "get", p("raus/geheim"), nil); r.code == 0 {
		t.Errorf("read through the link: %q", r.out)
	}
	if r := call(t, data, u, "put", p("raus/neu"), []byte("x")); r.code == 0 {
		t.Errorf("wrote through the link")
	}
}

// A broken upload leaves no file under the real name and no temporary one.
func TestPutTooLarge(t *testing.T) {
	u := user
	data := t.TempDir()
	call(t, data, u, "put", p("a.txt"), []byte("alt"))
	defer func(n int64) { maxFile = n }(maxFile)
	maxFile = 4
	if r := call(t, data, u, "put", p("a.txt"), []byte("zu lang")); r.code != 1 {
		t.Fatalf("too large: %d", r.code)
	}
	if r := call(t, data, u, "get", p("a.txt"), nil); string(r.out) != "alt" {
		t.Fatalf("the old file changed: %q", r.out)
	}
	ents, _ := os.ReadDir(filepath.Join(data, "users", u))
	if len(ents) != 1 {
		t.Fatalf("left behind: %v", ents)
	}
}

func FuzzValidPath(f *testing.F) {
	for _, s := range []string{"", "a", "a/b", "..", "a/../b", "/a", "a\\b", " a", "ä ö", "a\x00b"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !ValidPath(s) {
			return
		}
		// Whatever passes stays inside: cleaning changes nothing and never
		// climbs or starts at the top.
		if c := path.Clean(s); s != "" && c != s {
			t.Fatalf("%q cleans to %q", s, c)
		}
		if strings.HasPrefix(s, "/") || strings.Contains("/"+s+"/", "/../") || strings.ContainsAny(s, "\\\x00") {
			t.Fatalf("%q passed", s)
		}
	})
}

func FuzzParams(f *testing.F) {
	for _, s := range []string{"", `{"path":"a"}`, `{"path":"a","to":"b"}`, `{"path":".."}`, `{"x":1}`, `null`} {
		f.Add(s)
	}
	u := user
	f.Fuzz(func(t *testing.T, params string) {
		data := t.TempDir()
		env := map[string]string{service.EnvUser: u, service.EnvAction: "list", service.EnvParams: params}
		var out, errb bytes.Buffer
		if code := Run(func(k string) string { return env[k] }, data, nil, &out, &errb); code != 0 {
			t.Fatalf("list exit %d", code)
		}
		var a Answer
		if err := json.Unmarshal(out.Bytes(), &a); err != nil {
			t.Fatalf("answer %q: %v", out.String(), err)
		}
	})
}
