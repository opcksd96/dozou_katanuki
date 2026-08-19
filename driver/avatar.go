package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

// AuditAndResolveAvatar はアバター世代を監査し、必要に応じて履歴を追加して仮想キーを返します
func AuditAndResolveAvatar(tx *gorm.DB, acc *models.Account) (string, error) {
	var latest models.AccountProfileHistory
	err := tx.Where("account_id = ?", acc.NumericID).Order("avatar_seq DESC").First(&latest).Error

	seq := 1
	if err == nil {
		// URLに変更がない場合は既存の仮想キーをそのまま維持
		if latest.AvatarOriginalURL == acc.AvatarURL {
			return latest.AvatarVirtualKey, nil
		}
		seq = latest.AvatarSeq + 1
	}

	// 3桁世代キーを生成 (例: msluo14_avatar_001)
	virtualKey := fmt.Sprintf("%s_avatar_%03d", acc.Username, seq)
	history := models.AccountProfileHistory{
		AccountID:         acc.NumericID,
		DisplayName:       acc.DisplayName,
		AvatarOriginalURL: acc.AvatarURL,
		AvatarSeq:         seq,
		AvatarVirtualKey:  virtualKey,
		ObservedAt:        time.Now(),
	}

	if err := tx.Create(&history).Error; err != nil {
		return "", fmt.Errorf("failed to insert avatar history: %w", err)
	}

	return virtualKey, nil
}
