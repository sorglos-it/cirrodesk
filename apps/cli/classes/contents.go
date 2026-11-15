package classes

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sorglos-it/nextgencloud/shared/store"
)

// PartNames are the parts a release knows, in the order releaseParts lists
// them, each with the architectures store.Release.Validate allows for it.
// It exists so that a message about a wrong part can name the right ones:
// the tool knows the closed set, and a tool that knows the valid values
// says them.
//
// It is derived from releaseParts and not written out a second time, so a
// part that is added there cannot be missing from a message.
func PartNames() []string {
	var out []string
	for _, p := range releaseParts {
		if i := slices.Index(out, p.part); i < 0 {
			out = append(out, p.part)
		}
	}
	return out
}

// PartList is PartNames with the architectures behind each part, as one
// line for a message: "ngc-web (all), ngc-api (amd64, arm64), …".
func PartList() string {
	var parts []string
	for _, name := range PartNames() {
		var arch []string
		for _, p := range releaseParts {
			if p.part == name {
				arch = append(arch, p.arch)
			}
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", name, strings.Join(arch, ", ")))
	}
	return strings.Join(parts, ", ")
}

// ReleaseFileNames are the file names Release looks for in a directory, in
// its own order. Whoever builds a release had to guess these names until
// now — they lived in releaseParts and nowhere a reader could see them.
func ReleaseFileNames() []string {
	out := make([]string, 0, len(releaseParts))
	for _, p := range releaseParts {
		out = append(out, p.part+"-"+p.arch+releaseExt)
	}
	return out
}

// ReleaseFileName is the file name of one part and architecture, and
// whether the two name a file of a release at all.
func ReleaseFileName(part, arch string) (string, bool) {
	for _, p := range releaseParts {
		if p.part == part && p.arch == arch {
			return p.part + "-" + p.arch + releaseExt, true
		}
	}
	return "", false
}

// KnownPart reports whether part is a part of a release, whatever
// architecture.
func KnownPart(part string) bool { return slices.Contains(PartNames(), part) }

// ErrContents reports a release file whose tar does not hold what the
// installation scripts take out of it.
var ErrContents = errors.New("ngc-store: the release file holds the wrong content")

// partContents is what a release tar of one part must hold. The form is
// the one the ulh scripts rely on; it is written down here because until
// now it existed only in those scripts, and a directory packed against
// another idea of it was only found on the server it was installed on.
//
//   - need: names that must be in the tar.
//   - closed: the tar may hold nothing besides need (no further files and
//     no directories). An archive of plain programs is a closed set; the
//     web archive is not, it carries a whole interface.
//
// Sources, all in this repo: shared/ulh.md names what each script takes out
// of a release (ngc-api carries "ngc-api und ngc-store", ngc-controller
// carries "ngc-controller", ngc-web puts the interface under /var/www/ngc
// with public/ as the web root), and docs/STAND.md, "Auslieferungskette",
// names templates/apache-vhost.conf as the file ngc-web.sh reads out of the
// release. The scripts themselves live in the ulh repository and are not
// readable from here.
var partContents = map[string]struct {
	need   []string
	closed bool
	ui     bool // the tar must pass store.WalkUI as a whole
}{
	// S7 (decided 15.11.2026): the web archive is a UI tar, everything under
	// ui/, the form the controller unpacks and store.WalkUI reads.
	store.ReleaseWeb:        {need: []string{"ui/VERSION", "ui/public/index.php", "ui/templates/apache-vhost.conf"}, ui: true},
	store.ReleaseAPI:        {need: []string{"ngc-api", "ngc-store"}, closed: true},
	store.ReleaseController: {need: []string{"ngc-controller"}, closed: true},
}

// CheckContents judges the tar at path against the form the part's
// installation script expects, and returns nil or an error whose text names
// what is missing or what is there too much. It reads the file and writes
// nothing.
//
// Everything a message names comes out of the archive that was just read,
// so the answer is about this file and not about a rule somebody once
// remembered. A name that is no name of the contract (store.ValidPath, the
// rule every reader of a package applies) is refused before it is compared:
// a release file is unpacked on a server, and an entry that walks out of
// its directory is the one defect that must not be reported as a missing
// file.
func CheckContents(path, part string) error {
	want, ok := partContents[part]
	if !ok {
		return fmt.Errorf("%w: unknown part %q — known parts: %s", ErrContents, part, PartList())
	}
	names, err := tarNames(path)
	if err != nil {
		return err
	}
	var missing []string
	for _, n := range want.need {
		if !slices.Contains(names, n) {
			missing = append(missing, n)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s misses %s (needed: %s)", ErrContents,
			filepathBase(path), strings.Join(missing, ", "), strings.Join(want.need, ", "))
	}
	if want.closed {
		var extra []string
		for _, n := range names {
			if !slices.Contains(want.need, n) {
				extra = append(extra, n)
			}
		}
		if len(extra) > 0 {
			return fmt.Errorf("%w: %s holds %s besides %s — this archive holds those files and nothing else, no directories",
				ErrContents, filepathBase(path), strings.Join(extra, ", "), strings.Join(want.need, ", "))
		}
	}
	if want.ui {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("%w: %s is unreadable", ErrContents, filepathBase(path))
		}
		defer f.Close()
		if err := store.WalkUI(f, func(store.PackFile, io.Reader) error { return nil }); err != nil {
			return fmt.Errorf("%w: %s is no UI tar (store.WalkUI refuses it; pack it with ngc-store web)",
				ErrContents, filepathBase(path))
		}
		v, err := WebVersion(path)
		if err != nil {
			return err
		}
		if !store.ValidVersion(v) {
			return fmt.Errorf("%w: %s: ui/VERSION is no version line (MAJOR.MINOR.PATCH[.BUILD] and a newline)",
				ErrContents, filepathBase(path))
		}
	}
	return nil
}

// maxVersionFile bounds ui/VERSION: a version of four numbers and a newline.
const maxVersionFile = 64

// WebVersion returns the version a ngc-web release file names in ui/VERSION,
// without its newline. The controller reads the same file to know which
// version is installed, so the form is one line, store.ValidVersion and
// exactly one newline, as apps/web/VERSION. A file of another form returns "",
// which no version check passes.
func WebVersion(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s is unreadable", ErrContents, filepathBase(path))
	}
	defer f.Close()
	var v []byte
	found := false
	err = store.WalkUI(f, func(pf store.PackFile, body io.Reader) error {
		if pf.Path != "ui/VERSION" {
			return nil
		}
		found = true
		b, err := io.ReadAll(io.LimitReader(body, maxVersionFile+1))
		v = b
		return err
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s is no UI tar", ErrContents, filepathBase(path))
	}
	if !found {
		return "", fmt.Errorf("%w: %s misses ui/VERSION", ErrContents, filepathBase(path))
	}
	s := string(v)
	if len(v) > maxVersionFile || !strings.HasSuffix(s, "\n") || strings.Count(s, "\n") != 1 {
		return "", nil
	}
	return strings.TrimSuffix(s, "\n"), nil
}

// filepathBase is filepath.Base under its own name, so a message always
// names the file and never the path of the machine it was checked on.
func filepathBase(p string) string { return filepath.Base(p) }

// tarNames are the names in the tar at p, directories as they stand and
// without a trailing slash, sorted. A file that is no tar, a name the
// contract refuses and an entry that is neither a file nor a directory are
// each their own reason: they are three different mistakes while packing.
func tarNames(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %s is unreadable", ErrContents, filepathBase(p))
	}
	defer f.Close()
	r := tar.NewReader(f)
	var names []string
	for {
		h, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %s is no tar", ErrContents, filepathBase(p))
		}
		name := strings.TrimSuffix(path.Clean(h.Name), "/")
		if !store.ValidPath(name) {
			return nil, fmt.Errorf("%w: %s holds the name %q, which no installation unpacks",
				ErrContents, filepathBase(p), h.Name)
		}
		switch h.Typeflag {
		case tar.TypeReg, tar.TypeDir:
		default:
			return nil, fmt.Errorf("%w: %s holds %q, which is neither a file nor a directory",
				ErrContents, filepathBase(p), name)
		}
		names = append(names, name)
	}
	slices.Sort(names)
	return names, nil
}

// CheckDirContents judges every release file that lies in dir. It returns
// the names it looked at and the first reason it found, so a caller can say
// how much was judged: a directory where nothing was looked at is the one
// answer that must not read like "ok" (docs: messen.md, "Ein Lauf, der
// keinen Fall trifft, meldet trotzdem alles in Ordnung").
func CheckDirContents(dir string) ([]string, error) {
	var seen []string
	for _, p := range releaseParts {
		name := p.part + "-" + p.arch + releaseExt
		full := filepath.Join(dir, name)
		if _, err := os.Stat(full); errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return seen, err
		}
		seen = append(seen, name)
		if err := CheckContents(full, p.part); err != nil {
			return seen, err
		}
	}
	return seen, nil
}

// listedParts names the parts a release carries, for the message of a
// check that asked for one it does not. It reads the release that was just
// verified, so it says what is there and not what should be.
func listedParts(r *store.Release) string {
	if len(r.Files) == 0 {
		return "no file at all"
	}
	out := make([]string, 0, len(r.Files))
	for i := range r.Files {
		out = append(out, r.Files[i].Part+" "+r.Files[i].Arch)
	}
	slices.Sort(out)
	return strings.Join(out, ", ")
}
