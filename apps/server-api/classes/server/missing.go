package server

import (
	"log/slog"

	"github.com/sorglos-it/cirrodesk/shared/trust"
)

// missingAfter is how many minute rounds in a row a controller of the trust
// list may be without a tunnel to this API before the journal says so. One
// round is not enough: a controller that restarts, or an API that has just
// started, is without a tunnel for a moment, and a line for that would be the
// noise that teaches people to stop reading the log.
const missingAfter = 2

// absence answers A8a ("a part can drop out without anybody seeing it"): with
// the controller of svc-1 off, every app was gone and the installation still
// reported 40 good logins out of 40. Every controller keeps a tunnel to EVERY
// API (server-controller node/run.go), so each API can tell on its own which
// one is missing without asking its neighbours.
//
// It says a node once when it goes missing and once when it is back, never
// every minute. Only Server.minute uses it, under Server.once, so it needs no
// lock of its own.
type absence struct {
	rounds map[string]int
}

// check takes one minute round: nodes is the trust list's, up says whether a
// node has a tunnel to this API right now.
func (a *absence) check(nodes []trust.Node, up func(string) bool) {
	if a.rounds == nil {
		a.rounds = map[string]int{}
	}
	listed := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		if n.Kind != trust.NodeController {
			continue
		}
		listed[n.ID] = true
		if up(n.ID) {
			if a.rounds[n.ID] >= missingAfter {
				slog.Info("controller is back", "node", n.ID)
			}
			delete(a.rounds, n.ID)
			continue
		}
		a.rounds[n.ID]++
		if a.rounds[n.ID] == missingAfter {
			slog.Warn("controller has no tunnel to this API", "node", n.ID)
		}
	}
	// A node taken out of the trust list is not missing; it is gone on
	// purpose, and the next list that names it again starts from nothing.
	for id := range a.rounds {
		if !listed[id] {
			delete(a.rounds, id)
		}
	}
}
