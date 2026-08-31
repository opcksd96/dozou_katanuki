// app/app_rpc_downloader_reserve_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/models"
)

func TestMediaCandidateURLs(t *testing.T) {
	// 1. 画像テスト (プレーン、コロンorig、クエリorig、Waybackプレーン等を網羅)
	rawURL := "https://pbs.twimg.com/media/GX3sE17aAAA4f4g.jpg"
	candidates := BuildMediaCandidateURLs(rawURL)

	if len(candidates) < 5 {
		t.Fatalf("expected at least 5 candidates, got %d", len(candidates))
	}
	urlMap := make(map[string]models.ThunderResolutionType)
	for _, c := range candidates {
		urlMap[c.URL] = c.Type
	}

	// プレーンURL
	if _, ok := urlMap["https://pbs.twimg.com/media/GX3sE17aAAA4f4g.jpg"]; !ok {
		t.Errorf("missing plain URL: %+v", urlMap)
	}
	// コロン :orig
	if _, ok := urlMap["https://pbs.twimg.com/media/GX3sE17aAAA4f4g.jpg:orig"]; !ok {
		t.Errorf("missing colon :orig URL: %+v", urlMap)
	}
	// Wayback プレーン
	if _, ok := urlMap["https://web.archive.org/web/2id_/https://pbs.twimg.com/media/GX3sE17aAAA4f4g.jpg"]; !ok {
		t.Errorf("missing wayback plain URL: %+v", urlMap)
	}
	// クエリ orig
	if _, ok := urlMap["https://pbs.twimg.com/media/GX3sE17aAAA4f4g?format=jpg&name=orig"]; !ok {
		t.Errorf("missing query orig URL: %+v", urlMap)
	}

	// 2. 動画テスト
	videoURL := "https://video-s.twimg.com/ext_tw_video/1807739525416001536/pu/vid/avc1/720x1280/P5HCLz4CxtcJZRIv.mp4"
	vCandidates := BuildMediaCandidateURLs(videoURL)
	if len(vCandidates) < 3 {
		t.Fatalf("expected at least 3 video candidates, got %d", len(vCandidates))
	}
}

func TestSafePurgeEmpty(t *testing.T) {
	app := &App{}
	count, err := app.SafePurgeWithBackup([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 purged, got %d", count)
	}
}
