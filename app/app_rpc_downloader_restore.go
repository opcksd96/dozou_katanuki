// app/app_rpc_downloader_restore.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"

	"dozou_katanuki/models"
)

// FetchDownloadReserves は退避予約一覧を取得します
func (a *App) FetchDownloadReserves(status string) ([]models.DownloadReserve, error) {
	if a.Repo == nil {
		return nil, nil
	}
	list, _, err := a.Repo.FetchDownloadReserves(status, 100, 0)
	return list, err
}

// RestoreReserveToMotrix は退避レコードを Motrix へ再投入します
func (a *App) RestoreReserveToMotrix(id uint) (bool, error) {
	if a.Repo == nil {
		return false, fmt.Errorf("repo unavailable")
	}
	reserves, _, err := a.Repo.FetchDownloadReserves("all", 1000, 0)
	if err != nil {
		return false, err
	}
	for _, r := range reserves {
		if r.ID == id {
			urls := []string{r.URL}
			if r.MirrorURLs != "" {
				var mirrors []string
				if json.Unmarshal([]byte(r.MirrorURLs), &mirrors) == nil {
					urls = append(urls, mirrors...)
				}
			}
			gid, err := a.AddMotrixDownload(urls, "", r.FileName)
			if err == nil && gid != "" {
				_ = a.Repo.UpdateReserveStatus(id, "retrying", "Restored to Motrix: "+gid)
				return true, nil
			}
			return false, err
		}
	}
	return false, fmt.Errorf("reserve record not found")
}
