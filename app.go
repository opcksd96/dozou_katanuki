package main

import (
	"context"
	"fmt"
	"log"

	"dozou_katanuki/driver"
	"dozou_katanuki/services"
)

// App struct
type App struct {
	ctx             context.Context
	timelineService *services.TimelineService
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
	a.timelineService = services.NewTimelineService(repo)
	log.Println("[App] Core services and middleware initialized successfully")
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	log.Println("[App] Application shutting down safely...")
}

// GetTimeline はフロントエンドへRenderTree配列を供給するWailsバインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]services.RenderTree, error) {
	if a.timelineService == nil {
		return nil, fmt.Errorf("timeline service is not initialized")
	}
	return a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
}
