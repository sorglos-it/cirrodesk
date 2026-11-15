// Command cd-files is the service of the store app org.sorglos.files. The
// controller starts it once per call ("serve"); "hold" keeps the container up
// between calls, as cd-controller hold does for the terminal.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sorglos-it/cirrodesk/apps/store-files/classes"
	"github.com/sorglos-it/cirrodesk/shared/service"
)

// Version is set when building.
var Version = "0.0.0-dev"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: cd-files serve|hold|version")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "serve":
		os.Exit(classes.Run(os.Getenv, service.DataDir, os.Stdin, os.Stdout, os.Stderr))
	case "hold":
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, os.Interrupt)
		<-c
	case "version":
		fmt.Println("cd-files", Version)
	default:
		fmt.Fprintln(os.Stderr, "usage: cd-files serve|hold|version")
		os.Exit(2)
	}
}
