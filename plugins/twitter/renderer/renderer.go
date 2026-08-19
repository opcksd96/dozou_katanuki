// plugins/twitter/renderer/renderer.go (100行以下)
package twitter_renderer

import (
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
		mediaItems = append(mediaItems, models.RenderMedia{
			ID:          m.ID,
			Type:        m.Type,
			DirectURL:   m.URL,
			LocalPath:   "/media-local/" + m.ID,
			StashStream: "/stash-proxy/scene/" + m.StashSceneID + "/stream",
			Width:       m.Width,
			Height:      m.Height,
			AspectRatio: calculateAspectRatio(m.Width, m.Height),
		})
	}

	author := models.RenderAuthor{
		ID:          art.Account.ID,
		Handle:      art.Account.Username,
		DisplayName: art.Account.DisplayName,
		AvatarURL:   art.Account.AvatarURL,
	}

	return &models.RenderTree{
		ID:             art.ID,
		Platform:       "twitter",
		Author:         author,
		OriginalPost:   models.RenderPost{FullText: art.FullText, CreatedAt: art.CreatedAt.Format("2006-01-02T15:04:05Z")},
		DecoratedDOM:   models.RenderDOM{Japanese: art.FullTextJA, English: art.FullTextEN, Chinese: art.FullTextZH},
		Media:          mediaItems,
		IsLiked:        art.IsLiked,
		IsRepost:       art.IsRepost,
		ConversationID: art.ConversationID,
	}
}

func calculateAspectRatio(w, h int) string {
	if w <= 0 || h <= 0 {
		return "16/9"
	}
	if w == h {
		return "1/1"
	}
	return "16/9"
}
