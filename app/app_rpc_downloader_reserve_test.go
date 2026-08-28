// app/app_rpc_downloader_reserve_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/models"
)

func TestThunderTop3CandidateURLs(t *testing.T) {
	// 1. 画像テスト
	rawURL := "https://pbs.twimg.com/media/GX3sE17aAAA4f4g.jpg"
	candidates := BuildThunderTop3CandidateURLs(rawURL)

	if len(candidates) != 3 {
		t.Fatalf("expected 3 candidates, got %d", len(candidates))
	}
	if candidates[0].Type != models.ResolutionOrig || candidates[0].URL != "https://pbs.twimg.com/media/GX3sE17aAAA4f4g?format=jpg&name=orig" {
		t.Errorf("candidate[0] mismatch: %+v", candidates[0])
	}
	if candidates[2].Type != models.ResolutionWaybackOrig || candidates[2].URL != "https://web.archive.org/web/2id_/https://pbs.twimg.com/media/GX3sE17aAAA4f4g?format=jpg&name=orig" {
		t.Errorf("candidate[2] mismatch: %+v", candidates[2])
	}

	// 2. 動画テスト (?tag=14, ?tag=12, wayback ?tag=14)
	videoURL := "https://video-s.twimg.com/ext_tw_video/1807739525416001536/pu/vid/avc1/720x1280/P5HCLz4CxtcJZRIv.mp4"
	vCandidates := BuildThunderTop3CandidateURLs(videoURL)
	if len(vCandidates) != 3 {
		t.Fatalf("expected 3 video candidates, got %d", len(vCandidates))
	}
	if vCandidates[0].URL != videoURL+"?tag=14" {
		t.Errorf("video candidate[0] mismatch: expected %s, got %s", videoURL+"?tag=14", vCandidates[0].URL)
	}
	if vCandidates[1].URL != videoURL+"?tag=12" {
		t.Errorf("video candidate[1] mismatch: expected %s, got %s", videoURL+"?tag=12", vCandidates[1].URL)
	}
	if vCandidates[2].URL != "https://web.archive.org/web/2id_/"+videoURL+"?tag=14" {
		t.Errorf("video candidate[2] mismatch: expected %s, got %s", "https://web.archive.org/web/2id_/"+videoURL+"?tag=14", vCandidates[2].URL)
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
