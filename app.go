// app.go (100行以下)
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
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
	stashProber      *middleware.StashProber
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
	return &App{unifiedHandler: handler, distFS: distFS, ready: make(chan struct{})}
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
	stashEnabled := true
	if cfg != nil {
		schedCfg, netCfg, bcastCfg = cfg.Scheduler, cfg.Network, cfg.Broadcast
		stashEnabled = cfg.Storage.StashEnabled
		a.jobOrchestrator.SetStorageConfig(cfg.Storage)
		if a.unifiedHandler != nil {
			a.unifiedHandler.SetMediaDir(cfg.Storage.LocalMediaDir)
		}
	}

	a.scheduler = middleware.NewSchedulerService(schedCfg, a.repo, a.jobOrchestrator, emitter)
	a.scheduler.Start(ctx)

	a.broadcastService = middleware.NewBroadcastService(netCfg, bcastCfg, a.unifiedHandler, a.timelineService, emitter)
	if a.distFS != nil {
		a.broadcastService.SetDistFS(a.distFS)
	}
	_ = a.broadcastService.Start(ctx)

	if stashEnabled {
		a.stashProber = middleware.NewStashProber("./bin/stash/stash-win.exe", "http://127.0.0.1:9999/", emitter)
		a.stashProber.Start(ctx)
	}

	a.readyOnce.Do(func() { close(a.ready) })
	runtime.EventsEmit(ctx, "app:ready", true)
}

func (a *App) domReady(ctx context.Context) {
	runtime.WindowShow(ctx)
}

func (a *App) IsStashReady() bool {
	if a.stashProber != nil {
		return a.stashProber.IsConnected()
	}
	return false
}

func (a *App) shutdown(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if a.stashProber != nil { a.stashProber.Stop() }
		if a.broadcastService != nil { _ = a.broadcastService.Stop() }
		if a.scheduler != nil { a.scheduler.Stop() }
		if a.jobOrchestrator != nil { a.jobOrchestrator.Close() }
	}()

	select {
	case <-done:
	case <-time.After(300 * time.Millisecond):
	}
	os.Exit(0)
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
