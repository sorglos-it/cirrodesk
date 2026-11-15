//go:build !linux

package node

// CirroDesk is built for Linux alone; elsewhere nothing is locked, so the tests of
// this package also run where they are written.
func lockRun(string) (func(), error) { return func() {}, nil }
