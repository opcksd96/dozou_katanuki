// driver/repo_media_reset.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"dozou_katanuki/models"
)

// ResetVideosToRetained は失敗・エスカレーション中の動画メディアを RETAINED ステータスへ差し戻します
func (r *Repository) ResetVideosToRetained() (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	res := r.db.Model(&models.Media{}).
		Where("type != 'image' AND (download_status = 'ESCALATED' OR download_status = 'DEAD_404' OR download_status = 'FAILED' OR download_status = 'QUEUED')").
		Updates(map[string]interface{}{
			"download_status": "RETAINED",
			"failed_reason":   "ユーザー指示による動画差し戻し (Reset to RETAINED)",
		})
	return res.RowsAffected, res.Error
}

// ResetAllFailedToRetained は全失敗メディアを RETAINED に差し戻します
func (r *Repository) ResetAllFailedToRetained() (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	res := r.db.Model(&models.Media{}).
		Where("download_status = 'ESCALATED' OR download_status = 'DEAD_404' OR download_status = 'FAILED'").
		Updates(map[string]interface{}{
			"download_status": "RETAINED",
			"failed_reason":   "エスカレーション失敗タスクの再投入差し戻し",
		})
	return res.RowsAffected, res.Error
}
