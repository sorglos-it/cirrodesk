package login

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sorglos-it/nextgencloud/apps/server-api/classes/state"
	"github.com/sorglos-it/nextgencloud/shared/dbapi/memstore"
	"github.com/sorglos-it/nextgencloud/shared/session"
	"github.com/sorglos-it/nextgencloud/shared/wire"
)

// TestALoginBeforeTheSetupSaysNotSetUp is A3: an installation that is migrated
// and never set up answered every login with E_UNREACHABLE, about a database
// that answers perfectly well. Both paths read the list before any name, so
// the answer is the same for a name that exists and one that does not.
func TestALoginBeforeTheSetupSaysNotSetUp(t *testing.T) {
	full, fx := newFixture(t)
	// A real finish body, made on a set-up installation: one that fails to
	// parse would be refused before the list is read and test nothing here.
	finish := string(finishBody(t, fx, fx.Password, start(t, full, fx.PasswordName), ""))
	s := state.New(full.S.Cfg, memstore.New(), full.S.Box, full.S.Pin)
	s.Now = full.S.Now
	s.AnchorFile = filepath.Join(t.TempDir(), "trust-anchor.json")
	h := New(s)
	for _, c := range []struct {
		name string
		call func(http.ResponseWriter, *http.Request)
		path string
		body string
	}{
		{"start", h.Start, session.PathLoginStart, `{"name":"anna"}`},
		{"start other name", h.Start, session.PathLoginStart, `{"name":"nobody"}`},
		{"finish", h.Finish, session.PathLoginFinish, finish},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.call(w, jsonPost(c.path, c.body))
			if w.Code != session.HTTPStatus(wire.CodeNotSetUp) {
				t.Fatalf("status = %d, want %d: %s", w.Code, session.HTTPStatus(wire.CodeNotSetUp), w.Body)
			}
			if !strings.Contains(w.Body.String(), `"`+string(wire.CodeNotSetUp)+`"`) {
				t.Fatalf("body = %s", w.Body)
			}
		})
	}
}
