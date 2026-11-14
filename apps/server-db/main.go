// Command ngc-db runs on the database servers of a MariaDB cluster: it backs
// up and restores the local node with mariabackup.
package main

import (
	"os"

	"github.com/sorglos-it/nextgencloud/apps/server-db/classes/dbcli"
)

// Version is set when building.
var Version = "0.0.0-dev"

func main() { os.Exit(dbcli.Main(Version, os.Args[1:], os.Stdout, os.Stderr)) }
