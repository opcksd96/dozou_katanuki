// middleware/avatar_resolver.go (100行以下)
package middleware

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

// AuditAndResolveAvatar は投稿日時に合致する世代のアバター (Base64または原本URL) を解決します
func AuditAndResolveAvatar(platform string, tweetAt time.Time, histories []models.AccountProfileHistory) string {
	if len(histories) > 0 {
		for _, h := range histories {
			if h.ObservedAt.Before(tweetAt) || h.ObservedAt.Equal(tweetAt) {
				if h.AvatarBase64 != "" { return h.AvatarBase64 }
				if h.AvatarOriginalURL != "" { return h.AvatarOriginalURL }
			}
		}
		latest := histories[len(histories)-1]
		if latest.AvatarBase64 != "" { return latest.AvatarBase64 }
		if latest.AvatarOriginalURL != "" { return latest.AvatarOriginalURL }
	}
	return ""
}

// ResolveAccountAvatar はアカウント情報からアバターを解決します (Account.AvatarBase64 -> 世代履歴 -> Account.AvatarURL -> デフォルト)
func ResolveAccountAvatar(platform string, tweetAt time.Time, acc models.Account) string {
	if acc.AvatarBase64 != "" {
		return acc.AvatarBase64
	}
	if res := AuditAndResolveAvatar(platform, tweetAt, acc.ProfileHistory); res != "" {
		return res
	}
	if acc.AvatarURL != "" {
		return acc.AvatarURL
	}
	return fmt.Sprintf("/avatars/%s/default_avatar.jpg", platform)
}


