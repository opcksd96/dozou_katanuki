// app_rpc_database.go (100行以下)
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// GetAccountDetail は指定されたアカウントの詳細とアバター世代履歴を取得する Wails バインドメソッドです
func (a *App) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	res, err := a.repo.GetAccountDetail(numericID)
	if err != nil {
		log.Printf("[Wails RPC] GetAccountDetail error: %v", err)
		return nil, err
	}
	return res, nil
}

// GetMediaList はメディア一覧（Stashステータス・アカウント情報・種別フィルタ付き）を取得する Wails バインドメソッドです
func (a *App) GetMediaList(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	_, _ = a.repo.MigrateExcludedMedia() // レガシーWhitelist外DEAD_404をEXCLUDEDへ安全移行
	res, err := a.repo.SearchMediaDetails(accountID, status, mediaType, limit, offset)
	if err != nil {
		log.Printf("[Wails RPC] GetMediaList error: %v", err)
		return nil, err
	}
	return res, nil
}

// PurgeMedia は指定された単一メディアレコードをデータベースから物理削除する Wails バインドメソッドです
func (a *App) PurgeMedia(mediaID string) error {
	if err := a.waitForReady(); err != nil { return err }
	return a.repo.PurgeMedia(mediaID)
}

// PurgeMediaByStatus は指定ステータス（EXCLUDED, UNLINKED, DEAD_404等）のメディアを一括削除する Wails バインドメソッドです
func (a *App) PurgeMediaByStatus(status, accountID string) (int64, error) {
	if err := a.waitForReady(); err != nil { return 0, err }
	return a.repo.PurgeMediaByStatus(status, accountID)
}

// GetTableRecords は汎用テーブルの生カラム・ロウ形式スプレッドシートデータを取得する Wails バインドメソッドです
func (a *App) GetTableRecords(tableName string, limit, offset int, search string) (*models.TableRecordResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	res, err := a.repo.GetTableRecords(tableName, limit, offset, search)
	if err != nil {
		log.Printf("[Wails RPC] GetTableRecords error: %v", err)
		return nil, err
	}
	return res, nil
}

// ListAllAccounts は登録済みアカウント一覧を取得する Wails バインドメソッドです
func (a *App) ListAllAccounts() ([]models.Account, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	return a.repo.ListAccounts()
}

// UpdateAccount はアカウント情報（表示名・ユーザー名・アバターURL・一言コメント）を更新する Wails バインドメソッドです
func (a *App) UpdateAccount(numericID, displayName, username, avatarURL, description string) error {
	if err := a.waitForReady(); err != nil { return err }

	// セキュリティ: パストラバーサル・不正入力の遮断
	if strings.Contains(avatarURL, "..") || strings.Contains(avatarURL, "\\") || strings.Contains(avatarURL, "\x00") {
		return fmt.Errorf("invalid avatar_url: path traversal characters are forbidden")
	}

	err := a.repo.UpdateAccount(numericID, displayName, username, avatarURL, description)
	if err != nil {
		log.Printf("[Wails RPC] UpdateAccount error: %v", err)
		return err
	}
	return nil
}

// SaveAvatarImage は Base64 形式の画像データをデコードし assets/{platform}/{virtualKey}.jpg として保存します
func (a *App) SaveAvatarImage(platform, virtualKey, base64Data string) (string, error) {
	if err := a.waitForReady(); err != nil { return "", err }
	if platform == "" { platform = "twitter" }
	if virtualKey == "" { return "", fmt.Errorf("virtualKey is required") }

	// セキュリティ: virtualKey のサニタイズ（英数字・アンダースコア・ハイフンのみ許可）
	for _, c := range virtualKey {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.') {
			return "", fmt.Errorf("invalid characters in virtualKey")
		}
	}
	if strings.Contains(virtualKey, "..") {
		return "", fmt.Errorf("path traversal is forbidden")
	}

	// 拡張子が付いていない場合は .jpg を付与
	filename := virtualKey
	if !strings.HasSuffix(strings.ToLower(filename), ".jpg") &&
		!strings.HasSuffix(strings.ToLower(filename), ".jpeg") &&
		!strings.HasSuffix(strings.ToLower(filename), ".png") &&
		!strings.HasSuffix(strings.ToLower(filename), ".webp") {
		filename += ".jpg"
	}

	// Base64 プレフィックス (data:image/...;base64,) の除去
	idx := strings.Index(base64Data, ";base64,")
	rawB64 := base64Data
	if idx != -1 {
		rawB64 = base64Data[idx+8:]
	}

	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	targetDir := filepath.Join("./assets", platform)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create asset dir: %w", err)
	}

	targetPath := filepath.Join(targetDir, filename)
	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write avatar file: %w", err)
	}

	// DBに直接 Base64 イメージを格納（Same Source, Same Flow）
	cleanKey := strings.TrimSuffix(virtualKey, filepath.Ext(virtualKey))
	normalizedB64 := base64Data
	if !strings.HasPrefix(normalizedB64, "data:") {
		normalizedB64 = fmt.Sprintf("data:image/jpeg;base64,%s", rawB64)
	}
	_ = a.repo.UpdateAvatarBase64ByVirtualKey(cleanKey, normalizedB64)

	log.Printf("[Wails RPC] Saved avatar image to disk & DB: %s (%d bytes)", targetPath, len(data))
	return fmt.Sprintf("/avatars/%s/%s", platform, filename), nil
}

// ListAvailableAvatars は assets/{platform} 内に存在する利用可能なアバター画像一覧を取得します
func (a *App) ListAvailableAvatars(platform string) ([]string, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	if platform == "" { platform = "twitter" }

	targetDir := filepath.Join("./assets", platform)
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		if os.IsNotExist(err) { return []string{}, nil }
		return nil, err
	}

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".svg": true, ".gif": true}
	var results []string
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if allowedExts[ext] {
				results = append(results, fmt.Sprintf("/avatars/%s/%s", platform, entry.Name()))
			}
		}
	}
	return results, nil
}

