// middleware/timeline_media_actions.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"dozou_katanuki/models"
)

// ToggleMediaBookmark はメディアのブックマーク状態を反転します
func (s *TimelineService) ToggleMediaBookmark(mediaID string) (bool, error) {
	return s.repo.ToggleMediaBookmark(mediaID)
}

// UpdateMediaMetadata はメディアのメタデータを更新します
func (s *TimelineService) UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason string) error {
	return s.repo.UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason)
}

// PurgeMedia は指定された単一メディアを物理削除します
func (s *TimelineService) PurgeMedia(mediaID string) error {
	return s.repo.PurgeMedia(mediaID)
}

// PurgeMediaByStatus は指定ステータスのメディアを一括削除します
func (s *TimelineService) PurgeMediaByStatus(status, accountID string) (int64, error) {
	return s.repo.PurgeMediaByStatus(status, accountID)
}

// RequeueMediaByStatus は指定ステータスのメディアを一括で QUEUED 状態に戻します
func (s *TimelineService) RequeueMediaByStatus(status, accountID string) (int64, error) {
	return s.repo.RequeueMediaByStatus(status, accountID)
}

// ResolveMediaFilePath はメディアのローカル実体ファイルパスを解決します
func (s *TimelineService) ResolveMediaFilePath(mediaID string) (string, error) {
	return s.repo.ResolveMediaFilePath(mediaID)
}

// GetMediaByID はメディア情報を1件取得します
func (s *TimelineService) GetMediaByID(mediaID string) (*models.Media, error) {
	return s.repo.GetMediaByID(mediaID)
}
