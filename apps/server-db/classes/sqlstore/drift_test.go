package sqlstore_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/sorglos-it/nextgencloud/apps/server-db/classes/sqlstore"
	"github.com/sorglos-it/nextgencloud/shared/dbapi"
)

// A migrated database has no drift, and a bent column is the one finding -
// the second half is the effect check: a CheckColumns that answered nothing
// at all would pass the first.
func TestCheckColumnsFindsTheBentColumn(t *testing.T) {
	s, db := openDB(t)
	drift, err := sqlstore.CheckColumns(ctx, db, sqlstore.SQLite)
	if err != nil || len(drift) != 0 {
		t.Fatalf("fresh database: %v, %v", drift, err)
	}
	value := bytes.Repeat([]byte{0x7d}, 30)
	if err := s.OnceCreate(ctx, dbapi.OnceTask, value, "u-admin", onceGoodUntil); err != nil {
		t.Fatal(err)
	}
	bendExpires(t, db, dbapi.OnceTask, value)
	drift, err = sqlstore.CheckColumns(ctx, db, sqlstore.SQLite)
	want := []sqlstore.Drift{{Table: "ngc_once", Column: "expires", Want: "INTEGER", Found: "BLOB"}}
	if err != nil || !reflect.DeepEqual(drift, want) {
		t.Fatalf("bent: %v, %v; want %v", drift, err, want)
	}
}
