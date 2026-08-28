package driver

import (
	"dozou_katanuki/models"
)

// SaveDownloadReserve は退避メタデータをDBに保存または更新します
func (r *Repository) SaveDownloadReserve(reserve *models.DownloadReserve) error {
	if reserve == nil || r.db == nil {
		return nil
	}
	var existing models.DownloadReserve
	if reserve.GID != "" {
		if err := r.db.Where("gid = ?", reserve.GID).First(&existing).Error; err == nil {
			reserve.ID = existing.ID
			return r.db.Save(reserve).Error
		}
	}
	if reserve.URL != "" {
		if err := r.db.Where("url = ?", reserve.URL).First(&existing).Error; err == nil {
			reserve.ID = existing.ID
			return r.db.Save(reserve).Error
		}
	}
	return r.db.Create(reserve).Error
}

// FetchDownloadReserves は退避予約一覧を取得します
func (r *Repository) FetchDownloadReserves(status string, limit, offset int) ([]models.DownloadReserve, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	var list []models.DownloadReserve
	var total int64

	query := r.db.Model(&models.DownloadReserve{})
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 50
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

// UpdateReserveStatus は退避レコードのステータスを更新します
func (r *Repository) UpdateReserveStatus(id uint, status, reason string) error {
	if r.db == nil {
		return nil
	}
	updates := map[string]interface{}{"status": status}
	if reason != "" {
		updates["reason"] = reason
	}
	return r.db.Model(&models.DownloadReserve{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteDownloadReserve は退避レコードを削除します
func (r *Repository) DeleteDownloadReserve(id uint) error {
	if r.db == nil {
		return nil
	}
	return r.db.Delete(&models.DownloadReserve{}, id).Error
}
