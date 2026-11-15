package main

import (
	"os"

	"github.com/sorglos-it/cirrodesk/tools/release"
)

func main() {
	f, err := os.Create("tools/release/testdata/cd-web-all.tar")
	if err != nil {
		panic(err)
	}
	if err := release.WebFile(f); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}
