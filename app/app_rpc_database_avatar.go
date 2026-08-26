// app/app_rpc_database_avatar.go (100行以下)
package app

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// UpdateAccount はアカウント情報（表示名・ユーザー名・アバターURL・一言コメント・名寄せ先・グループ名）を更新する Wails バインドメソッドです
func (a *App) UpdateAccount(numericID, displayName, username, avatarURL, description, aliasOf, groupName string) error {
	if err := a.WaitForReady(); err != nil { return err }
	if strings.Contains(avatarURL, "..") || strings.Contains(avatarURL, "\\") || strings.Contains(avatarURL, "\x00") {
		return fmt.Errorf("invalid avatar_url: path traversal characters are forbidden")
	}
	err := a.Repo.UpdateAccount(numericID, displayName, username, avatarURL, description, aliasOf, groupName)
	if err != nil {
		log.Printf("[Wails RPC] UpdateAccount error: %v", err)
		return err
	}
	return nil
}

// SaveAvatarImage は Base64 形式の画像データをデコードし assets/{platform}/{virtualKey}.jpg として保存します
func (a *App) SaveAvatarImage(platform, virtualKey, base64Data string) (string, error) {
	if err := a.WaitForReady(); err != nil { return "", err }
	if platform == "" { platform = "twitter" }
	if virtualKey == "" { return "", fmt.Errorf("virtualKey is required") }

	for _, c := range virtualKey {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.') {
			return "", fmt.Errorf("invalid characters in virtualKey")
		}
	}
	if strings.Contains(virtualKey, "..") { return "", fmt.Errorf("path traversal is forbidden") }

	filename := virtualKey
	if !strings.HasSuffix(strings.ToLower(filename), ".jpg") && !strings.HasSuffix(strings.ToLower(filename), ".png") && !strings.HasSuffix(strings.ToLower(filename), ".webp") {
		filename += ".jpg"
	}

	idx := strings.Index(base64Data, ";base64,")
	rawB64 := base64Data
	if idx != -1 { rawB64 = base64Data[idx+8:] }
	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil { return "", fmt.Errorf("failed to decode base64: %w", err) }

	targetDir := filepath.Join("./assets", platform)
	if err := os.MkdirAll(targetDir, 0755); err != nil { return "", fmt.Errorf("failed to create asset dir: %w", err) }
	targetPath := filepath.Join(targetDir, filename)
	if err := os.WriteFile(targetPath, data, 0644); err != nil { return "", fmt.Errorf("failed to write avatar file: %w", err) }

	cleanKey := strings.TrimSuffix(virtualKey, filepath.Ext(virtualKey))
	normalizedB64 := base64Data
	if !strings.HasPrefix(normalizedB64, "data:") { normalizedB64 = fmt.Sprintf("data:image/jpeg;base64,%s", rawB64) }
	_ = a.Repo.UpdateAvatarBase64ByVirtualKey(cleanKey, normalizedB64)

	log.Printf("[Wails RPC] Saved avatar image: %s (%d bytes)", targetPath, len(data))
	return fmt.Sprintf("/avatars/%s/%s", platform, filename), nil
}

// ListAvailableAvatars は assets/{platform} 内に存在する利用可能なアバター画像一覧を取得します
func (a *App) ListAvailableAvatars(platform string) ([]string, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
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
		if !entry.IsDir() && allowedExts[strings.ToLower(filepath.Ext(entry.Name()))] {
			results = append(results, fmt.Sprintf("/avatars/%s/%s", platform, entry.Name()))
		}
	}
	return results, nil
}
