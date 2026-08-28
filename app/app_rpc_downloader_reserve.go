package app

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"time"

	"dozou_katanuki/models"
)

// SafePurgeWithBackup はタスクメタデータをDBに安全退避してからAria2キューから削除します
func (a *App) SafePurgeWithBackup(gids []string) (int, error) {
	if len(gids) == 0 {
		return 0, nil
	}

	purgedCount := 0
	for _, gid := range gids {
		if gid == "" {
			continue
		}

		// 1. tellStatus でタスクの詳細メタデータを取得
		rawStatus, err := callMotrixRPC("aria2.tellStatus", []interface{}{gid, []string{"gid", "status", "totalLength", "completedLength", "files", "errorMessage"}})
		if err == nil {
			var res struct {
				Result struct {
					GID, Status, TotalLength, CompletedLength, ErrorMessage string
					Files []struct {
						Path string `json:"path"`
						URIs []struct {
							URI string `json:"uri"`
						} `json:"uris"`
					} `json:"files"`
				} `json:"result"`
			}
			if json.Unmarshal(rawStatus, &res) == nil && res.Result.GID != "" {
				fileName, uri := "", ""
				if len(res.Result.Files) > 0 {
					fileName = filepath.Base(res.Result.Files[0].Path)
					if len(res.Result.Files[0].URIs) > 0 {
						uri = res.Result.Files[0].URIs[0].URI
					}
				}
				tl, _ := strconv.ParseInt(res.Result.TotalLength, 10, 64)

				// 派生フォールバックURL群を生成
				fallbacks := BuildThunderFallbackURLs(uri)
				mirrorJSON, _ := json.Marshal(fallbacks)

				reserve := &models.DownloadReserve{
					GID:         gid,
					URL:         uri,
					FileName:    fileName,
					MirrorURLs:  string(mirrorJSON),
					Status:      "reserved",
					Reason:      "Safe purged from Motrix/Aria2 active queue",
					TotalLength: tl,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				// DBリポジトリへ退避保存
				if a.Repo != nil {
					_ = a.Repo.SaveDownloadReserve(reserve)
				}
			}
		}

		// 2. Aria2 キューから安全に除去 (remove または removeDownloadResult)
		_, _ = callMotrixRPC("aria2.remove", []interface{}{gid})
		_, _ = callMotrixRPC("aria2.removeDownloadResult", []interface{}{gid})
		purgedCount++
	}

	return purgedCount, nil
}
