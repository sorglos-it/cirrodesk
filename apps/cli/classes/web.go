package classes

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/sorglos-it/nextgencloud/shared/store"
)

// webSkip are the paths of apps/web that never go into a release: the tests
// and the test link, the same two tools/pack-release.sh excludes, and
// package.json, which only makes Node read the test mock as ESM; Apache and the
// shell never read it.
var webSkip = []string{"tests", "public/assets/js/net/testlink.js", "package.json"}

// PackWeb writes the tree of apps/web in dir as a release file of ngc-web: a
// UI tar (store.UIWriter), every file under ui/, in sorted order, no
// directory entries. That is the form the controller unpacks and
// store.WalkUI reads (S7, decided 15.11.2026); the old form with the names
// in the root was refused by both. A link, a device or any other non-regular
// file stops the run: nothing an installation would unpack differently from
// what was packed.
//
// shortcut: it packs the working tree, so a Windows checkout with
// core.autocrlf rewrites line ends; build releases on Linux or from a clean
// checkout until the release packer takes its files from git.
func PackWeb(dir, out string) error {
	var files []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if slices.Contains(webSkip, rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("%w: %s is no regular file", ErrPack, rel)
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return err
	}
	slices.Sort(files)
	tmp := out + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	w := store.NewUIWriter(f)
	err = func() error {
		for _, rel := range files {
			if err := addWebFile(w, dir, rel); err != nil {
				return err
			}
		}
		return w.Close()
	}()
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(tmp)
		if errors.Is(err, store.ErrPackage) {
			return fmt.Errorf("%w: the contract refuses a name or a size in %s", ErrPack, dir)
		}
		return err
	}
	return os.Rename(tmp, out)
}

func addWebFile(w *store.UIWriter, dir, rel string) error {
	src, err := os.Open(filepath.Join(dir, filepath.FromSlash(rel)))
	if err != nil {
		return err
	}
	defer src.Close()
	st, err := src.Stat()
	if err != nil {
		return err
	}
	if err := w.File("ui/"+rel, st.Size(), src); err != nil {
		return fmt.Errorf("%w: %s", err, rel)
	}
	return nil
}
