// driver/migrate_avatars.go (100行以下)
package driver

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

// MigrateAvatarsToBase64 は assets/ フォルダのアバター画像を読み込み、DB の avatar_base64 に格納します
func MigrateAvatarsToBase64(db *gorm.DB) error {
	baseDirs := []string{
		"./assets/twitter",
		"./assets",
		"../assets/twitter",
		"../../assets/twitter",
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		baseDirs = append(baseDirs, filepath.Join(exeDir, "assets", "twitter"), filepath.Join(exeDir, "..", "assets", "twitter"))
	}

	// 1. account_profile_histories の Base64 格納
	var histories []models.AccountProfileHistory
	if err := db.Where("avatar_base64 IS NULL OR avatar_base64 = ''").Find(&histories).Error; err != nil {
		return err
	}

	migratedCount := 0
	for _, h := range histories {
		if h.AvatarVirtualKey == "" { continue }
		cleanKey := strings.TrimSpace(h.AvatarVirtualKey)

		b64Data := findAndEncodeImage(baseDirs, cleanKey)
		if b64Data != "" {
			_ = db.Model(&models.AccountProfileHistory{}).Where("id = ?", h.ID).Update("avatar_base64", b64Data).Error
			migratedCount++
		}
	}

	// 2. accounts の avatar_base64 を最新世代またはアバター画像から同期
	var accounts []models.Account
	if err := db.Where("avatar_base64 IS NULL OR avatar_base64 = ''").Find(&accounts).Error; err == nil {
		for _, acc := range accounts {
			var latestHist models.AccountProfileHistory
			if err := db.Where("account_id = ? AND avatar_base64 != ''", acc.NumericID).Order("avatar_seq DESC").First(&latestHist).Error; err == nil {
				_ = db.Model(&models.Account{}).Where("numeric_id = ?", acc.NumericID).Update("avatar_base64", latestHist.AvatarBase64).Error
			} else if acc.AvatarURL != "" {
				key := filepath.Base(acc.AvatarURL)
				key = strings.TrimSuffix(key, filepath.Ext(key))
				b64Data := findAndEncodeImage(baseDirs, key)
				if b64Data != "" {
					_ = db.Model(&models.Account{}).Where("numeric_id = ?", acc.NumericID).Update("avatar_base64", b64Data).Error
				}
			}
		}
	}

	if migratedCount > 0 {
		log.Printf("[Migration] Successfully migrated %d avatar images into DB as Base64", migratedCount)
	}
	return nil
}

func findAndEncodeImage(baseDirs []string, key string) string {
	exts := []string{"", ".jpg", ".jpeg", ".png", ".webp"}
	for _, dir := range baseDirs {
		for _, ext := range exts {
			fp := filepath.Join(dir, key+ext)
			if info, err := os.Stat(fp); err == nil && !info.IsDir() {
				bytes, err := os.ReadFile(fp)
				if err == nil && len(bytes) > 0 {
					mime := "image/jpeg"
					lfp := strings.ToLower(fp)
					if strings.HasSuffix(lfp, ".png") { mime = "image/png" }
					if strings.HasSuffix(lfp, ".webp") { mime = "image/webp" }
					return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(bytes))
				}
			}
		}
	}
	return ""
}
