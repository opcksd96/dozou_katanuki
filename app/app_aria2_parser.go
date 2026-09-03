// app/app_aria2_parser.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"path/filepath"
	"strconv"

	"dozou_katanuki/models"
)

func parseAria2Tasks(raw []byte) []models.DownloaderTaskInfo {
	var res struct {
		Result []struct {
			GID, Status, TotalLength, CompletedLength, DownloadSpeed, ErrorMessage string
			Files                                                                  []struct {
				Path string `json:"path"`
				URIs []struct {
					URI string `json:"uri"`
				} `json:"uris"`
			} `json:"files"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil
	}
	var tasks []models.DownloaderTaskInfo
	for _, t := range res.Result {
		tl, _ := strconv.ParseInt(t.TotalLength, 10, 64)
		cl, _ := strconv.ParseInt(t.CompletedLength, 10, 64)
		ds, _ := strconv.ParseInt(t.DownloadSpeed, 10, 64)
		fn, uri := "", ""
		if len(t.Files) > 0 {
			fn = filepath.Base(t.Files[0].Path)
			if len(t.Files[0].URIs) > 0 {
				uri = t.Files[0].URIs[0].URI
			}
		}
		prog := 0.0
		if tl > 0 {
			prog = float64(cl) / float64(tl) * 100
		}
		tasks = append(tasks, models.DownloaderTaskInfo{
			GID:             t.GID,
			Status:          t.Status,
			FileName:        fn,
			URL:             uri,
			TotalLength:     tl,
			CompletedLength: cl,
			DownloadSpeed:   ds,
			Progress:        prog,
			ErrorMessage:    t.ErrorMessage,
		})
	}
	return tasks
}
