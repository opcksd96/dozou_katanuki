package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App は Wails アプリケーションの基底構造体です
type App struct {
	ctx             context.Context
	timelineService *middleware.TimelineService
	stashManager    *StashManager
	ready           chan struct{}
	readyOnce       sync.Once
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		stashManager: NewStashManager(),
		ready:        make(chan struct{}),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 1. データベース初期化 (WALモード & AutoMigrate) を最優先で即時実行
	db, err := driver.InitDB("archive.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. リポジトリとミドルウェアサービスの組み立て
	repo := driver.NewRepository(db)
	a.timelineService = middleware.NewTimelineService(repo)
	log.Println("[App] Core services and middleware initialized successfully")

	// 3. バックエンドの初期化完了を通知（RPC待機解除）
	a.readyOnce.Do(func() {
		close(a.ready)
	})
	runtime.EventsEmit(ctx, "app:ready", true)

	// 4. Stash プロセスの自動ヘッドレス起動（非同期で実行し、DBアクセスやUI初期表示をブロックしない）
	go func() {
		if a.stashManager != nil {
			if err := a.stashManager.Start(ctx, "./bin/stash-win.exe"); err != nil {
				log.Printf("[App] StashManager 起動失敗 (非致命的): %v", err)
			}
		}
	}()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	log.Println("[App] Frontend DOM Ready. Displaying window smoothly...")
	runtime.WindowShow(ctx)
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	log.Println("[App] Application shutting down safely...")

	// Stash プロセスの道連れ終了 (Kill)
	if a.stashManager != nil {
		a.stashManager.Stop()
	}
}

// waitForReady はコアサービス（DB/TimelineService）の初期化完了を安全に待機します
func (a *App) waitForReady() error {
	select {
	case <-a.ready:
		if a.timelineService == nil {
			return fmt.Errorf("timeline service is not initialized")
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for core services initialization")
	}
}

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetAccounts(platform)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	res, err := a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}

// GetSystemLanguage は config.json (SPEC-CONFIG-001) からシステム言語設定を取得します
func (a *App) GetSystemLanguage() string {
	type Config struct {
		System struct {
			Language string `json:"language"`
		} `json:"system"`
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return "ja" // デフォルト言語
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.System.Language == "" {
		return "ja"
	}
	return cfg.System.Language
}
