package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sorglos-it/cirrodesk/apps/server-api/classes/join"
	"github.com/sorglos-it/cirrodesk/shared/session"
	"github.com/sorglos-it/cirrodesk/shared/tunnel"
)

// loggedPaths are the only paths a request line names. Anything else is
// "other": the path is chosen by whoever sends the request, and the journal
// must not be a place a stranger writes into.
var loggedPaths = map[string]bool{
	session.PathSetup: true, session.PathLoginStart: true, session.PathLoginFinish: true,
	session.PathLogout: true, session.PathTrust: true, session.PathSocket: true,
	tunnel.PathTunnel: true, tunnel.PathPair: true, tunnel.PathJoin: true, join.PathKeys: true,
}

// requestLog writes one line per request (A8c): with two API servers behind
// one balancer, an error could not be traced to the node that answered,
// because no line existed for a single request. The line names the node, the
// method, the path from the fixed list, the status and the time. Never the
// query (PathJoin carries the hash of a one-time secret there), never a
// header, never a body. A WebSocket's line comes when it ends, with status
// 0 when the library wrote its 101 on the hijacked connection itself.
func requestLog(node string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		route := r.URL.EscapedPath()
		if !loggedPaths[route] {
			route = "other"
		}
		method := r.Method
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodDelete,
			http.MethodOptions, http.MethodPatch:
		default:
			method = "other"
		}
		slog.Info("request", "node", node, "method", method, "path", route, "status", rec.status,
			"ms", time.Since(start).Milliseconds())
	})
}

// statusRecorder keeps the status a handler sent. Unwrap lets the WebSocket
// library and http.ResponseController reach the Hijacker and Flusher below.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	if s.status == 0 {
		s.status = code
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }
