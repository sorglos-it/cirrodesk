package sqlstore_test

// The one assurance the pairing paths of the API rest on, held where the
// property lives.
//
// join.Accept and controllers.Pair are to answer a busy database differently
// from a code that does not verify (busy: 503 with Retry-After, a code that
// does not verify: E_VERIFY). For controllers.Pair that opens no oracle only
// because OnceUse meets the file lock BEFORE it reads a row: if it read first,
// a hash that is a secret of the installation would come back "busy" and one
// that never was would come back "not found", and the two would be tellable
// apart. The API's own cases cannot hold that: the store under them is a stub
// that blocks every call, so it makes the sentence true by construction.

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/sorglos-it/cirrodesk/apps/server-db/classes/sqlite"
	"github.com/sorglos-it/cirrodesk/shared/config"
	"github.com/sorglos-it/cirrodesk/shared/dbapi"
)

// TestOnceUseMeetsTheLockBeforeItReadsARow runs the store the program runs
// (sqlite.Open: the split handles, the write turn, BEGIN IMMEDIATE) against a
// real foreign lock and asks whether a value that EXISTS and one that never
// did get the same answer.
//
// The intermediate check comes first and is the half that makes the rest mean
// anything: without the lock the two answers DIFFER, so the case can tell them
// apart. Only then is the lock taken and the sameness measured. Without that
// step a store that answered "not found" to everything would pass.
func TestOnceUseMeetsTheLockBeforeItReadsARow(t *testing.T) {
	const kind = "pair"
	known := bytes.Repeat([]byte{3}, 32)
	unknown := bytes.Repeat([]byte{4}, 32) // the same shape, nothing behind it
	const now = int64(1790000000)

	path := filepath.Join(t.TempDir(), "state.db")
	s, err := sqlite.Open(onceLockConfig(path), "oncelock", now)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	if err := s.OnceCreate(ctx, kind, known, "a secret of this installation", now+3600); err != nil {
		t.Fatalf("the value under test could not be created: %v", err)
	}

	// Free database: the two have to be tellable apart here, or nothing below
	// is an assurance.
	if _, err := s.OnceUse(ctx, kind, known, now); err != nil {
		t.Fatalf("a value that exists answered %v on a free database", err)
	}
	if _, err := s.OnceUse(ctx, kind, unknown, now); !errors.Is(err, dbapi.ErrNotFound) {
		t.Fatalf("a value that never existed answered %v on a free database, want %v",
			err, dbapi.ErrNotFound)
	}

	// A second value that exists, because the first one is spent now.
	known2 := bytes.Repeat([]byte{5}, 32)
	if err := s.OnceCreate(ctx, kind, known2, "a second secret", now+3600); err != nil {
		t.Fatalf("the second value could not be created: %v", err)
	}

	holdForeignLock(t, path)

	// A deadline of its own, because what is measured is the ANSWER and not
	// the retry chain: without one both calls would take a whole chain. The
	// same figure for both, so the deadline cannot be the difference.
	ask := func(value []byte) error {
		c, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()
		_, err := s.OnceUse(c, kind, value, now)
		return err
	}
	gotKnown, gotUnknown := ask(known2), ask(unknown)

	// Neither may be an answer about the VALUE: under the lock the row was
	// never reached, so saying "not found" or "used" would be saying whether
	// the value is a secret of this installation.
	for name, err := range map[string]error{"known": gotKnown, "unknown": gotUnknown} {
		if err == nil {
			t.Fatalf("the %s value got through although the file is locked; nothing measured", name)
		}
		if errors.Is(err, dbapi.ErrNotFound) || errors.Is(err, dbapi.ErrUsed) {
			t.Errorf("under a foreign lock the %s value is answered with %v -- that is an "+
				"answer about the value, and the pairing paths of the API rest on there "+
				"being none", name, err)
		}
	}
	if gotKnown.Error() != gotUnknown.Error() {
		t.Errorf("under a foreign lock a value that exists is answered %v and one that "+
			"never did %v; the two must not be tellable apart", gotKnown, gotUnknown)
	}
}

// onceLockConfig is the configuration of the case's file, made afresh for
// every opening so that no opening depends on what an earlier one did with
// the key.
func onceLockConfig(path string) config.DB {
	return config.DB{
		Driver: config.DriverSQLite,
		Path:   path,
		Key:    config.NewSecret(bytes.Repeat([]byte{9}, 64)),
	}
}

// holdForeignLock is the foreign writer: a second program that holds the
// write lock of the file until the case is over. BEGIN IMMEDIATE takes the
// lock the moment it begins, and a lock that was never taken is a case that
// measures nothing - so it is checked and not assumed.
func holdForeignLock(t *testing.T, path string) {
	t.Helper()
	db, err := sqlite.OpenDB(onceLockConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.ExecContext(context.Background(), "BEGIN IMMEDIATE"); err != nil {
		t.Fatalf("the foreign lock was not taken, so this case measures nothing: %v", err)
	}
	t.Cleanup(func() {
		conn.ExecContext(context.Background(), "ROLLBACK")
		conn.Close()
		db.Close()
	})
}
