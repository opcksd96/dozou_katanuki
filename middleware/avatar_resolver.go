// middleware/avatar_resolver.go (100行以下)
package middleware

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

// AuditAndResolveAvatar は投稿日時に基づくアバター専用パス (/avatars/...) を解決します
func AuditAndResolveAvatar(platform string, tweetAt time.Time, histories []models.AccountProfileHistory) string {
	if len(histories) == 0 {
		return fmt.Sprintf("/avatars/%s/default_avatar.jpg", platform)
	}

	for _, h := range histories {
		if h.ObservedAt.Before(tweetAt) || h.ObservedAt.Equal(tweetAt) {
			return fmt.Sprintf("/avatars/%s/%s.jpg", platform, h.AvatarVirtualKey)
		}
	}

	return fmt.Sprintf("/avatars/%s/%s.jpg", platform, histories[len(histories)-1].AvatarVirtualKey)
}
