package server

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sorglos-it/cirrodesk/shared/session"
	"github.com/sorglos-it/cirrodesk/shared/tunnel"
)

// TestRequestLogNamesNodePathAndStatusOnly holds A8c: one line per request
// with node, a path from the fixed list, the status -- and no query, and no
// path a stranger chose.
func TestRequestLogNamesNodePathAndStatusOnly(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(old) })
	h := requestLog("api-2", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, tunnel.PathJoin+"?secret=SECRETHASH", nil))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/evil%0Aline", nil))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, session.PathTrust, nil))
	got := buf.String()
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got:\n%s", got)
	}
	if !strings.Contains(lines[0], "node=api-2") || !strings.Contains(lines[0], "path="+tunnel.PathJoin) ||
		!strings.Contains(lines[0], "status=403") {
		t.Fatalf("first line: %s", lines[0])
	}
	if strings.Contains(got, "SECRETHASH") || strings.Contains(got, "evil") {
		t.Fatalf("query or stranger's path reached the journal:\n%s", got)
	}
	if !strings.Contains(lines[1], "path=other") || !strings.Contains(lines[2], "path="+session.PathTrust) {
		t.Fatalf("paths: %s", got)
	}
}
