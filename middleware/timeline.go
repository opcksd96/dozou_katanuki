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

// GetAccounts は登録されている全アカウントの RenderAuthor リストを返却します
func (s *TimelineService) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	accounts, err := s.repo.GetAccounts()
	if err != nil {
		return nil, err
	}
	authors := make([]models.RenderAuthor, 0, len(accounts))
	for _, acc := range accounts {
		avatarURL := AuditAndResolveAvatar(platform, time.Now(), acc.ProfileHistory)
		authors = append(authors, models.RenderAuthor{
			NumericID:   acc.NumericID,
			Handle:      acc.Username,
			DisplayName: acc.DisplayName,
			AvatarURL:   avatarURL,
			Bio:         acc.Description,
		})
	}
	return authors, nil
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

// GetArticleDetail は指定された記事および同一会話スレッド全体のツリーを取得します
func (s *TimelineService) GetArticleDetail(platform, id string) (*models.ArticleDetailResult, error) {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		return nil, err
	}
	targetTree := ToRenderTree(*article, platform)

	var threadTrees []models.RenderTree
	if article.ConversationID != "" {
		threadArticles, err := s.repo.GetArticlesByConversationID(article.ConversationID)
		if err == nil {
			threadTrees = make([]models.RenderTree, 0, len(threadArticles))
			for _, ta := range threadArticles {
				threadTrees = append(threadTrees, ToRenderTree(ta, platform))
			}
		}
	}
	return &models.ArticleDetailResult{
		Article: targetTree,
		Thread:  threadTrees,
	}, nil
}

// SearchArticles は検索クエリに基づいて記事を検索し RenderTree 形式で返却します
func (s *TimelineService) SearchArticles(query, accountID, filter string, limit, offset int) (*models.ArticleSearchResult, error) {
	articles, total, err := s.repo.SearchArticles(query, accountID, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	trees := make([]models.RenderTree, 0, len(articles))
	for _, article := range articles {
		trees = append(trees, ToRenderTree(article, "twitter"))
	}
	return &models.ArticleSearchResult{
		Items: trees,
		Total: total,
	}, nil
}

// SearchMediaDetails は RawMediaItems を取得し AuditAndResolveAvatar で完全解決した MediaSearchResult を返却します (Same Source, Same Flow)
func (s *TimelineService) SearchMediaDetails(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	rawItems, total, stats, err := s.repo.FetchRawMediaItems(accountID, status, mediaType, limit, offset)
	if err != nil {
		return nil, err
	}

	details := make([]models.MediaItemDetail, 0, len(rawItems))
	for _, item := range rawItems {
		rm := models.BuildRenderMedia(item.Media)
		hasStash := item.Media.StashSceneID.Valid && item.Media.StashSceneID.String != "" || item.Media.StashImageID.Valid && item.Media.StashImageID.String != ""
		avatarURL := AuditAndResolveAvatar("twitter", item.CreatedAt, item.ProfileHistory)
		title := fmt.Sprintf("X (@%s): Tweet %s", item.Username, item.ArticleID)
		tweetDate := item.CreatedAt.Format("2006-01-02")

		details = append(details, models.MediaItemDetail{
			RenderMedia: rm,
			MediaID:     item.Media.MediaID,
			ArticleID:   item.ArticleID,
			AccountID:   item.AccountID,
			Username:    item.Username,
			DisplayName: item.DisplayName,
			AvatarURL:   avatarURL,
			RawStatus:   item.Media.DownloadStatus,
			HasStash:    hasStash,
			CreatedAt:   item.CreatedAt,
			Title:       title,
			FullText:    item.FullText,
			FullTextJA:  item.FullTextJA,
			TweetDate:   tweetDate,
		})
	}

	return &models.MediaSearchResult{
		Items: details,
		Total: total,
		Stats: stats,
	}, nil
}

