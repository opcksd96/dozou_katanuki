// driver/repo_thunder_tasks_counts.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

// ThunderTaskCounts は thunder_tasks テーブルの実態集計
type ThunderTaskCounts struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`   // RUNNING + ONBOARDED
	Holding   int `json:"holding"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`    // RETIRED + DEPLETED + REAPED
}

// GetThunderTaskCounts は thunder_tasks テーブルから実態のステータス別件数をリアルタイム集計します
func (r *Repository) GetThunderTaskCounts() (ThunderTaskCounts, error) {
	var counts ThunderTaskCounts
	if r.db == nil { return counts, nil }

	type statusRow struct {
		Status string
		Count  int
	}
	var rows []statusRow
	err := r.db.Raw("SELECT status, count(*) as count FROM thunder_tasks GROUP BY status").Scan(&rows).Error
	if err != nil { return counts, err }

	for _, row := range rows {
		counts.Total += row.Count
		switch row.Status {
		case "PENDING":
			counts.Pending += row.Count
		case "RUNNING", "ONBOARDED":
			counts.Running += row.Count
		case "HOLDING":
			counts.Holding += row.Count
		case "COMPLETED":
			counts.Completed += row.Count
		case "RETIRED", "DEPLETED", "REAPED":
			counts.Failed += row.Count
		}
	}
	return counts, nil
}
