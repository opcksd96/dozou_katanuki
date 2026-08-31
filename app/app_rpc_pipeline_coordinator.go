// app/app_rpc_pipeline_coordinator.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// ExpandMediaCandidateTasks は 1つのメディアに対して指定ステージの解像度候補タスク群を生成します
func ExpandMediaCandidateTasks(m models.Media, stage models.PipelineStage) []models.DownloadTask {
	var tasks []models.DownloadTask
	if m.DownloadURL == "" { return tasks }

	cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
	ext := ".jpg"
	if m.Type == "video" || strings.Contains(m.DownloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(m.MediaID), ".mp4") {
		ext = ".mp4"
	}

	candidates := BuildMediaCandidateURLs(m.DownloadURL)
	for _, c := range candidates {
		tID := fmt.Sprintf("%s-%s-%s", cleanID, strings.ToLower(string(stage)), c.Type)
		fileName := fmt.Sprintf("%s_%s%s", cleanID, c.Type, ext)

		tasks = append(tasks, models.DownloadTask{
			ID:             tID,
			MediaID:        m.MediaID,
			ArticleID:      m.ArticleID,
			Stage:          stage,
			ResolutionType: string(c.Type),
			URL:            c.URL,
			FileName:       fileName,
			Status:         models.TaskPending,
		})
	}
	return tasks
}

// CoordinateTaskCompletion は どのステージであれ1つの候補が完了した際、他全候補の刈り取りとStashパイプラインを実行します
func (a *App) CoordinateTaskCompletion(mediaID, completedFileName string, stage models.PipelineStage) {
	if a.Repo == nil || mediaID == "" { return }

	// 1. メディアDBのチェックポイント通過時刻とステータスを更新
	_ = a.Repo.UpdateMediaCheckpointTime(mediaID, stage)
	_ = a.Repo.UpdateMediaMetadata(mediaID, "COMPLETED", "", "", "")

	// 2. download_tasks 上で他全タスクを REAPED に一括更新
	reaps, _ := a.Repo.MarkTaskCompletedAndReapAllOthers(mediaID, completedFileName)

	// 3. 迅雷のキューからもサイレント削除
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err == nil && wsURL != "" {
		for _, task := range reaps {
			a.deleteTaskByFileNameSilent(wsURL, task.FileName)
		}
	}

	// 4. Stash 全自動パイプライン起動
	a.AppendPipelineLog(string(stage), "SUCCESS", fmt.Sprintf("救出完了: %s (%s)", mediaID, completedFileName))
	a.TriggerStashFullPipeline()
}

// CoordinateTaskDepletion は 特定ステージのタスク枯渇を記録し、全滅時の次ステージエスカレーションを実行します
func (a *App) CoordinateTaskDepletion(fileName string, currentStage models.PipelineStage) {
	if a.Repo == nil || fileName == "" { return }

	allDepleted, mediaID, err := a.Repo.MarkStageTaskDepleted(fileName, currentStage)
	if err != nil || !allDepleted || mediaID == "" { return }

	// 現在のステージの候補が全滅した場合のエスカレーション制御
	switch currentStage {
	case models.StageRequests, models.StageMotrix:
		// Motrix/Requests 全滅 ➔ 自動で迅雷へエスカレーション！
		_ = a.Repo.UpdateMediaCheckpointTime(mediaID, models.StageThunder)
		_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", "Motrix枯渇: 迅雷P2SPへ自動エスカレーション")
		a.AppendPipelineLog(string(currentStage), "WARN", fmt.Sprintf("全候補枯渇 ➔ 迅雷エスカレーション: %s", mediaID))
	case models.StageThunder:
		// 迅雷も全滅 ➔ RETAINED (長期待機) へ退避
		_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "迅雷全候補枯渇: 司令塔リキュー待機")
		a.AppendPipelineLog("THUNDER", "ERROR", fmt.Sprintf("全候補全滅 ➔ RETAINED退避: %s", mediaID))
	}
}
