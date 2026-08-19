package main

import (
	"context"
	"fmt"
	"log"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

// App は Wails アプリケーションの基底構造体です
type App struct {
	ctx             context.Context
	timelineService *middleware.TimelineService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 1. データベース初期化 (WALモード & AutoMigrate)
	db, err := driver.InitDB("archive.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. リポジトリとミドルウェアサービスの組み立て
	repo := driver.NewRepository(db)
	a.timelineService = middleware.NewTimelineService(repo)
	log.Println("[App] Core services and middleware initialized successfully")
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	log.Println("[App] Application shutting down safely...")
}

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	if a.timelineService == nil {
		return nil, fmt.Errorf("timeline service is not initialized")
	}
	return a.timelineService.GetAccounts(platform)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if a.timelineService == nil {
		return nil, fmt.Errorf("timeline service is not initialized")
	}
	res, err := a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}
