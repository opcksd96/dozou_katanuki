// app/app_rpc_pipeline_metrics_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestCalculatePipelineMetrics(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "metrics_test.db")
	gdb, err := driver.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	sqlDB, _ := gdb.DB()
	defer sqlDB.Close()
	repo := driver.NewRepository(gdb)
	app := &App{Repo: repo}

	// 外部キー制約を満たすためダミーの Account と Article を登録
	_ = gdb.Create(&models.Account{NumericID: "acc_1", Username: "user1", DisplayName: "User 1", AvatarURL: "http://av"}).Error
	_ = gdb.Create(&models.Article{ID: "art_1", AccountID: "acc_1", ConversationID: "c1", FullText: "text", Via: "web", WaybackURL: "http://wb"}).Error

	// 10作品作成: 2作品ESCALATED, 7作品COMPLETED, 1作品RETAINED
	for i := 0; i < 7; i++ {
		_ = gdb.Create(&models.Media{
			MediaID: fmt.Sprintf("comp_%d", i), ArticleID: "art_1", Type: "video", DownloadURL: "http://example.com/1", DownloadStatus: "COMPLETED",
		}).Error
	}
	_ = gdb.Create(&models.Media{MediaID: "esc_1", ArticleID: "art_1", Type: "video", DownloadURL: "http://example.com/2", DownloadStatus: "ESCALATED"}).Error
	_ = gdb.Create(&models.Media{MediaID: "esc_2", ArticleID: "art_1", Type: "video", DownloadURL: "http://example.com/3", DownloadStatus: "ESCALATED"}).Error
	_ = gdb.Create(&models.Media{MediaID: "ret_1", ArticleID: "art_1", Type: "video", DownloadURL: "http://example.com/4", DownloadStatus: "RETAINED"}).Error

	// thunder_tasks
	// 2作品 * 7 = 14タスク
	// 稼働数(ONBOARDED): 3
	// 継続探査数(PENDING): 5
	// 成功数(REAPED/COMPLETED): 4
	// 失敗数(RETIRED): 2
	for i := 0; i < 3; i++ {
		_ = gdb.Create(&models.ThunderTask{
			ID: fmt.Sprintf("t_on_%d", i), MediaID: "esc_1", ArticleID: "art_1", ResolutionType: "1080p", URL: "http://a", FileName: fmt.Sprintf("on_%d.mp4", i), Status: models.ThunderTaskOnboarded,
		}).Error
	}
	for i := 0; i < 5; i++ {
		_ = gdb.Create(&models.ThunderTask{
			ID: fmt.Sprintf("t_pen_%d", i), MediaID: "esc_1", ArticleID: "art_1", ResolutionType: "720p", URL: "http://b", FileName: fmt.Sprintf("pen_%d.mp4", i), Status: models.ThunderTaskPending,
		}).Error
	}
	for i := 0; i < 4; i++ {
		_ = gdb.Create(&models.ThunderTask{
			ID: fmt.Sprintf("t_rea_%d", i), MediaID: "esc_2", ArticleID: "art_1", ResolutionType: "480p", URL: "http://c", FileName: fmt.Sprintf("rea_%d.mp4", i), Status: models.ThunderTaskReaped,
		}).Error
	}
	for i := 0; i < 2; i++ {
		_ = gdb.Create(&models.ThunderTask{
			ID: fmt.Sprintf("t_ret_%d", i), MediaID: "esc_2", ArticleID: "art_1", ResolutionType: "360p", URL: "http://d", FileName: fmt.Sprintf("ret_%d.mp4", i), Status: models.ThunderTaskRetired,
		}).Error
	}

	metrics := app.CalculatePipelineMetrics()

	if metrics.MediaCount != 2 {
		t.Errorf("expected MediaCount 2, got %d", metrics.MediaCount)
	}
	if metrics.TotalMediaCount != 10 {
		t.Errorf("expected TotalMediaCount 10, got %d", metrics.TotalMediaCount)
	}
	if metrics.TotalTasks != 14 {
		t.Errorf("expected TotalTasks 14, got %d", metrics.TotalTasks)
	}
	if metrics.ActiveTasks != 3 {
		t.Errorf("expected ActiveTasks 3, got %d", metrics.ActiveTasks)
	}
	if metrics.PendingTasks != 5 {
		t.Errorf("expected PendingTasks 5, got %d", metrics.PendingTasks)
	}
	if metrics.RegisteredTasks != 8 {
		t.Errorf("expected RegisteredTasks 8, got %d", metrics.RegisteredTasks)
	}
	if metrics.SuccessTasks != 4 {
		t.Errorf("expected SuccessTasks 4, got %d", metrics.SuccessTasks)
	}
	if metrics.FailedTasks != 2 {
		t.Errorf("expected FailedTasks 2, got %d", metrics.FailedTasks)
	}

	// 関係式の検証
	// タスク数 ＝ 下載登録数 + 成功数 + 失敗数
	if metrics.TotalTasks != metrics.RegisteredTasks+metrics.SuccessTasks+metrics.FailedTasks {
		t.Errorf("formula mismatch: %d != %d + %d + %d",
			metrics.TotalTasks, metrics.RegisteredTasks, metrics.SuccessTasks, metrics.FailedTasks)
	}
}
