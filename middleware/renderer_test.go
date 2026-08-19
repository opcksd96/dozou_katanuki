package middleware

import (
	"database/sql"
	"testing"
	"time"

	"dozou_katanuki/models"
)

func TestToRenderTree_ExpandRedirects(t *testing.T) {
	art := models.Article{
		ID:             "1001",
		ConversationID: "1001",
		CreatedAt:      time.Now(),
		FullText:       "Check out this site! https://t.co/xyz789 and https://t.co/abc123",
		FullTextJA:     sql.NullString{String: "このサイトを見て！ https://t.co/xyz789", Valid: true},
		Account: models.Account{
			NumericID:   "9999",
			Username:    "mash_kyrielight",
			DisplayName: "Mash",
		},
		UrlRedirects: []models.UrlRedirect{
			{
				ShortURL:    "https://t.co/xyz789",
				ExpandedURL: "https://example.com/nes-apu-manual",
				ArticleID:   "1001",
			},
			{
				ShortURL:    "https://t.co/abc123",
				ExpandedURL: "https://example.com/famitracker",
				ArticleID:   "1001",
			},
		},
	}

	renderTree := ToRenderTree(art, "twitter")

	// Original の検証
	expectedOrig := "Check out this site! https://example.com/nes-apu-manual and https://example.com/famitracker"
	if renderTree.Content.Original != expectedOrig {
		t.Errorf("expected Content.Original = %q, got %q", expectedOrig, renderTree.Content.Original)
	}

	// JA の検証
	expectedJA := "このサイトを見て！ https://example.com/nes-apu-manual"
	if renderTree.Content.JA != expectedJA {
		t.Errorf("expected Content.JA = %q, got %q", expectedJA, renderTree.Content.JA)
	}
}
