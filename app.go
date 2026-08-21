// app.go (100行以下)
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx              context.Context
	repo             *driver.Repository
	timelineService  *middleware.TimelineService
	stashManager     *StashManager
	jobOrchestrator  *middleware.JobOrchestrator
	scheduler        *middleware.SchedulerService
	auditService     *middleware.AuditService
	broadcastService *middleware.BroadcastService
	unifiedHandler   *middleware.UnifiedHandler
	distFS           fs.FS
	ready            chan struct{}
	readyOnce        sync.Once
}

func NewApp(handler *middleware.UnifiedHandler, distFS fs.FS) *App {
	return &App{stashManager: NewStashManager(), unifiedHandler: handler, distFS: distFS, ready: make(chan struct{})}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := driver.InitDB("archive.db")
	if err != nil { log.Fatalf("Failed to init database: %v", err) }

	a.repo = driver.NewRepository(db)
	a.timelineService = middleware.NewTimelineService(a.repo)
	emitter := func(event string, data ...interface{}) { runtime.EventsEmit(ctx, event, data...) }
	a.jobOrchestrator = middleware.NewJobOrchestrator(ctx, emitter)
	if a.unifiedHandler != nil { a.unifiedHandler.SetJobOrchestrator(a.jobOrchestrator) }
	a.auditService = middleware.NewAuditService(a.repo, emitter)

	cfg, _ := a.GetConfig()
	schedCfg, netCfg, bcastCfg := models.SchedulerConfig{}, models.NetworkConfig{}, models.BroadcastConfig{}
	if cfg != nil { schedCfg, netCfg, bcastCfg = cfg.Scheduler, cfg.Network, cfg.Broadcast }

	a.scheduler = middleware.NewSchedulerService(schedCfg, a.repo, a.jobOrchestrator, emitter)
	a.scheduler.Start(ctx)

	a.broadcastService = middleware.NewBroadcastService(netCfg, bcastCfg, a.unifiedHandler, a.timelineService, emitter)
	if a.distFS != nil {
		a.broadcastService.SetDistFS(a.distFS)
	}
	_ = a.broadcastService.Start(ctx)

	a.readyOnce.Do(func() { close(a.ready) })
	runtime.EventsEmit(ctx, "app:ready", true)
	go func() {
		if a.stashManager != nil { _ = a.stashManager.Start(ctx, "./bin/stash-win.exe") }
	}()
}

func (a *App) domReady(ctx context.Context) {
	if a.stashManager != nil {
		_ = a.stashManager.WaitForReady(ctx, 8*time.Second)
	}
	runtime.WindowShow(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	if a.broadcastService != nil { _ = a.broadcastService.Stop() }
	if a.scheduler != nil { a.scheduler.Stop() }
	if a.jobOrchestrator != nil { a.jobOrchestrator.Close() }
	if a.stashManager != nil { a.stashManager.Stop() }
}

func (a *App) waitForReady() error {
	select {
	case <-a.ready:
		if a.timelineService == nil { return fmt.Errorf("timeline service not initialized") }
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for core initialization")
	}
}
