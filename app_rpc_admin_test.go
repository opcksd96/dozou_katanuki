package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, string) {
	tempFile := "test_archive_" + time.Now().Format("20060102150405") + ".db"
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.Account{},
		&models.AccountProfileHistory{},
		&models.Article{},
		&models.Media{},
		&models.UrlRedirect{},
		&models.Whitelist{},
	)
	if err != nil {
		t.Fatalf("Failed to automigrate test db: %v", err)
	}

	return db, tempFile
}

func TestWhitelistCRUD(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)

	// 1. AddWhitelist
	item, err := repo.AddWhitelist("account", "mashu_dev")
	if err != nil {
		t.Fatalf("AddWhitelist failed: %v", err)
	}
	if item.ID == 0 || item.Value != "mashu_dev" || !item.IsActive {
		t.Fatalf("Unexpected item: %+v", item)
	}

	// 2. GetWhitelists
	list, err := repo.GetWhitelists()
	if err != nil || len(list) != 1 {
		t.Fatalf("GetWhitelists failed or unexpected len: %v, len: %d", err, len(list))
	}

	// 3. ToggleWhitelist
	if err := repo.ToggleWhitelist(item.ID); err != nil {
		t.Fatalf("ToggleWhitelist failed: %v", err)
	}
	list2, _ := repo.GetWhitelists()
	if list2[0].IsActive != false {
		t.Fatalf("Expected IsActive to be false, got true")
	}

	// 4. UpdateWhitelist
	if err := repo.UpdateWhitelist(item.ID, "keyword", "retro_famicom", true); err != nil {
		t.Fatalf("UpdateWhitelist failed: %v", err)
	}
	list3, _ := repo.GetWhitelists()
	if list3[0].Type != "keyword" || list3[0].Value != "retro_famicom" || !list3[0].IsActive {
		t.Fatalf("UpdateWhitelist not reflected: %+v", list3[0])
	}

	// 5. DeleteWhitelist
	if err := repo.DeleteWhitelist(item.ID); err != nil {
		t.Fatalf("DeleteWhitelist failed: %v", err)
	}
	list4, _ := repo.GetWhitelists()
	if len(list4) != 0 {
		t.Fatalf("Expected len 0 after deletion, got %d", len(list4))
	}
}

func TestArticleSearchAndTranslation(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)

	// アカウントと記事を作成
	acc := models.Account{
		NumericID:   "1001",
		Username:    "senpai_apu",
		DisplayName: "先輩マスター",
		AvatarURL:   "https://example.com/avatar.jpg",
		UpdatedAt:   time.Now(),
	}
	db.Create(&acc)

	art := models.Article{
		ID:             "art_001",
		AccountID:      "1001",
		ConversationID: "art_001",
		CreatedAt:      time.Now(),
		FullText:       "NES sound register $4000 test",
		Lang:           "ja",
		FullTextJA:     sql.NullString{String: "ファミコン音源レジスタ $4000 テスト", Valid: true},
		Via:            "Twitter for Web",
		IsRepost:       false,
		IsLiked:        true,
		WaybackURL:     "https://web.archive.org/web/...",
	}
	db.Create(&art)

	// 1. SearchArticles (query match in full_text)
	results, total, err := repo.SearchArticles("sound", "all", "all", 10, 0)
	if err != nil || total != 1 || len(results) != 1 {
		t.Fatalf("SearchArticles failed: err=%v, total=%d, len=%d", err, total, len(results))
	}

	// 2. SearchArticles (query match in translation)
	results2, total2, err := repo.SearchArticles("ファミコン", "all", "all", 10, 0)
	if err != nil || total2 != 1 || len(results2) != 1 {
		t.Fatalf("SearchArticles by translation failed: err=%v, total=%d, len=%d", err, total2, len(results2))
	}

	// 3. UpdateArticleTranslations
	err = repo.UpdateArticleTranslations("art_001", "ファミコン更新", "NES updated", "红白机更新")
	if err != nil {
		t.Fatalf("UpdateArticleTranslations failed: %v", err)
	}

	updatedArt, err := repo.GetArticleByID("art_001")
	if err != nil {
		t.Fatalf("GetArticleByID failed: %v", err)
	}
	if updatedArt.FullTextJA.String != "ファミコン更新" || updatedArt.FullTextEN.String != "NES updated" || updatedArt.FullTextZH.String != "红白机更新" {
		t.Fatalf("Translations not updated properly: %+v", updatedArt)
	}
}

func TestSchedulerRPC(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	readyChan := make(chan struct{})
	close(readyChan)

	app := &App{
		repo:            repo,
		timelineService: timeline,
		ready:           readyChan,
	}

	orch := middleware.NewJobOrchestrator(t.Context(), func(string, ...interface{}) {})
	defer orch.Close()
	app.jobOrchestrator = orch

	sched := middleware.NewSchedulerService(models.SchedulerConfig{
		PollIntervalSec:      10,
		BackupIntervalHours:  24,
		MaxBackupGenerations: 3,
	}, repo, orch, func(string, ...interface{}) {})
	app.scheduler = sched

	// 1. TriggerBackup RPC
	backupPath, err := app.TriggerBackup()
	if err != nil {
		t.Fatalf("TriggerBackup RPC failed: %v", err)
	}
	if backupPath == "" {
		t.Fatalf("Expected backup path, got empty")
	}
	defer os.Remove(backupPath)

	// 2. TriggerPoll RPC
	job, err := app.TriggerPoll()
	if err != nil {
		t.Fatalf("TriggerPoll RPC failed: %v", err)
	}
	if job == nil || job.Type != models.JobTypeMediaPoll {
		t.Fatalf("Expected JobTypeMediaPoll, got: %+v", job)
	}
}

func TestAuditRPC(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	readyChan := make(chan struct{})
	close(readyChan)

	app := &App{
		repo:            repo,
		timelineService: timeline,
		ready:           readyChan,
	}

	auditSvc := middleware.NewAuditService(repo, func(string, ...interface{}) {})
	app.auditService = auditSvc

	// 1. RunAudit
	report, err := app.RunAudit(false, false)
	if err != nil {
		t.Fatalf("RunAudit RPC failed: %v", err)
	}
	if report == nil || !report.IntegrityOK {
		t.Fatalf("Expected valid report with IntegrityOK true, got: %+v", report)
	}

	// 2. PurgeOrphanFiles
	purgedFiles, err := app.PurgeOrphanFiles([]string{})
	if err != nil || purgedFiles != 0 {
		t.Fatalf("PurgeOrphanFiles RPC failed: %v", err)
	}

	// 3. PurgeOrphanDBMedia
	purgedDB, err := app.PurgeOrphanDBMedia([]string{})
	if err != nil || purgedDB != 0 {
		t.Fatalf("PurgeOrphanDBMedia RPC failed: %v", err)
	}
}

func TestSkinCSSRPC(t *testing.T) {
	app := &App{}

	// 1. GetSkinCSS (twitter)
	css, err := app.GetSkinCSS("twitter")
	if err != nil {
		t.Fatalf("GetSkinCSS failed: %v", err)
	}
	if !strings.Contains(css, ".twitter-card") {
		t.Fatalf("GetSkinCSS expected '.twitter-card', got: %s", css)
	}

	// 2. SaveSkinCSS (テスト用カスタムプラットフォーム)
	testPlatform := "test_skin_platform"
	testCSS := "/* Test Skin CSS */\n.test-card { color: #abcdef; }\n"
	err = app.SaveSkinCSS(testPlatform, testCSS)
	if err != nil {
		t.Fatalf("SaveSkinCSS failed: %v", err)
	}
	defer os.RemoveAll(filepath.Join("plugins", testPlatform))

	// 3. GetSkinCSS で保存内容の読み出し確認
	readCSS, err := app.GetSkinCSS(testPlatform)
	if err != nil {
		t.Fatalf("GetSkinCSS for test platform failed: %v", err)
	}
	if readCSS != testCSS {
		t.Fatalf("GetSkinCSS content mismatch. expected: %s, got: %s", testCSS, readCSS)
	}
}

