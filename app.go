package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App は Wails アプリケーションの基底構造体です
type App struct {
	ctx             context.Context
	timelineService *middleware.TimelineService
	stashManager    *StashManager
	jobOrchestrator *middleware.JobOrchestrator
	unifiedHandler  *middleware.UnifiedHandler
	ready           chan struct{}
	readyOnce       sync.Once
}

// NewApp creates a new App application struct
func NewApp(handler *middleware.UnifiedHandler) *App {
	return &App{
		stashManager:   NewStashManager(),
		unifiedHandler: handler,
		ready:          make(chan struct{}),
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

	// 3. JobOrchestrator の初期化と Wails イベント発行コールバックのバインド
	a.jobOrchestrator = middleware.NewJobOrchestrator(ctx, func(eventName string, optionalData ...interface{}) {
		runtime.EventsEmit(ctx, eventName, optionalData...)
	})
	if a.unifiedHandler != nil {
		a.unifiedHandler.SetJobOrchestrator(a.jobOrchestrator)
	}
	log.Println("[App] Core services, middleware, and JobOrchestrator initialized successfully")

	// 4. バックエンドの初期化完了を通知（RPC待機解除）
	a.readyOnce.Do(func() {
		close(a.ready)
	})
	runtime.EventsEmit(ctx, "app:ready", true)

	// 5. Stash プロセスの自動ヘッドレス起動（非同期で実行し、DBアクセスやUI初期表示をブロックしない）
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
	log.Println("[App] Frontend DOM Ready. Stash サーバーの起動待機中...")

	// 🌟 Stash サーバー (:9999) が正常応答するまで待機（最大8秒）
	// これにより、起動直後のメディア 502 / ECONNREFUSED を完全防止し、完璧な状態でウィンドウを表示する
	if a.stashManager != nil {
		if err := a.stashManager.WaitForReady(ctx, 8*time.Second); err != nil {
			log.Printf("[App] Stash 待機タイムアウトまたは未起動 (ウィンドウをフォールバック表示): %v", err)
		} else {
			log.Println("[App] Stash サーバー疎通確認完了 (Ready)。ウィンドウを滑らかに表示します。")
		}
	}

	runtime.WindowShow(ctx)
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	log.Println("[App] Application shutting down safely...")

	// JobOrchestrator の安全停止
	if a.jobOrchestrator != nil {
		a.jobOrchestrator.Close()
	}

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
