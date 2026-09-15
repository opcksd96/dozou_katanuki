// app/app_rpc_thunder_controller_api_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/models"
)

func TestIsFailedEmptyTask(t *testing.T) {
	tests := []struct {
		name     string
		task     models.ActiveThunderTask
		expected bool
	}{
		{
			name: "0B/0B 停止中 ErrorCode 402 (.mp4) -> 対象",
			task: models.ActiveThunderTask{
				FileName: "sample_video.mp4", FileSize: 0, DownloadSize: 0,
				ErrorCode: 402, TaskStatusCode: 9, DownloadSpeed: 0,
				DetailText: "原始资源不存在，且未找到候选资源，无法继续下载",
			},
			expected: true,
		},
		{
			name: "⚡ 排队等待 (キュー待機中ジョブ) -> 絶対に除外",
			task: models.ActiveThunderTask{
				FileName: "queued_video.mp4", FileSize: 0, DownloadSize: 0,
				TaskStatusCode: 0, Status: "排队等待", DetailText: "排队中",
			},
			expected: false,
		},
		{
			name: "⚡ 待機中 (TaskStatusCode: 0, ピア待機) -> 絶対に除外",
			task: models.ActiveThunderTask{
				FileName: "waiting_video.mp4", FileSize: 0, DownloadSize: 0,
				TaskStatusCode: 0, Status: "⚪ 待機中", DetailText: "等待中",
			},
			expected: false,
		},
		{
			name: "⚡ msluo14.zip 等の外部zipファイル -> 絶対に除外",
			task: models.ActiveThunderTask{
				FileName: "msluo14.zip", FileSize: 0, DownloadSize: 0,
				ErrorCode: 402, TaskStatusCode: 9, DownloadSpeed: 0,
			},
			expected: false,
		},
		{
			name: "0B/0B だが ダウンロード速度発生中 (DownloadSpeed > 0) -> 除外",
			task: models.ActiveThunderTask{
				FileName: "video.mp4", FileSize: 0, DownloadSize: 0, ErrorCode: 402,
				TaskStatusCode: 9, DownloadSpeed: 1024,
			},
			expected: false,
		},
		{
			name: "0B/0B だが SpeedText が 0B/s 以外 -> 除外",
			task: models.ActiveThunderTask{
				FileName: "image.jpg", FileSize: 0, DownloadSize: 0, ErrorCode: 402,
				TaskStatusCode: 9, SpeedText: "500KB/s",
			},
			expected: false,
		},
		{
			name: "0B/0B だが 文言に 正在下载 含む -> 除外",
			task: models.ActiveThunderTask{
				FileName: "video.mp4", FileSize: 0, DownloadSize: 0, ErrorCode: 402,
				TaskStatusCode: 9, DetailText: "正在下载中", Status: "ダウンロード中",
			},
			expected: false,
		},
		{
			name: "0B/0B タイムアウト (ErrorCode: 5) -> 除外",
			task: models.ActiveThunderTask{
				FileName: "video.mp4", FileSize: 0, DownloadSize: 0, ErrorCode: 5,
				TaskStatusCode: 9, DetailText: "任务连接超时，无法继续下载",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isFailedEmptyTask(tt.task)
			if got != tt.expected {
				t.Errorf("isFailedEmptyTask() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
