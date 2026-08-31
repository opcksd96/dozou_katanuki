// app/app_media_candidates.go (100行以下 - SPEC-PRINCIPLE-001)
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

// BuildMediaCandidateURLs はメディアの全候補URL (orig, colon, plain, large, 各wayback) を網羅生成します
func BuildMediaCandidateURLs(rawURL string) []CandidateURL {
	if rawURL == "" {
		return nil
	}

	base := rawURL
	for _, sfx := range []string{":orig", ":large", ":medium", ":small", ":thumb", ":tiny", ":900x900", ":1200x1200"} {
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
		tag14URL := cleanBase + "?tag=14"
		tag12URL := cleanBase + "?tag=12"
		return []CandidateURL{
			{Type: models.ResolutionOrig, URL: tag14URL},
			{Type: models.ResolutionLarge, URL: tag12URL},
			{Type: models.ResolutionPlain, URL: cleanBase},
			{Type: models.ResolutionWaybackOrig, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", tag14URL)},
			{Type: models.ResolutionWaybackPlain, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", cleanBase)},
		}
	}

	cleanNoExt := cleanBase
	ext := ".jpg"
	for _, e := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		if strings.HasSuffix(strings.ToLower(cleanNoExt), e) {
			ext = e
			cleanNoExt = cleanNoExt[:len(cleanNoExt)-len(e)]
			break
		}
	}

	plainURL := cleanNoExt + ext
	colonOrigURL := plainURL + ":orig"
	paramOrigURL := fmt.Sprintf("%s?format=%s&name=orig", cleanNoExt, strings.TrimPrefix(ext, "."))
	paramLargeURL := fmt.Sprintf("%s?format=%s&name=large", cleanNoExt, strings.TrimPrefix(ext, "."))

	return []CandidateURL{
		{Type: models.ResolutionOrig, URL: paramOrigURL},
		{Type: models.ResolutionColonOrig, URL: colonOrigURL},
		{Type: models.ResolutionPlain, URL: plainURL},
		{Type: models.ResolutionWaybackPlain, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", plainURL)},
		{Type: models.ResolutionWaybackColon, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", colonOrigURL)},
		{Type: models.ResolutionLarge, URL: paramLargeURL},
		{Type: models.ResolutionWaybackOrig, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", paramOrigURL)},
	}
}

// BuildThunderTop3CandidateURLs は後方互換性のためのエイリアスです
func BuildThunderTop3CandidateURLs(rawURL string) []CandidateURL {
	return BuildMediaCandidateURLs(rawURL)
}
