package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

// AuditAndResolveAvatar はアバター世代を監査し、必要に応じて履歴を追加して仮想キーを返します
func AuditAndResolveAvatar(tx *gorm.DB, acc *models.Account) (string, error) {
	var histories []models.AccountProfileHistory
	err := tx.Where("account_id = ?", acc.NumericID).Order("avatar_seq DESC").Limit(1).Find(&histories).Error
	if err != nil {
		return "", fmt.Errorf("failed to query avatar history: %w", err)
	}

	seq := 1
	if len(histories) > 0 {
		latest := histories[0]
		if latest.AvatarOriginalURL == acc.AvatarURL {
			return latest.AvatarVirtualKey, nil
		}
		seq = latest.AvatarSeq + 1
	}

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
