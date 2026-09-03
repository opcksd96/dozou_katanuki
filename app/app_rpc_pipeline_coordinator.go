// app/app_rpc_pipeline_coordinator.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

// ExpandMediaCandidateTasks は 1つのメディアに対して原本URLおよび候補URLを展開してタスクレコード群を生成します
func ExpandMediaCandidateTasks(m models.Media, stage models.PipelineStage) []models.DownloadTask {
	if m.MediaID == "" { return nil }
	candidates := BuildCandidateURLsFromMediaWithArticle(m.MediaID, m.DownloadURL, m.Type, m.ArticleID)
	if len(candidates) == 0 {
		return []models.DownloadTask{
			{MediaID: m.MediaID, ArticleID: m.ArticleID, URL: m.DownloadURL, FileName: m.MediaID, Stage: stage, Status: models.TaskPending},
		}
	}
	tasks := make([]models.DownloadTask, 0, len(candidates))
	for _, c := range candidates {
		if c.URL == "" { continue }
		tasks = append(tasks, models.DownloadTask{
			MediaID: m.MediaID, ArticleID: m.ArticleID, URL: c.URL, FileName: m.MediaID, Stage: stage, Status: models.TaskPending,
		})
	}
	return tasks
}

// CoordinateTaskCompletion は 完了時の記録とStashパイプラインを実行します
func (a *App) CoordinateTaskCompletion(mediaID, completedFileName string, stage models.PipelineStage) {
	if a.Repo == nil || mediaID == "" { return }
	_ = a.Repo.UpdateMediaCheckpointTime(mediaID, stage)
	_ = a.Repo.UpdateMediaMetadata(mediaID, "COMPLETED", "", "", "")
	_ = a.Repo.MarkTaskCompleted(mediaID)
	a.AppendPipelineLog(string(stage), "SUCCESS", fmt.Sprintf("救出完了: %s (%s)", mediaID, completedFileName))
	a.emitToast("success", fmt.Sprintf("🎉 ダウンロード成功: %s", completedFileName))
	a.TriggerStashAllPipelines()
}

// CoordinateTaskDepletion は 特定ステージのタスク枯渇を記録し、次ステージエスカレーションを実行します
func (a *App) CoordinateTaskDepletion(fileName string, currentStage models.PipelineStage) {
	if a.Repo == nil || fileName == "" { return }
	mediaID := fileName

	// 現在のステージの候補が全滅した場合のエスカレーション制御
	switch currentStage {
	case models.StageRequests:
		// Requests 全滅 ➔ Motrix Next (第2ステージ) へエスカレーション
		_ = a.Repo.UpdateMediaMetadata(mediaID, "OUTSOURCED", "", "", "Requests枯渇: Motrix Nextへ移管")
		a.AppendPipelineLog("REQUESTS", "WARN", fmt.Sprintf("Requests全候補枯渇 ➔ Motrix移管: %s", mediaID))
		if m, err := a.Repo.GetMediaByID(mediaID); err == nil && m != nil {
			tasks := ExpandMediaCandidateTasks(*m, models.StageMotrix)
			_ = a.Repo.BatchUpsertDownloadTasks(tasks)
		}
	case models.StageMotrix:
		// Motrix 全滅 ➔ 迅雷 P2SP (第3ステージ) へエスカレーション
		_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", "Motrix枯渇: 迅雷P2SPへ自動エスカレーション")
		a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("Motrix全候補枯渇 ➔ 迅雷エスカレーション: %s", mediaID))
		if m, err := a.Repo.GetMediaByID(mediaID); err == nil && m != nil {
			tasks := ExpandMediaCandidateTasks(*m, models.StageThunder)
			_ = a.Repo.BatchUpsertDownloadTasks(tasks)
			_, _ = a.EscalateToThunder(m.MediaID, m.DownloadURL)
		}
	case models.StageThunder:
		// 迅雷も全滅 ➔ RETAINED (長期待機) へ退避
		_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "迅雷全候補枯渇: 司令塔リキュー待機")
		a.AppendPipelineLog("THUNDER", "ERROR", fmt.Sprintf("全候補全滅 ➔ RETAINED退避: %s", mediaID))
	}
}
