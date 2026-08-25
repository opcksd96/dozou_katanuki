// middleware/avatar_resolver.go (100行以下)
package middleware

import (
	"time"

	"dozou_katanuki/models"
)

// DefaultHumanAvatarSVG はアバター画像未取得時に安全に描画される標準の人型シルエットSVG (Data URI) です
const DefaultHumanAvatarSVG = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>"

// AuditAndResolveAvatar は投稿日時に合致する世代のアバター (Base64) を解決します
func AuditAndResolveAvatar(platform string, tweetAt time.Time, histories []models.AccountProfileHistory) string {
	if len(histories) > 0 {
		for _, h := range histories {
			if h.ObservedAt.Before(tweetAt) || h.ObservedAt.Equal(tweetAt) {
				if h.AvatarBase64 != "" { return h.AvatarBase64 }
			}
		}
		latest := histories[len(histories)-1]
		if latest.AvatarBase64 != "" { return latest.AvatarBase64 }
	}
	return ""
}

// ResolveAccountAvatar はアカウント情報からアバターを解決します (Account.AvatarBase64 -> 世代履歴 -> デフォルト人型SVG)
func ResolveAccountAvatar(platform string, tweetAt time.Time, acc models.Account) string {
	if acc.AvatarBase64 != "" {
		return acc.AvatarBase64
	}
	if res := AuditAndResolveAvatar(platform, tweetAt, acc.ProfileHistory); res != "" {
		return res
	}
	return DefaultHumanAvatarSVG
}


