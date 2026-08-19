package middleware

import (
	"fmt"
	"strings"

	"dozou_katanuki/models"
)

// expandRedirects は短縮URLマップを適用してテキスト内の短縮URLを展開します
func expandRedirects(text string, redirects []models.UrlRedirect) string {
	if text == "" || len(redirects) == 0 {
		return text
	}
	res := text
	for _, r := range redirects {
		if r.ShortURL != "" && r.ExpandedURL != "" && strings.Contains(res, r.ShortURL) {
			res = strings.ReplaceAll(res, r.ShortURL, r.ExpandedURL)
		}
	}
	return res
}

// ToRenderTree は装飾済みの中間データをフロントエンド描画用構造体へ即時マッピングします
func ToRenderTree(item models.Article, platform string) models.RenderTree {
	avatarURL := AuditAndResolveAvatar(platform, item.CreatedAt, item.Account.ProfileHistory)

	orig := expandRedirects(item.FullText, item.UrlRedirects)
	ja := expandRedirects(item.FullTextJA.String, item.UrlRedirects)
	en := expandRedirects(item.FullTextEN.String, item.UrlRedirects)
	zh := expandRedirects(item.FullTextZH.String, item.UrlRedirects)

	return models.RenderTree{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Content: models.RenderContent{
			Original: orig,
			JA:       ja,
			EN:       en,
			ZH:       zh,
		},
		Author: models.RenderAuthor{
			NumericID:   item.Account.NumericID,
			Handle:      item.Account.Username,
			DisplayName: item.Account.DisplayName,
			AvatarURL:   avatarURL,
		},
		Media: mapMediaToRenderMedia(item.Media),
		Metrics: models.RenderMetrics{
			Replies:  0,
			Retweets: 0,
			Likes:    0,
		},
		IsLiked:   item.IsLiked,
		IsPinned:  false,
		SourceURL: item.WaybackURL,
	}
}

func mapMediaToRenderMedia(mediaList []models.Media) []models.RenderMedia {
	result := make([]models.RenderMedia, 0, len(mediaList))
	for _, m := range mediaList {
		var mediaURLs models.RenderMediaURLs
		mediaURLs.Original = m.DownloadURL

		if m.DownloadStatus == "COMPLETED" && (m.StashSceneID.Valid || m.StashImageID.Valid) {
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = fmt.Sprintf("/stash-proxy/scene/%s/stream", m.StashSceneID.String)
			} else {
				mediaURLs.Image = fmt.Sprintf("/stash-proxy/image/%s/image", m.StashImageID.String)
				mediaURLs.Thumbnail = fmt.Sprintf("/stash-proxy/image/%s/thumbnail", m.StashImageID.String)
			}
		} else {
			// 未ダウンロード（Stash未登録）時は原本/Wayback外部URLを直引き
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = m.DownloadURL
			} else {
				mediaURLs.Image = m.DownloadURL
				mediaURLs.Thumbnail = m.DownloadURL
			}
		}

		result = append(result, models.RenderMedia{
			ID:             m.MediaID,
			Type:           m.Type,
			DownloadStatus: m.DownloadStatus,
			FailedReason:   m.FailedReason.String,
			URLs:           mediaURLs,
			Width:          m.Width,
			Height:         m.Height,
			StashSceneID:   m.StashSceneID.String,
			StashImageID:   m.StashImageID.String,
		})
	}
	return result
}
