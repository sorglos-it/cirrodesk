package server

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/shared/trust"
)

func TestAbsenceSaysOnceGoneAndOnceBack(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)

	nodes := []trust.Node{{ID: "ctl-1", Kind: trust.NodeController}, {ID: "api-1", Kind: trust.NodeAPI}}
	up := false
	var a absence
	round := func() { a.check(nodes, func(id string) bool { return up }) }

	round()
	if buf.Len() != 0 {
		t.Fatalf("one round without a tunnel is a restart, not a loss: %q", buf.String())
	}
	round()
	round()
	round()
	if n := strings.Count(buf.String(), "has no tunnel"); n != 1 {
		t.Fatalf("missing said %d times, want 1: %q", n, buf.String())
	}
	if !strings.Contains(buf.String(), "node=ctl-1") || strings.Contains(buf.String(), "api-1") {
		t.Fatalf("wrong node named: %q", buf.String())
	}
	up = true
	round()
	round()
	if n := strings.Count(buf.String(), "is back"); n != 1 {
		t.Fatalf("back said %d times, want 1: %q", n, buf.String())
	}

	// Gone for one round and back: nothing to say either way.
	buf.Reset()
	up = false
	round()
	up = true
	round()
	if buf.Len() != 0 {
		t.Fatalf("a short gap was reported: %q", buf.String())
	}

	// Taken out of the list while missing: forgotten, not reported back.
	up = false
	round()
	round()
	buf.Reset()
	nodes = nodes[1:]
	round()
	if buf.Len() != 0 || len(a.rounds) != 0 {
		t.Fatalf("an unlisted node is still tracked: %q %v", buf.String(), a.rounds)
	}
}
