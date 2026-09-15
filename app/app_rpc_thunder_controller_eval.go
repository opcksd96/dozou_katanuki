// app/app_rpc_thunder_controller_eval.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// isFailedEmptyTask は 0B/0B かつ「ダウンロードが完全に停止している」タイムアウト以外の失敗タスクかを判定します
func isFailedEmptyTask(t models.ActiveThunderTask) bool {
	// 0. 拡張子ガード: Twitterメディア以外の外部ファイルは絶対に除外
	ext := strings.ToLower(filepath.Ext(t.FileName))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".mp4" {
		return false
	}

	// 1. ⚡ 最優先待機ガード: 排队等待・待機中のジョブは絶対に削除しない！
	lowerDetail := strings.ToLower(t.DetailText + " " + t.RawText + " " + t.Status)
	if isQueuedOrWaiting(t, lowerDetail) { return false }

	// 2. 既に一部でもダウンロード完了しているタスクは誤削除防止のため除外
	if t.DownloadSize > 0 { return false }

	// 3. ダウンロード停止ガード: 通信中・速度発生中タスクは絶対に除外
	if !isTaskFullyStopped(t, lowerDetail) { return false }

	// 4. タイムアウト除外 (ErrorCode: 5, または文言に「超时」「timeout」)
	if t.ErrorCode == 5 || strings.Contains(lowerDetail, "超时") || strings.Contains(lowerDetail, "timeout") {
		return false
	}

	// 5. 失敗・枯渇判定 (明確に停止済みのErrorCode 402, TaskStatusCode 9, 枯渇文言)
	if t.ErrorCode == 402 || t.TaskStatusCode == 9 { return true }
	if strings.Contains(t.Status, "枯渇") || strings.Contains(lowerDetail, "无法继续下载") ||
		strings.Contains(lowerDetail, "原始资源不存在") || strings.Contains(lowerDetail, "暂无任何有效资源") {
		return true
	}
	return false
}

// isQueuedOrWaiting は タスクが排队等待(キュー待機)や接続準備中かを判定します
func isQueuedOrWaiting(t models.ActiveThunderTask, lowerText string) bool {
	// 迅雷ステータス: 0 (待機/ピア接続/排队待ち), 5 (リソース探索中) はアクティブ待機
	if t.TaskStatusCode == 0 || t.TaskStatusCode == 5 { return true }
	// 文言検査: 排队や等待が含まれていれば待機中ジョブ
	if strings.Contains(lowerText, "排队") || strings.Contains(lowerText, "等待") ||
		strings.Contains(lowerText, "待机") || strings.Contains(lowerText, "队列") {
		return true
	}
	return false
}

// isTaskFullyStopped は タスクのダウンロードが完全に停止しているかを検証します
func isTaskFullyStopped(t models.ActiveThunderTask, lowerDetail string) bool {
	if t.DownloadSpeed > 0 || t.VIPSpeed > 0 { return false }
	if t.SpeedText != "" && !strings.HasPrefix(t.SpeedText, "0B") { return false }
	if strings.Contains(lowerDetail, "正在下载") || strings.Contains(lowerDetail, "下载中") { return false }

	// 迅雷のステータスコード: 9 (失敗/停止), 7, 8 (一時停止) のみ停止状態
	if t.TaskStatusCode == 9 || t.TaskStatusCode == 7 || t.TaskStatusCode == 8 {
		return true
	}
	return false
}
