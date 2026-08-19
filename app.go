// app.go (100行以下)
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

type App struct {
	ctx             context.Context
	timelineService *middleware.TimelineService
	stashManager    *StashManager
	jobOrchestrator *middleware.JobOrchestrator
	unifiedHandler  *middleware.UnifiedHandler
	ready           chan struct{}
	readyOnce       sync.Once
}

func NewApp(handler *middleware.UnifiedHandler) *App {
	return &App{
		stashManager:   NewStashManager(),
		unifiedHandler: handler,
		ready:          make(chan struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := driver.InitDB("archive.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	a.timelineService = middleware.NewTimelineService(driver.NewRepository(db))
	a.jobOrchestrator = middleware.NewJobOrchestrator(ctx, func(event string, data ...interface{}) {
		runtime.EventsEmit(ctx, event, data...)
	})
	if a.unifiedHandler != nil {
		a.unifiedHandler.SetJobOrchestrator(a.jobOrchestrator)
	}

	a.readyOnce.Do(func() { close(a.ready) })
	runtime.EventsEmit(ctx, "app:ready", true)

	go func() {
		if a.stashManager != nil {
			_ = a.stashManager.Start(ctx, "./bin/stash-win.exe")
		}
	}()
}

func (a *App) domReady(ctx context.Context) {
	log.Println("[App] Frontend DOM Ready. Stash サーバー起動待機中...")
	if a.stashManager != nil {
		if err := a.stashManager.WaitForReady(ctx, 8*time.Second); err != nil {
			log.Printf("[App] Stash 待機タイムアウト: %v", err)
		} else {
			log.Println("[App] Stash サーバー正常起動確認完了 (Ready)")
		}
	}
	runtime.WindowShow(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	log.Println("[App] Application shutting down...")
	if a.jobOrchestrator != nil {
		a.jobOrchestrator.Close()
	}
	if a.stashManager != nil {
		a.stashManager.Stop()
	}
}

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
