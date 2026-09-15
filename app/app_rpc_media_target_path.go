// app/app_rpc_media_target_path.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveMediaTargetDir は メディアIDから正規のアカウント別保存フォルダを解決・事前作成します
func (a *App) ResolveMediaTargetDir(mediaID string) string {
	destRoot := a.getMediaDownloadDir()
	owner := "unknown"
	if a.Repo != nil && mediaID != "" {
		if o, err := a.Repo.GetMediaOwnerUsername(mediaID); err == nil && o != "" {
			owner = o
		}
	}
	targetDir := filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
	_ = os.MkdirAll(targetDir, 0755)
	return strings.ReplaceAll(filepath.Clean(targetDir), "/", `\`)
}

// CheckMediaFileExists は 正規のアカウント別フォルダまたはルートに実体が存在するかを確認します
func (a *App) CheckMediaFileExists(mediaID, fileName string) bool {
	if fileName == "" { return false }
	targetDir := a.ResolveMediaTargetDir(mediaID)
	if fi, err := os.Stat(filepath.Join(targetDir, fileName)); err == nil && fi.Size() > 0 {
		return true
	}
	if fi, err := os.Stat(filepath.Join(a.getMediaDownloadDir(), fileName)); err == nil && fi.Size() > 0 {
		return true
	}
	return false
}
