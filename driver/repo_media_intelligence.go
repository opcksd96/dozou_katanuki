// driver/repo_media_intelligence.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// MergeDuplicateMedia は同一ファイル名、重複URL、サフィックス付き(:large等)メディアを自動統合します
func (r *Repository) MergeDuplicateMedia() (int, error) {
	var suffixMedia []models.Media
	r.db.Where("media_id LIKE ? OR media_id LIKE ? OR media_id LIKE ?", "%:large", "%:orig", "%:small").Find(&suffixMedia)
	mergedCount := 0
	for _, m := range suffixMedia {
		baseID := strings.Split(m.MediaID, ":")[0]
		var base models.Media
		if err := r.db.Where("media_id = ?", baseID).First(&base).Error; err == nil {
			if m.DownloadStatus == "COMPLETED" && base.DownloadStatus != "COMPLETED" {
				base.DownloadStatus = "COMPLETED"
				base.StashImageID = m.StashImageID
				base.StashSceneID = m.StashSceneID
				r.db.Save(&base)
			}
			r.db.Where("media_id = ?", m.MediaID).Delete(&models.Media{})
			mergedCount++
		} else {
			r.db.Model(&models.Media{}).Where("media_id = ?", m.MediaID).Update("media_id", baseID)
			mergedCount++
		}
	}

	var duplicates []struct {
		DownloadURL string
		Count       int
	}
	if err := r.db.Model(&models.Media{}).Select("download_url, count(*) as count").Where("download_url != ''").Group("download_url").Having("count(*) > 1").Scan(&duplicates).Error; err == nil {
		for _, dup := range duplicates {
			var list []models.Media
			if err := r.db.Where("download_url = ?", dup.DownloadURL).Find(&list).Error; err == nil && len(list) >= 2 {
				var primary models.Media
				var secondaryIDs []string
				for _, m := range list {
					if primary.MediaID == "" { primary = m } else if m.DownloadStatus == "COMPLETED" && primary.DownloadStatus != "COMPLETED" {
						secondaryIDs = append(secondaryIDs, primary.MediaID); primary = m
					} else { secondaryIDs = append(secondaryIDs, m.MediaID) }
				}
				if len(secondaryIDs) > 0 {
					if err := r.db.Where("media_id IN ?", secondaryIDs).Delete(&models.Media{}).Error; err == nil {
						mergedCount += len(secondaryIDs)
					}
				}
			}
		}
	}
	return mergedCount, nil
}

// PurgeLowerResolutionDuplicates は同一コンテンツの低解像度バリエーションを検出しゴミ箱へ退避します
func (r *Repository) PurgeLowerResolutionDuplicates(trashDir string) (int, error) {
	var allMedia []models.Media
	if err := r.db.Find(&allMedia).Error; err != nil { return 0, err }
	groups := make(map[string][]models.Media)
	for _, m := range allMedia {
		base := strings.Split(strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID)), ":")[0]
		if base != "" { groups[base] = append(groups[base], m) }
	}
	purgedCount := 0
	for _, group := range groups {
		if len(group) <= 1 { continue }
		bestIdx := 0
		for i := 1; i < len(group); i++ {
			if (group[i].Width * group[i].Height) > (group[bestIdx].Width * group[bestIdx].Height) { bestIdx = i }
		}
		for i, m := range group {
			if i != bestIdx {
				r.db.Model(&models.Media{}).Where("media_id = ?", m.MediaID).Update("download_status", "EXCLUDED")
				purgedCount++
			}
		}
	}
	return purgedCount, nil
}
