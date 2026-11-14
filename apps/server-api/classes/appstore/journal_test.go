package appstore

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

// TestJournalSaysAChangeOnce holds A5: a store that keeps failing leaves one
// journal line and not one per round, a recovery leaves one more, and a store
// that never failed leaves none.
func TestJournalSaysAChangeOnce(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(old) })
	st := &Stores{}
	st.journal("good", nil)
	for range 3 {
		st.journal("shop", ErrStore)
	}
	st.journal("shop", nil)
	st.journal("shop", nil)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], "level=WARN") || !strings.Contains(lines[0], "store=shop") ||
		!strings.Contains(lines[1], "answers again") {
		t.Fatalf("journal said:\n%s", buf.String())
	}
	// The line carries the fixed sentence, never the error's own text.
	buf.Reset()
	st.journal("other", errors.Join(ErrStore, errors.New("https://x.test/?secret=1")))
	if strings.Contains(buf.String(), "secret") {
		t.Fatalf("the error text reached the journal: %s", buf.String())
	}
}
