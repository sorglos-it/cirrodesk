package release

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/shared/store"
)

// WHAT THIS FILE MEASURES, and why it still stands after S7.
//
// Until 15.11.2026 two halves of the release path could not meet:
// tools/pack-release.sh wrote the ngc-web file with "git archive" over apps/web,
// names in the root, and the controller (update/install.go, unpack) reads an
// ngc-web release file with store.WalkUI, a UI tar rooted at "ui/".
//
// store.WalkUI knows exactly ONE tar dialect (the rules at the top of
// shared/store/ui.go), and a real git archive breaks it in three independent
// places. All three are MEASURED, on 12.10.2026, by reading the header block of
// the archive git actually wrote
// ("git archive --format=tar --mtime=@1760000000 HEAD:docs", bytes 100-147):
//
//	mode   git writes "0000664", WalkUI requires "0000644" (a directory
//	       "0000755")
//	mtime  git writes the pinned date, WalkUI requires "00000000000".
//	       "--mtime=@0" does NOT help: the field then read 015263246526,
//	       also measured.
//	name   no "ui/" prefix at all
//
// S7 (Thomas, 15.11.2026) answered it the second of the two ways this file
// named: shared/ stays as it is, and the packer stops handing out git archive.
// pack-release.sh still lets git archive choose the bytes of the commit, unpacks
// them, and writes the release file with "ngc-store web" through
// store.UIWriter - fixed headers, sorted names, so reproducible without a pinned
// date (TestPackIsReproducible). TestPackedArchivePassesWalkUI holds that the
// result is what the controller reads.
//
// The cases below stay, because they say WHY git archive alone can never be
// the release file again. They take each cause ONE AT A TIME, starting from
// bytes store.WalkUI accepts. Without the base case the others would all be red
// for whatever reason happened to be first.

// uiBlocks builds a one-file UI tar block by block, so that each field can
// be set on its own. mode and mtime are the fields git archive spells
// differently; name is the whole name including any prefix.
//
// The checksum is computed the way the contract spells it (shared/store/ui.go:
// six octal digits, NUL and a space, over the 512 bytes with these eight
// counted as spaces), so a patched header is a VALID header of whatever it
// now says - the refusals below are about the field, never about a checksum
// gone stale.
func uiBlocks(name, mode, mtime string, body []byte) []byte {
	var out []byte
	hdr := func(flag byte, nm, md, mt string, size int64) []byte {
		b := make([]byte, 512)
		copy(b[0:100], nm)
		copy(b[100:108], md)
		copy(b[108:124], "0000000\x000000000\x00")
		copy(b[124:136], octal(size, 12))
		copy(b[136:148], mt)
		b[156] = flag
		copy(b[257:265], "ustar\x0000")
		// The checksum last, over the block with its own field as spaces.
		for i := 148; i < 156; i++ {
			b[i] = ' '
		}
		sum := 0
		for _, c := range b {
			sum += int(c)
		}
		copy(b[148:156], octal(int64(sum), 7)+" ")
		return b
	}
	// A directory entry, then the file. The entry is NOT required - measured
	// by switching this block off, and the base case stayed green, because
	// WalkUI implies a parent from the path ("ui/a/b" implies "ui" and
	// "ui/a"). It is written because a real archive carries one, so the bytes
	// below are the shape of a real archive and not the thinnest one that
	// passes.
	if i := strings.LastIndexByte(name, '/'); i > 0 {
		out = append(out, hdr('5', name[:i+1], "0000755\x00", "00000000000\x00", 0)...)
	}
	out = append(out, hdr('0', name, mode, mtime, int64(len(body)))...)
	out = append(out, body...)
	if pad := (512 - len(body)%512) % 512; pad > 0 {
		out = append(out, make([]byte, pad)...)
	}
	return append(out, make([]byte, 1024)...) // the end marker
}

// octal formats v as octal digits and a NUL, width bytes in all.
func octal(v int64, width int) string {
	s := strconv.FormatInt(v, 8)
	return strings.Repeat("0", width-1-len(s)) + s + "\x00"
}

const (
	uiMode    = "0000644\x00"
	uiTime    = "00000000000\x00"
	gitMode   = "0000664\x00"     // measured: what git archive writes
	gitTime   = "15071674000\x00" // MEASURED: the field git archive wrote
	uiName    = "ui/public/index.php"
	plainName = "public/index.php" // measured: what pack-release.sh wrote before S7
)

// walk hands ar to store.WalkUI and says how many files came through.
func walk(ar []byte) (int, error) {
	n := 0
	err := store.WalkUI(bytes.NewReader(ar), func(f store.PackFile, body io.Reader) error {
		n++
		_, err := io.Copy(io.Discard, body)
		return err
	})
	return n, err
}

// 1. THE BASE CASE, and nothing below means anything without it: bytes this
// file builds, with every field as the contract spells it, go through
// store.WalkUI. If this is red, uiBlocks is broken and the four refusals
// below measure that instead of the thing they name.
func TestTheBuiltArchiveIsOneWalkUIAccepts(t *testing.T) {
	n, err := walk(uiBlocks(uiName, uiMode, uiTime, []byte("<?php\n")))
	if err != nil {
		t.Fatalf("uiBlocks is broken, not the contract: a correct UI tar must pass, got %v", err)
	}
	if n != 1 {
		t.Fatalf("the archive delivered %d files, expected 1", n)
	}
}

// 2. The mode git archive writes. Measured 0000664; the contract wants
// 0000644.
func TestTheModeGitArchiveWritesIsRefused(t *testing.T) {
	refused(t, uiBlocks(uiName, gitMode, uiTime, []byte("<?php\n")),
		"the mode of a git archive entry (0000664)")
}

// 3. The date the packer pins. The contract wants the field all zeros, and
// the packer pins it to the commit date ON PURPOSE - that is how it makes
// two builds of one commit give one checksum. The two requirements are
// opposed, and that is the heart of the decision.
func TestThePinnedDateIsRefused(t *testing.T) {
	refused(t, uiBlocks(uiName, uiMode, gitTime, []byte("<?php\n")),
		"a pinned modification time")
}

// 4. The missing prefix, which is the one cause the BLOCKED paragraph in the
// controller already names.
func TestTheMissingPrefixIsRefused(t *testing.T) {
	refused(t, uiBlocks(plainName, uiMode, uiTime, []byte("<?php\n")),
		"an entry name without the \"ui/\" prefix")
}

// 5. All three together, which is what a release built with git archive alone
// was until S7. It is not a fourth cause; it is the shape that would arrive on
// a server if a packer went back to it.
func TestTheArchiveOfGitArchiveAloneIsRefused(t *testing.T) {
	refused(t, uiBlocks(plainName, gitMode, gitTime, []byte("<?php\n")),
		"the shape tools/pack-release.sh wrote before S7")
}

// 6. And the archive pack-release.sh writes since S7 is one store.WalkUI
// takes, whole, with the entry point ngc-web.sh looks for and ui/VERSION.
func TestPackedArchivePassesWalkUI(t *testing.T) {
	t.Parallel()
	needTools(t)
	repo := webFixture(t, nil)
	out, log, err := pack(t, repo, "main", "3")
	if err != nil {
		t.Fatalf("pack-release.sh: %v\n%s", err, log)
	}
	ar, err := os.ReadFile(filepath.Join(out, "ngc-web-all.tar"))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	if err := store.WalkUI(bytes.NewReader(ar), func(f store.PackFile, body io.Reader) error {
		seen[f.Path] = true
		_, err := io.Copy(io.Discard, body)
		return err
	}); err != nil {
		t.Fatalf("store.WalkUI refuses the archive pack-release.sh wrote: %v", err)
	}
	for _, want := range []string{uiName, "ui/VERSION", "ui/templates/apache-vhost.conf"} {
		if !seen[want] {
			t.Errorf("the archive lacks %s", want)
		}
	}
}

// refused checks that store.WalkUI says no, and says no as a package error.
func refused(t *testing.T, ar []byte, what string) {
	t.Helper()
	n, err := walk(ar)
	if err == nil {
		t.Fatalf("store.WalkUI accepted %s (%d files). If that is now the agreed form, "+
			"docs/release-weg.md has to say so and this case has to be rewritten", what, n)
	}
	if !errors.Is(err, store.ErrPackage) {
		t.Fatalf("%s was refused, but not as a package error: %v", what, err)
	}
	t.Logf("refused as expected: %s", what)
}
