// driver/repo_media_intelligence.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// MergeDuplicateMedia は同一ファイル名や重複URLを持つメディアレコードを統合します
func (r *Repository) MergeDuplicateMedia() (int, error) {
	var duplicates []struct {
		DownloadURL string
		Count       int
	}
	err := r.db.Model(&models.Media{}).
		Select("download_url, count(*) as count").
		Where("download_url != ''").
		Group("download_url").
		Having("count(*) > 1").
		Scan(&duplicates).Error
	if err != nil {
		return 0, err
	}

	mergedCount := 0
	for _, dup := range duplicates {
		var list []models.Media
		if err := r.db.Where("download_url = ?", dup.DownloadURL).Find(&list).Error; err != nil || len(list) < 2 {
			continue
		}
		// 主レコード（COMPLETEDまたは最初のもの）を決定
		var primary models.Media
		var secondaryIDs []string
		for _, m := range list {
			if primary.MediaID == "" {
				primary = m
			} else if m.DownloadStatus == "COMPLETED" && primary.DownloadStatus != "COMPLETED" {
				secondaryIDs = append(secondaryIDs, primary.MediaID)
				primary = m
			} else {
				secondaryIDs = append(secondaryIDs, m.MediaID)
			}
		}

		if len(secondaryIDs) > 0 {
			// 重複レコードのパージ
			if err := r.db.Where("media_id IN ?", secondaryIDs).Delete(&models.Media{}).Error; err == nil {
				mergedCount += len(secondaryIDs)
			}
		}
	}
	return mergedCount, nil
}

// PurgeLowerResolutionDuplicates は同一コンテンツの低解像度バリエーションを検出しゴミ箱へ退避します
func (r *Repository) PurgeLowerResolutionDuplicates(trashDir string) (int, error) {
	var allMedia []models.Media
	if err := r.db.Find(&allMedia).Error; err != nil {
		return 0, err
	}

	groups := make(map[string][]models.Media)
	for _, m := range allMedia {
		base := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		base = strings.Split(base, ":")[0] // :large, :small の除去
		if base != "" {
			groups[base] = append(groups[base], m)
		}
	}

	purgedCount := 0
	for _, group := range groups {
		if len(group) <= 1 {
			continue
		}
		// 最大解像度を持つものを探す
		maxIdx := 0
		maxArea := group[0].Width * group[0].Height
		for i, m := range group {
			area := m.Width * m.Height
			if area > maxArea {
				maxArea = area
				maxIdx = i
			}
		}

		for i, m := range group {
			if i == maxIdx {
				continue
			}
			// 低解像度版を退避
			filePath, err := r.ResolveMediaFilePath(m.MediaID)
			if err == nil && filePath != "" {
				_ = MoveToRecycleBin(filePath)
			}
			_ = r.db.Where("media_id = ?", m.MediaID).Delete(&models.Media{})
			purgedCount++
		}
	}
	return purgedCount, nil
}
