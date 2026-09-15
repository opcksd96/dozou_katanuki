// driver/repo_thunder_tasks_reconcile_retained.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

// ReconcileDepletedMediaToRetained は 候補タスクがすべて全滅(REAPED/RETIRED/DEPLETED)しているのに
// download_status が ESCALATED のままになっているメディアを検出し、一括で RETAINED に退避します。
func (r *Repository) ReconcileDepletedMediaToRetained() (int64, error) {
	if r.db == nil {
		return 0, nil
	}

	query := `
		UPDATE media
		SET download_status = 'RETAINED',
		    failed_reason = '全候補枯渇(REAPED/RETIRED)によりRETAINED退避'
		WHERE download_status = 'ESCALATED'
		  AND media_id IN (
		      SELECT m.media_id
		      FROM media m
		      JOIN thunder_tasks t ON t.media_id = m.media_id
		      WHERE m.download_status = 'ESCALATED'
		      GROUP BY m.media_id
		      HAVING COUNT(t.id) > 0
		         AND SUM(CASE WHEN t.status = 'COMPLETED' THEN 1 ELSE 0 END) = 0
		         AND SUM(CASE WHEN t.status IN ('PENDING', 'ONBOARDED', 'RUNNING', 'HOLDING') THEN 1 ELSE 0 END) = 0
		         AND SUM(CASE WHEN t.status IN ('REAPED', 'RETIRED', 'DEPLETED') THEN 1 ELSE 0 END) = COUNT(t.id)
		  )
	`
	res := r.db.Exec(query)
	return res.RowsAffected, res.Error
}
