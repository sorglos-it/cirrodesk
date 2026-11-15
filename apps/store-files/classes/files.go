// Package classes is the service of org.sorglos.files: it runs once per call
// inside the app's container (contract V3, shared/service) and works on the
// files of the calling user below /data/users/<user>.
package classes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sorglos-it/cirrodesk/shared/enc"
	"github.com/sorglos-it/cirrodesk/shared/service"
	"github.com/sorglos-it/cirrodesk/shared/trust"
	"github.com/sorglos-it/cirrodesk/shared/wire"
)

// maxFile is the largest upload; a var only so the test need not send 1 GiB.
var maxFile int64 = 1 << 30

// Limits. A listing larger than MaxEntries is cut and says so.
const (
	MaxEntries  = 2000
	maxPath     = 1024
	maxSegment  = 255
	maxSegments = 32

	uploadPrefix = ".upload-" // an upload in progress; never listed
)

// Error codes the answer carries; the UI translates them.
const (
	ErrBadCall  = "call"     // unknown action, bad params, no user
	ErrBadPath  = "path"     // a name the service does not accept
	ErrNotFound = "notfound" // nothing at that path
	ErrExists   = "exists"   // the target is already there
	ErrNotEmpty = "notempty" // a folder that still holds something
	ErrTooLarge = "toolarge" // more than maxFile bytes
	ErrFailed   = "failed"   // anything else; the details go to stderr only
)

// Params are the params of every action; which members are needed depends on
// the action.
type Params struct {
	Path string `json:"path"`
	To   string `json:"to,omitempty"`
}

// Entry is one line of a listing.
type Entry struct {
	Name  string `json:"name"`
	Dir   bool   `json:"dir"`
	Size  int64  `json:"size"`
	MTime int64  `json:"mtime"` // Unix seconds
}

// Answer is what list, mkdir, remove and rename write to stdout.
type Answer struct {
	Error   string  `json:"error,omitempty"`
	Entries []Entry `json:"entries,omitempty"`
	Cut     bool    `json:"cut,omitempty"`
}

// failure carries an error code out of an action.
type failure string

func (f failure) Error() string { return string(f) }

// ValidPath reports whether p is a path the service accepts: "" is the
// user's top folder, else up to maxSegments names joined by "/". A name is
// valid UTF-8 without control characters, "/" or "\", not "." or "..", and
// does not start or end with a space. os.Root is the second wall behind this
// one: no path, and no link, leaves the user's folder whatever this allows.
func ValidPath(p string) bool {
	if p == "" {
		return true
	}
	if len(p) > maxPath || !utf8.ValidString(p) {
		return false
	}
	segs := strings.Split(p, "/")
	if len(segs) > maxSegments {
		return false
	}
	for _, s := range segs {
		if s == "" || s == "." || s == ".." || len(s) > maxSegment ||
			strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
			return false
		}
		for _, r := range s {
			if r == '\\' || unicode.IsControl(r) || r == utf8.RuneError {
				return false
			}
		}
	}
	return true
}

// Run is one call. getenv is os.Getenv in the container; data is
// service.DataDir there and a temporary directory in the tests. It returns
// the exit code: 0 also for a refused call of list, mkdir, remove and rename,
// whose answer says why; get and put have no answer to say it in and end
// with 1 instead, so the channel closes with an error.
func Run(getenv func(string) string, data string, stdin io.Reader, stdout, stderr io.Writer) int {
	action := getenv(service.EnvAction)
	stream := action == "get" || action == "put"
	err := run(getenv, data, action, stdin, stdout)
	if err == nil {
		return 0
	}
	code := ErrFailed
	var f failure
	if errors.As(err, &f) {
		code = string(f)
	} else {
		fmt.Fprintf(stderr, "cd-files: %s: %v\n", action, err)
	}
	if stream {
		fmt.Fprintln(stderr, "cd-files:", code)
		return 1
	}
	json.NewEncoder(stdout).Encode(Answer{Error: code})
	return 0
}

func run(getenv func(string) string, data, action string, stdin io.Reader, stdout io.Writer) error {
	user := getenv(service.EnvUser)
	if !trust.ValidID(user) {
		return failure(ErrBadCall)
	}
	var p Params
	if raw := getenv(service.EnvParams); raw != "" {
		if enc.DecodeStrict([]byte(raw), wire.MaxParams, &p) != nil {
			return failure(ErrBadCall)
		}
	}
	if !ValidPath(p.Path) || !ValidPath(p.To) {
		return failure(ErrBadPath)
	}
	home := path.Join("users", user)
	top, err := os.OpenRoot(data)
	if err != nil {
		return err
	}
	defer top.Close()
	if err := top.MkdirAll(home, 0o700); err != nil {
		return err
	}
	root, err := top.OpenRoot(home)
	if err != nil {
		return err
	}
	defer root.Close()
	name := p.Path
	if name == "" {
		name = "."
	}

	switch action {
	case "list":
		return list(root, name, stdout)
	case "mkdir":
		if p.Path == "" {
			return failure(ErrBadPath)
		}
		return answer(stdout, mapErr(root.Mkdir(name, 0o700)))
	case "remove":
		if p.Path == "" {
			return failure(ErrBadPath)
		}
		// Checked here and not read from the error: Linux and Windows name a
		// full folder differently, and the tests run on both.
		if d, err := root.Open(name); err == nil {
			fi, serr := d.Stat()
			var more []fs.FileInfo
			if serr == nil && fi.IsDir() {
				more, _ = d.Readdir(1)
			}
			d.Close()
			if len(more) > 0 {
				return failure(ErrNotEmpty)
			}
		}
		return answer(stdout, mapErr(root.Remove(name)))
	case "rename":
		if p.Path == "" || p.To == "" {
			return failure(ErrBadPath)
		}
		if _, err := root.Lstat(p.To); err == nil {
			return failure(ErrExists)
		}
		return answer(stdout, mapErr(root.Rename(name, p.To)))
	case "get":
		return get(root, name, stdout)
	case "put":
		if p.Path == "" {
			return failure(ErrBadPath)
		}
		return put(root, name, stdin)
	}
	return failure(ErrBadCall)
}

func answer(stdout io.Writer, err error) error {
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(Answer{})
}

// mapErr turns what the file system said into an error code; what has no code
// stays an error and is logged.
func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, fs.ErrNotExist):
		return failure(ErrNotFound)
	case errors.Is(err, fs.ErrExist):
		return failure(ErrExists)
	}
	return err
}

func list(root *os.Root, name string, stdout io.Writer) error {
	d, err := root.Open(name)
	if err != nil {
		return mapErr(err)
	}
	defer d.Close()
	infos, err := d.Readdir(MaxEntries + 1)
	if err != nil && err != io.EOF {
		return mapErr(err)
	}
	a := Answer{Entries: []Entry{}}
	if len(infos) > MaxEntries {
		infos, a.Cut = infos[:MaxEntries], true
	}
	for _, fi := range infos {
		// Links are not shown: os.Root would not follow one out of the folder,
		// and one that stays inside is nothing a user of this app made.
		if fi.Mode()&fs.ModeSymlink != 0 || !ValidPath(fi.Name()) || strings.HasPrefix(fi.Name(), uploadPrefix) {
			continue
		}
		a.Entries = append(a.Entries, Entry{Name: fi.Name(), Dir: fi.IsDir(), Size: fi.Size(), MTime: fi.ModTime().Unix()})
	}
	slices.SortFunc(a.Entries, func(x, y Entry) int {
		if x.Dir != y.Dir {
			if x.Dir {
				return -1
			}
			return 1
		}
		return strings.Compare(x.Name, y.Name)
	})
	return json.NewEncoder(stdout).Encode(a)
}

func get(root *os.Root, name string, stdout io.Writer) error {
	f, err := root.Open(name)
	if err != nil {
		return mapErr(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return failure(ErrNotFound)
	}
	_, err = io.Copy(stdout, f)
	return err
}

// put writes stdin to a temporary name beside the target and renames it at
// the end: a broken upload never leaves half a file under the real name, and
// an old file stays as it was until the new one is complete.
func put(root *os.Root, name string, stdin io.Reader) error {
	dir, base := path.Split(name)
	tmp := dir + uploadPrefix + base
	f, err := root.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return mapErr(err)
	}
	n, err := io.Copy(f, io.LimitReader(stdin, maxFile+1))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil && n > maxFile {
		err = failure(ErrTooLarge)
	}
	if err == nil {
		err = root.Rename(tmp, name)
	}
	if err != nil {
		root.Remove(tmp)
		return mapErr(err)
	}
	return nil
}
