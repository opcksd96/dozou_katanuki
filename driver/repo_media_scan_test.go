package driver

import (
	"testing"
)

func TestSearchMediaDetails_RealDB(t *testing.T) {
	db, err := InitDB("../archive.db")
	if err != nil {
		t.Skipf("archive.db not found: %v", err)
	}
	repo := NewRepository(db)
	res, err := repo.SearchMediaDetails("all", "all", "all", 24, 0)
	if err != nil {
		t.Fatalf("SearchMediaDetails failed: %v", err)
	}
	if len(res.Items) == 0 {
		t.Fatalf("expected items, got 0")
	}
	t.Logf("Fetched %d items, Total: %d, Stats: %+v", len(res.Items), res.Total, res.Stats)
	t.Logf("First item: MediaID=%s, ArticleID=%s, Username=%s, DisplayName=%s",
		res.Items[0].MediaID, res.Items[0].ArticleID, res.Items[0].Username, res.Items[0].DisplayName)

	resAcc, err := repo.SearchMediaDetails("msluo14", "all", "all", 24, 0)
	if err != nil {
		t.Fatalf("SearchMediaDetails with account failed: %v", err)
	}
	t.Logf("Fetched msluo14 items: %d, Total: %d", len(resAcc.Items), resAcc.Total)
}
