// app/app_rpc_thunder_escalate_test.go
package app

import (
	"os"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestEscalateToThunder_EmptyURL(t *testing.T) {
	a := &App{}
	ok, err := a.EscalateToThunder("", "")
	if ok || err == nil {
		t.Errorf("expected error for empty URL, got ok=%v, err=%v", ok, err)
	}
}

func TestGiveUpRetainedMedia_EmptyID(t *testing.T) {
	a := &App{}
	ok, err := a.GiveUpRetainedMedia("")
	if ok || err == nil {
		t.Errorf("expected error for empty mediaID, got ok=%v, err=%v", ok, err)
	}
}

func TestEscalateToThunder_PipelineIntegration(t *testing.T) {
	dbPath := "./test_thunder_pipeline.db"
	_ = os.Remove(dbPath)
	defer os.Remove(dbPath)

	gormDB, err := driver.InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	repo := driver.NewRepository(gormDB)

	// 1. テスト用のアカウント・記事・メディアレコードを作成
	acc := models.Account{
		NumericID:   "test_user_001",
		Username:    "test_user",
		DisplayName: "Test User",
		AvatarURL:   "avatar.jpg",
		UpdatedAt:   time.Now(),
	}
	if err := gormDB.Create(&acc).Error; err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	art := models.Article{
		ID:         "test_article_001",
		AccountID:  "test_user_001",
		CreatedAt:  time.Now(),
		FullText:   "Thunder pipeline test",
		Lang:       "ja",
		WaybackURL: "http://example.com",
	}
	if err := gormDB.Create(&art).Error; err != nil {
		t.Fatalf("failed to create article: %v", err)
	}

	testMediaID := "test_thunder_media_001.jpg"
	testURL := "https://pbs.twimg.com/media/Go4PdTDb0AAmdRV.jpg"
	m := models.Media{
		MediaID:        testMediaID,
		ArticleID:      "test_article_001",
		Type:           "image",
		DownloadURL:    testURL,
		DownloadStatus: "RETAINED",
	}
	if err := gormDB.Create(&m).Error; err != nil {
		t.Fatalf("failed to create test media: %v", err)
	}

	app := &App{Repo: repo}

	// 2. RETAINED から OUTSOURCED に変更
	err = repo.UpdateMediaMetadata(testMediaID, "OUTSOURCED", "", "", "Motrix外注中テスト")
	if err != nil {
		t.Fatalf("failed to set OUTSOURCED: %v", err)
	}
	mOutsourced, err := repo.GetMediaByID(testMediaID)
	if err != nil || mOutsourced.DownloadStatus != "OUTSOURCED" {
		t.Fatalf("expected OUTSOURCED, got %s (err: %v)", mOutsourced.DownloadStatus, err)
	}

	// 3. EscalateToThunder を呼び出して Thunder へ昇格投入
	ok, err := app.EscalateToThunder(testMediaID, testURL)
	if !ok || err != nil {
		t.Fatalf("EscalateToThunder failed: ok=%v, err=%v", ok, err)
	}

	// 4. ステータスが ESCALATED に昇格したことを検証
	mEscalated, err := repo.GetMediaByID(testMediaID)
	if err != nil || mEscalated.DownloadStatus != "ESCALATED" {
		t.Fatalf("expected status ESCALATED, got %s (err: %v)", mEscalated.DownloadStatus, err)
	}
	t.Logf("Successfully escalated media %s to ESCALATED with Thunder invocation!", testMediaID)
}
