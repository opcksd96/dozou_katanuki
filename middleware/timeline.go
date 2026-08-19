// middleware/timeline.go (100行以下)
package middleware

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

// FetchTimeline はパラメータ検証を経て RenderTree 配列を一方向データフローで供給します
func (s *TimelineService) FetchTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if platform == "" || accountID == "" {
		return nil, fmt.Errorf("invalid parameters: platform and account_id are required")
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	articles, err := FetchInChunks(s.repo, accountID, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	trees := make([]models.RenderTree, 0, len(articles))
	for _, article := range articles {
		trees = append(trees, ToRenderTree(article, platform))
	}
	return trees, nil
}

func (s *TimelineService) ResolveAvatar(platform, accountID string, tweetAt time.Time) (string, error) {
	histories, err := s.repo.GetAccountHistories(accountID)
	if err != nil {
		return "", err
	}
	return AuditAndResolveAvatar(platform, tweetAt, histories), nil
}
