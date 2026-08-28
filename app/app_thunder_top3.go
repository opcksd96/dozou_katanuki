// app/app_thunder_top3.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"strings"

	"dozou_katanuki/models"
)

// CandidateURL は解像度種別付きのURL構造です
type CandidateURL struct {
	Type models.ThunderResolutionType
	URL  string
}

// BuildThunderTop3CandidateURLs は高成功率の厳選3種類 (orig, large, wayback_orig) のURLを生成します
func BuildThunderTop3CandidateURLs(rawURL string) []CandidateURL {
	if rawURL == "" {
		return nil
	}

	base := rawURL
	for _, sfx := range []string{":orig", ":large", ":medium", ":small", ":thumb", ":tiny"} {
		if strings.HasSuffix(base, sfx) {
			base = strings.TrimSuffix(base, sfx)
			break
		}
	}
	cleanBase := base
	if idx := strings.Index(base, "?"); idx != -1 {
		cleanBase = base[:idx]
	}

	isVideo := strings.Contains(base, "video") || strings.Contains(base, ".mp4") || strings.Contains(base, ".m3u8")

	if isVideo {
		// 動画の場合: ?tag=14 (最高画質), ?tag=12 (標準画質), Waybackアーカイブ
		tag14URL := cleanBase + "?tag=14"
		tag12URL := cleanBase + "?tag=12"
		return []CandidateURL{
			{Type: models.ResolutionOrig, URL: tag14URL},
			{Type: models.ResolutionLarge, URL: tag12URL},
			{Type: models.ResolutionWaybackOrig, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", tag14URL)},
		}
	}

	cleanNoExt := cleanBase
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		if strings.HasSuffix(strings.ToLower(cleanNoExt), ext) {
			cleanNoExt = cleanNoExt[:len(cleanNoExt)-len(ext)]
			break
		}
	}

	var origURL, largeURL string
	if strings.Contains(base, "pbs.twimg.com") || strings.Contains(base, "/media/") {
		origURL = fmt.Sprintf("%s?format=jpg&name=orig", cleanNoExt)
		largeURL = fmt.Sprintf("%s?format=jpg&name=large", cleanNoExt)
	} else {
		origURL = cleanBase + ":orig"
		largeURL = cleanBase + ":large"
	}

	return []CandidateURL{
		{Type: models.ResolutionOrig, URL: origURL},
		{Type: models.ResolutionLarge, URL: largeURL},
		{Type: models.ResolutionWaybackOrig, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", origURL)},
	}
}
