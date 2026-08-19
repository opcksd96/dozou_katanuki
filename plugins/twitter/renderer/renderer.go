// plugins/twitter/renderer/renderer.go (100行以下)
package twitter_renderer

import (
	"fmt"

	"dozou_katanuki/models"
)

type TwitterRenderer struct{}

func NewTwitterRenderer() *TwitterRenderer {
	return &TwitterRenderer{}
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

		if m.DownloadStatus == "COMPLETED" && (m.StashSceneID.Valid || m.StashImageID.Valid) {
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = fmt.Sprintf("/stash-proxy/scene/%s/stream", m.StashSceneID.String)
			} else {
				mediaURLs.Image = fmt.Sprintf("/stash-proxy/image/%s/image", m.StashImageID.String)
				mediaURLs.Thumbnail = fmt.Sprintf("/stash-proxy/image/%s/thumbnail", m.StashImageID.String)
			}
		} else {
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = m.DownloadURL
			} else {
				mediaURLs.Image = m.DownloadURL
				mediaURLs.Thumbnail = m.DownloadURL
			}
		}

		mediaItems = append(mediaItems, models.RenderMedia{
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

	author := models.RenderAuthor{
		NumericID:   art.Account.NumericID,
		Handle:      art.Account.Username,
		DisplayName: art.Account.DisplayName,
		AvatarURL:   art.Account.AvatarURL,
	}

	return &models.RenderTree{
		ID:             art.ID,
		ConversationID: art.ConversationID,
		CreatedAt:      art.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Content: models.RenderContent{
			Original: art.FullText,
			JA:       art.FullTextJA.String,
			EN:       art.FullTextEN.String,
			ZH:       art.FullTextZH.String,
		},
		Author:    author,
		Media:     mediaItems,
		IsLiked:   art.IsLiked,
		SourceURL: art.WaybackURL,
	}
}

