package services

import (
	"fmt"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

type TimelineService struct {
	repo *driver.Repository
}

func NewTimelineService(repo *driver.Repository) *TimelineService {
	return &TimelineService{repo: repo}
}

// FetchTimeline はDBレコードを取得し、RenderTree配列へ高速変換して返します
func (s *TimelineService) FetchTimeline(platform, accountID, filter string, limit, offset int) ([]RenderTree, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	if filter == "" {
		filter = "all"
	}

	articles, err := s.repo.GetArticles(accountID, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("timeline fetch error: %w", err)
	}

	trees := make([]RenderTree, len(articles))
	for i, art := range articles {
		trees[i] = s.toRenderTree(art, platform)
	}
	return trees, nil
}

func (s *TimelineService) toRenderTree(art models.Article, platform string) RenderTree {
	avatarURL := fmt.Sprintf("/assets/%s/%s.jpg", platform, art.Account.AvatarURL)
	if art.Account.AvatarURL == "" {
		avatarURL = "/assets/default_avatar.png"
	}

	mediaList := make([]RenderMediaItem, len(art.Media))
	for i, m := range art.Media {
		mediaList[i] = RenderMediaItem{
			MediaID: m.MediaID,
			Type:    m.Type,
			Width:   m.Width,
			Height:  m.Height,
			URLs: RenderMediaURLs{
				Stream: fmt.Sprintf("/stash-proxy/scene/%s", m.StashSceneID.String),
				Image:  fmt.Sprintf("/stash-proxy/image/%s", m.StashImageID.String),
			},
		}
		if m.FailedReason.Valid {
			mediaList[i].FailedReason = m.FailedReason.String
		}
	}

	return RenderTree{
		ID:             art.ID,
		ConversationID: art.ConversationID,
		CreatedAt:      art.CreatedAt.Format(time.RFC3339),
		Content: RenderContent{
			Original: art.FullText,
			JA:       art.FullTextJA.String,
			EN:       art.FullTextEN.String,
			ZH:       art.FullTextZH.String,
		},
		Author: RenderAuthor{
			NumericID:   art.Account.NumericID,
			Handle:      art.Account.Username,
			DisplayName: art.Account.DisplayName,
			AvatarURL:   avatarURL,
		},
		Media:     mediaList,
		IsLiked:   art.IsLiked,
		SourceURL: art.WaybackURL,
	}
}
