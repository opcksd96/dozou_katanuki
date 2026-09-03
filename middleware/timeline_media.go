// middleware/timeline_media.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"
	"time"

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
		accMock := models.Account{AvatarBase64: item.AvatarBase64, ProfileHistory: item.ProfileHistory}
		avatarURL := ResolveAccountAvatar("twitter", item.CreatedAt, accMock)
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
	accs, err := s.repo.ListAccounts()
	if err != nil {
		return nil, err
	}
	for i := range accs {
		accs[i].AvatarURL = ResolveAccountAvatar("twitter", time.Now(), accs[i])
	}
	return accs, nil
}

// GetAccountDetail は指定アカウントの詳細・世代変遷履歴・投稿件数を取得します
func (s *TimelineService) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	detail, err := s.repo.GetAccountDetail(numericID)
	if err != nil {
		return nil, err
	}
	if detail != nil {
		detail.Account.AvatarURL = ResolveAccountAvatar("twitter", time.Now(), detail.Account)
	}
	return detail, nil
}
