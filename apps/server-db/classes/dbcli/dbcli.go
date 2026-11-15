// Package dbcli is cd-db: backup and restore of the local MariaDB node.
//
// It drives mariabackup and nothing else. A physical copy and no dump,
// because a dump played back may run code (register.backup says the same).
// The copy keeps the tables encrypted: mariabackup reads the key plugin from
// the server's own configuration, so the keys under /etc/cd/db-keys must be
// kept apart from the copy, or it can never be opened again.
//
// SQLite installations are not this program's: they are backed up with
// "cd-api backup", since their file lies on the API's machine.
package dbcli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MaxPath bounds the directory argument; mariabackup gets it verbatim.
const MaxPath = 1024

// run starts mariabackup with an argument list, never through a shell. Tests
// replace it.
var run = func(ctx context.Context, stdout, stderr io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "mariabackup", args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	return cmd.Run()
}

const usage = `Aufruf:
  cd-db backup <Ordner>    sichert den MariaDB-Knoten dieses Rechners in einen neuen Ordner
  cd-db restore <Ordner>   spielt eine Sicherung zurück (MariaDB gestoppt, Datenordner leer)
  cd-db version
`

// Main is the whole program; it answers the exit code.
func Main(version string, args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && args[0] == "version" {
		fmt.Fprintf(stdout, "cd-db %s\n", version)
		return 0
	}
	if len(args) != 2 || (args[0] != "backup" && args[0] != "restore") {
		fmt.Fprint(stderr, usage)
		return 2
	}
	dir, err := checkDir(args[1])
	if err != nil {
		fmt.Fprintf(stderr, "cd-db: %v\n", err)
		return 2
	}
	ctx := context.Background()
	if args[0] == "backup" {
		err = Backup(ctx, dir, stdout, stderr)
	} else {
		err = Restore(ctx, dir, stdout, stderr)
	}
	if err != nil {
		fmt.Fprintf(stderr, "cd-db: %v\n", err)
		return 1
	}
	return 0
}

// checkDir accepts only an absolute, clean path without control characters,
// so that what mariabackup gets is what the admin typed and nothing relative
// to wherever the program happened to start.
func checkDir(p string) (string, error) {
	if p == "" || len(p) > MaxPath {
		return "", errors.New("der Ordner fehlt oder ist zu lang")
	}
	if strings.ContainsFunc(p, func(r rune) bool { return r < 0x20 || r == 0x7f }) {
		return "", errors.New("der Ordner enthält Steuerzeichen")
	}
	// filepath.IsAbs would accept "C:\..." on Windows; the target is Linux.
	if !strings.HasPrefix(p, "/") || filepath.ToSlash(filepath.Clean(p)) != p || p == "/" {
		return "", fmt.Errorf("der Ordner %q muss ein absoluter Pfad ohne \"..\" sein", p)
	}
	return p, nil
}

// Backup copies the running node into dir, which must not exist yet: a backup
// written over an older one would be neither.
func Backup(ctx context.Context, dir string, stdout, stderr io.Writer) error {
	if _, err := os.Lstat(dir); err == nil {
		return fmt.Errorf("%s ist schon da; eine Sicherung kommt in einen neuen Ordner", dir)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// --galera-info records the cluster position, which a node restored from
	// this copy needs to rejoin without a full transfer.
	if err := run(ctx, stdout, stderr, "--backup", "--galera-info", "--target-dir="+dir); err != nil {
		return fmt.Errorf("mariabackup --backup: %w", err)
	}
	// Prepared right away, so a copy that cannot be made consistent fails now
	// and not on the day it is needed.
	if err := run(ctx, stdout, stderr, "--prepare", "--target-dir="+dir); err != nil {
		return fmt.Errorf("mariabackup --prepare: %w", err)
	}
	fmt.Fprintf(stdout, "Sicherung in %s geschrieben.\n", dir)
	return nil
}

// Restore copies a prepared backup back into the data directory. It does not
// stop or start MariaDB: mariabackup --copy-back refuses a non-empty data
// directory itself, and that refusal is the guard against overwriting a
// live node.
func Restore(ctx context.Context, dir string, stdout, stderr io.Writer) error {
	st, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("die Sicherung %s ist nicht da: %w", dir, err)
	}
	if !st.IsDir() {
		return fmt.Errorf("%s ist kein Ordner", dir)
	}
	if err := run(ctx, stdout, stderr, "--copy-back", "--target-dir="+dir); err != nil {
		return fmt.Errorf("mariabackup --copy-back: %w", err)
	}
	fmt.Fprintln(stdout, "Sicherung zurückgespielt. Dateirechte des Datenordners auf mysql setzen, dann MariaDB starten.")
	return nil
}
