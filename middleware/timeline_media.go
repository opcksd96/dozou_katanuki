// middleware/timeline_media.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"

	"dozou_katanuki/models"
)

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
			WaybackURL:  item.WaybackURL,
		})
	}

	return &models.MediaSearchResult{
		Items: details,
		Total: total,
		Stats: stats,
	}, nil
}

// FetchDownloadStatusStats はダウンロードキュー全体のステータス別集計を取得します
func (s *TimelineService) FetchDownloadStatusStats(accountID string) (*models.DownloadStatusStats, error) {
	return s.repo.FetchDownloadStatusStats(accountID)
}

// ListRawAccounts は登録済みアカウントRawエンティティ（usernameフィールド付き）一覧を取得します
func (s *TimelineService) ListRawAccounts() ([]models.Account, error) {
	return s.repo.ListAccounts()
}

// GetAccountDetail は指定アカウントの詳細・世代変遷履歴・投稿件数を取得します
func (s *TimelineService) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	return s.repo.GetAccountDetail(numericID)
}
