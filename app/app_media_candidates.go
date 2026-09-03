// app/app_media_candidates.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// CandidateURL は解像度種別付きのURL構造です
type CandidateURL struct {
	Type models.ThunderResolutionType
	URL  string
}

// BuildCandidateURLsFromMedia は 画像/動画の構造差を厳密に区別して候補URL群を網羅生成します
func BuildCandidateURLsFromMedia(mediaID, downloadURL, mediaType string) []CandidateURL {
	return BuildCandidateURLsFromMediaWithArticle(mediaID, downloadURL, mediaType, "")
}

// BuildCandidateURLsFromMediaWithArticle は article_id を活用した動画の正規階層トレイリング候補も生成します
func BuildCandidateURLsFromMediaWithArticle(mediaID, downloadURL, mediaType, articleID string) []CandidateURL {
	isVideo := mediaType == "video" || strings.Contains(downloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(mediaID), ".mp4") || strings.Contains(downloadURL, "video")
	cleanID := mediaID
	if cleanID == "" {
		cleanID = filepath.Base(downloadURL)
	}
	if idx := strings.Index(cleanID, "?"); idx != -1 {
		cleanID = cleanID[:idx]
	}
	ext := filepath.Ext(cleanID)
	rawID := strings.TrimSuffix(cleanID, ext)
	for _, sfx := range []string{":orig", ":large", ":medium", ":small", ":thumb", ":tiny", "_motrix", "_requests", "_thunder", "_plain", "_orig", "_large"} {
		rawID = strings.TrimSuffix(rawID, sfx)
	}

	if isVideo {
		if ext == "" {
			ext = ".mp4"
			cleanID = rawID + ext
		}
		if downloadURL != "" && strings.Contains(downloadURL, "http") {
			base := downloadURL
			if idx := strings.Index(base, "?"); idx != -1 {
				base = base[:idx]
			}
			tag14, tag12 := base+"?tag=14", base+"?tag=12"
			return []CandidateURL{
				{Type: models.ResolutionOrig, URL: tag14},
				{Type: models.ResolutionLarge, URL: tag12},
				{Type: models.ResolutionPlain, URL: base},
				{Type: models.ResolutionWaybackOrig, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", tag14)},
				{Type: models.ResolutionWaybackPlain, URL: fmt.Sprintf("https://web.archive.org/web/2id_/%s", base)},
			}
		}
		if articleID != "" {
			// [仕様変更 2026-09-02]: 解像度の推測ループを削除 (Plan B)
			// Twitterの動画は解像度ごとにファイル名(ハッシュ)が異なるため、同じcleanIDで解像度だけを
			// 変更してURLを生成しても100% 404になります。これが429の原因でした。
			// 正規のdownloadURLがない場合は推測を諦め、空リストを返します。
			return []CandidateURL{}
		}
		return []CandidateURL{}
	}

	if ext == "" {
		ext = ".jpg"
	}
	plainURL := fmt.Sprintf("https://pbs.twimg.com/media/%s%s", rawID, ext)
	colonOrigURL := plainURL + ":orig"
	paramOrigURL := fmt.Sprintf("https://pbs.twimg.com/media/%s?format=%s&name=orig", rawID, strings.TrimPrefix(ext, "."))
	paramLargeURL := fmt.Sprintf("https://pbs.twimg.com/media/%s?format=%s&name=large", rawID, strings.TrimPrefix(ext, "."))

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

// BuildMediaCandidateURLs は後方互換用のラッパーです
func BuildMediaCandidateURLs(rawURL string) []CandidateURL {
	fn := filepath.Base(rawURL)
	return BuildCandidateURLsFromMedia(fn, rawURL, "")
}
