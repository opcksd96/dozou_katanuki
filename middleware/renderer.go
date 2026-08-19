// middleware/renderer.go (100行以下)
package middleware

import (
	"fmt"
	"dozou_katanuki/models"
)

// ToRenderTree は装飾済みの中間データをフロントエンド描画用構造体へ即時マッピングします
func ToRenderTree(item models.Article, platform string) models.RenderTree {
	avatarURL := AuditAndResolveAvatar(platform, item.CreatedAt, item.Account.ProfileHistory)

	return models.RenderTree{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Content: models.RenderContent{
			Original: item.FullText,
			JA:       item.FullTextJA.String,
			EN:       item.FullTextEN.String,
			ZH:       item.FullTextZH.String,
		},
		Author: models.RenderAuthor{
			NumericID:   item.Account.NumericID,
			Handle:      item.Account.Username,
			DisplayName: item.Account.DisplayName,
			AvatarURL:   avatarURL,
		},
		Media: mapMediaToRenderMedia(item.Media),
		Metrics: models.RenderMetrics{
			Replies: 0,
			Reposts: 0,
			Likes:   0,
		},
		IsLiked:   item.IsLiked,
		SourceURL: item.WaybackURL,
	}
}

func mapMediaToRenderMedia(mediaList []models.Media) []models.RenderMedia {
	result := make([]models.RenderMedia, 0, len(mediaList))
	for _, m := range mediaList {
		var urls []string
		if m.DownloadStatus == "COMPLETED" {
			if m.Type == "video" || m.Type == "gif" {
				urls = append(urls, fmt.Sprintf("/stash-proxy/scene/%s/stream", m.StashSceneID.String))
			} else {
				urls = append(urls, fmt.Sprintf("/stash-proxy/image/%s/file", m.StashImageID.String))
			}
		}
		result = append(result, models.RenderMedia{
			Type:         m.Type,
			URLs:         urls,
			FailedReason: m.FailedReason.String,
		})
	}
	return result
}
