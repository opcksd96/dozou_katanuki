package twitter_renderer

import (
	"fmt"
	"strings"

	"dozou_katanuki/models"
)

type TwitterRenderer struct{}

func NewTwitterRenderer() *TwitterRenderer {
	return &TwitterRenderer{}
}

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

// RenderTweet は Twitter 生記事モデルをフロントエンド用 RenderTree に装飾補正します
func (r *TwitterRenderer) RenderTweet(art *models.Article) *models.RenderTree {
	if art == nil {
		return nil
	}

	mediaItems := make([]models.RenderMedia, 0, len(art.Media))
	for _, m := range art.Media {
		var mediaURLs models.RenderMediaURLs
		mediaURLs.Original = m.DownloadURL

		effectiveStatus := m.DownloadStatus
		if m.DownloadStatus == "COMPLETED" && (m.StashSceneID.Valid || m.StashImageID.Valid) {
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = fmt.Sprintf("/stash-proxy/scene/%s/stream", m.StashSceneID.String)
			} else {
				mediaURLs.Image = fmt.Sprintf("/stash-proxy/image/%s/image", m.StashImageID.String)
				mediaURLs.Thumbnail = fmt.Sprintf("/stash-proxy/image/%s/thumbnail", m.StashImageID.String)
			}
		} else {
			if effectiveStatus == "COMPLETED" {
				effectiveStatus = "DEAD_404"
			}
		}

		mediaItems = append(mediaItems, models.RenderMedia{
			ID:             m.MediaID,
			Type:           m.Type,
			DownloadStatus: effectiveStatus,
			FailedReason:   m.FailedReason.String,
			URLs:           mediaURLs,
			Width:          m.Width,
			Height:         m.Height,
			StashSceneID:   m.StashSceneID.String,
			StashImageID:   m.StashImageID.String,
		})
	}

	author := models.RenderAuthor{
		NumericID:   art.Account.NumericID,
		Handle:      art.Account.Username,
		DisplayName: art.Account.DisplayName,
		AvatarURL:   art.Account.AvatarURL,
	}

	orig := expandRedirects(art.FullText, art.UrlRedirects)
	ja := expandRedirects(art.FullTextJA.String, art.UrlRedirects)
	en := expandRedirects(art.FullTextEN.String, art.UrlRedirects)
	zh := expandRedirects(art.FullTextZH.String, art.UrlRedirects)

	return &models.RenderTree{
		ID:             art.ID,
		ConversationID: art.ConversationID,
		CreatedAt:      art.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Content: models.RenderContent{
			Original: orig,
			JA:       ja,
			EN:       en,
			ZH:       zh,
		},
		Author:    author,
		Media:     mediaItems,
		IsLiked:   art.IsLiked,
		SourceURL: art.WaybackURL,
	}
}

