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

func TestToRenderTree_MediaSanitization(t *testing.T) {
	art := models.Article{
		ID: "1002",
		Media: []models.Media{
			{
				MediaID:        "vid_ok",
				Type:           "video",
				DownloadURL:    "https://video.twimg.com/amplify_video/123.mp4",
				DownloadStatus: "COMPLETED",
				StashSceneID:   sql.NullString{String: "scene_123", Valid: true},
			},
			{
				MediaID:        "vid_missing_stash",
				Type:           "video",
				DownloadURL:    "https://video.twimg.com/amplify_video/456.mp4",
				DownloadStatus: "COMPLETED",
			},
			{
				MediaID:        "vid_dead",
				Type:           "video",
				DownloadURL:    "https://video.twimg.com/amplify_video/789.mp4",
				DownloadStatus: "DEAD_404",
			},
		},
	}
	tree := ToRenderTree(art, "twitter")
	if len(tree.Media) != 3 {
		t.Fatalf("expected 3 media items, got %d", len(tree.Media))
	}
	if tree.Media[0].URLs.Stream != "/stash-proxy/scene/scene_123/stream" || tree.Media[0].DownloadStatus != "COMPLETED" {
		t.Errorf("expected media[0] stream proxy, got %v", tree.Media[0])
	}
	if tree.Media[1].URLs.Stream != "" || tree.Media[1].DownloadStatus != "DEAD_404" {
		t.Errorf("expected media[1] stream to be empty & status DEAD_404, got %v", tree.Media[1])
	}
	if tree.Media[2].URLs.Stream != "" || tree.Media[2].DownloadStatus != "DEAD_404" {
		t.Errorf("expected media[2] stream to be empty & status DEAD_404, got %v", tree.Media[2])
	}
}
