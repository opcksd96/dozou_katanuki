// middleware/avatar_resolver.go (100行以下)
package middleware

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

// AuditAndResolveAvatar は投稿日時に基づくアバター (Base64またはURL) を解決します
func AuditAndResolveAvatar(platform string, tweetAt time.Time, histories []models.AccountProfileHistory) string {
	if len(histories) > 0 {
		for _, h := range histories {
			if h.ObservedAt.Before(tweetAt) || h.ObservedAt.Equal(tweetAt) {
				if h.AvatarBase64 != "" { return h.AvatarBase64 }
				if h.AvatarVirtualKey != "" { return fmt.Sprintf("/avatars/%s/%s.jpg", platform, h.AvatarVirtualKey) }
			}
		}
		latest := histories[len(histories)-1]
		if latest.AvatarBase64 != "" { return latest.AvatarBase64 }
		if latest.AvatarVirtualKey != "" { return fmt.Sprintf("/avatars/%s/%s.jpg", platform, latest.AvatarVirtualKey) }
	}
	return fmt.Sprintf("/avatars/%s/default_avatar.jpg", platform)
}

// ResolveAccountAvatar はアカウント情報全体からアバターを解決します (世代履歴 -> Account.AvatarBase64 -> Account.AvatarURL -> デフォルト)
func ResolveAccountAvatar(platform string, tweetAt time.Time, acc models.Account) string {
	if res := AuditAndResolveAvatar(platform, tweetAt, acc.ProfileHistory); res != fmt.Sprintf("/avatars/%s/default_avatar.jpg", platform) {
		return res
	}
	if acc.AvatarBase64 != "" { return acc.AvatarBase64 }
	if acc.AvatarURL != "" { return acc.AvatarURL }
	return fmt.Sprintf("/avatars/%s/default_avatar.jpg", platform)
}
