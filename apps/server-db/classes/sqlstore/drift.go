package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Drift is one column whose type in the database is not the one in
// ColumnType. It carries names and types only, never a value: a value from
// cd_session or cd_user is a secret.
type Drift struct {
	Table, Column, Want, Found string
}

// CheckColumns reads the real type of every column in ColumnType and answers
// those that differ (docs/schema-pruefung-beim-start.md, Weg B). A column
// that is not there is not reported: an API that starts before the migration
// that adds it is the rolling update, not a fault.
func CheckColumns(ctx context.Context, db *sql.DB, d Dialect) ([]Drift, error) {
	var out []Drift
	tables := make([]string, 0, len(ColumnType))
	for t := range ColumnType {
		tables = append(tables, t)
	}
	slices.Sort(tables)
	for _, t := range tables {
		found, err := columnTypes(ctx, db, d, t)
		if err != nil {
			return nil, err
		}
		cols := make([]string, 0, len(ColumnType[t]))
		for c := range ColumnType[t] {
			cols = append(cols, c)
		}
		slices.Sort(cols)
		for _, c := range cols {
			got, ok := found[c]
			if !ok {
				continue
			}
			if want := expectedType(d, ColumnType[t][c]); got != want {
				out = append(out, Drift{t, c, want, got})
			}
		}
	}
	return out, nil
}

// expectedType is the type as the database reports it back. MariaDB: the
// measured rule of the design - DATA_TYPE in lower case, the length only for
// varbinary, because BIGINT comes back as bigint(20) and the BLOB sorts carry
// a length nobody wrote.
func expectedType(d Dialect, placeholder string) string {
	t := d.Type(placeholder)
	if d == MariaDB {
		return strings.ToLower(t)
	}
	return t
}

// columnTypes answers column name -> type of one table. t is always a key of
// ColumnType, never input from outside, and goes in as a placeholder anyway.
func columnTypes(ctx context.Context, db *sql.DB, d Dialect, t string) (map[string]string, error) {
	const (
		maria = "SELECT COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH FROM information_schema.COLUMNS " +
			"WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?"
		lite = "SELECT name, type, NULL FROM pragma_table_info(?)"
	)
	var rows *sql.Rows
	var err error
	if d == MariaDB {
		rows, err = db.QueryContext(ctx, maria, t)
	} else {
		rows, err = db.QueryContext(ctx, lite, t)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlstore: column types of %s: %w", t, err)
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var name, typ string
		var length sql.NullInt64
		if err := rows.Scan(&name, &typ, &length); err != nil {
			return nil, err
		}
		if d == MariaDB {
			typ = strings.ToLower(typ)
			if typ == "varbinary" && length.Valid {
				typ += "(" + strconv.FormatInt(length.Int64, 10) + ")"
			}
		} else {
			typ = strings.ToUpper(typ)
		}
		out[name] = typ
	}
	return out, rows.Err()
}

// LogDrift runs CheckColumns at a start and writes each finding as an ERROR.
// It never stops the start - that is Weg B, chosen so that a wrong check
// cannot take a healthy installation down; a check that cannot run is a
// WARN for the same reason.
func LogDrift(db *sql.DB, d Dialect) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	drift, err := CheckColumns(ctx, db, d)
	if err != nil {
		slog.Warn("db: cannot check the column types", "error", err)
		return
	}
	for _, x := range drift {
		slog.Error("db: Schema weicht ab", "table", x.Table, "column", x.Column, "want", x.Want, "found", x.Found)
	}
}
