package appstore

import "testing"

// TestRefreshNowFetchesAtMostOnceAMinute holds A25: the button fetches every
// store at once, and a second press within RefreshEvery gets the last state
// without knocking again.
func TestRefreshNowFetchesAtMostOnceAMinute(t *testing.T) {
	fx, v := newFixture(t)
	rows, err := fx.St.RefreshNow(ctx)
	if err != nil {
		t.Fatalf("RefreshNow: %v", err)
	}
	first := fx.hits()
	if first == 0 {
		t.Fatal("the store saw no request")
	}
	found := false
	for _, r := range rows {
		if r.ID == v.Store && r.Seq == v.CatalogLastSeq {
			found = true
		}
	}
	if !found {
		t.Fatalf("the answer does not carry the fetched catalog: %+v", rows)
	}
	if _, err := fx.St.RefreshNow(ctx); err != nil {
		t.Fatalf("second RefreshNow: %v", err)
	}
	if fx.hits() != first {
		t.Fatalf("a second press within the minute fetched again: %d requests, want %d", fx.hits(), first)
	}
	// Once the minute is over, a press fetches again.
	fx.St.manual = fx.St.manual.Add(-RefreshEvery)
	if _, err := fx.St.RefreshNow(ctx); err != nil {
		t.Fatalf("third RefreshNow: %v", err)
	}
	if fx.hits() == first {
		t.Fatal("a press after the minute did not fetch")
	}
}
