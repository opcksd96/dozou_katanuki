// app/app_lifecycle.go (100行以下)
package app

import (
	"context"
	"log"
	"os"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
	db, err := driver.InitDB("archive.db")
	if err != nil { log.Fatalf("Failed to init database: %v", err) }

	a.Repo = driver.NewRepository(db)
	a.TimelineService = middleware.NewTimelineService(a.Repo)
	emitter := func(event string, data ...interface{}) { runtime.EventsEmit(ctx, event, data...) }
	a.JobOrchestrator = middleware.NewJobOrchestrator(ctx, emitter)
	if a.UnifiedHandler != nil { a.UnifiedHandler.SetJobOrchestrator(a.JobOrchestrator) }
	a.AuditService = middleware.NewAuditService(a.Repo, emitter)

	cfg, _ := a.GetConfig()
	schedCfg, netCfg, bcastCfg := models.SchedulerConfig{}, models.NetworkConfig{}, models.BroadcastConfig{}
	stashEnabled := true
	if cfg != nil {
		schedCfg, netCfg, bcastCfg = cfg.Scheduler, cfg.Network, cfg.Broadcast
		stashEnabled = cfg.Storage.StashEnabled
		a.JobOrchestrator.SetStorageConfig(cfg.Storage)
		if a.UnifiedHandler != nil {
			a.UnifiedHandler.SetMediaDir(cfg.Storage.LocalMediaDir)
		}
	}

	a.schedulerStart(ctx, schedCfg, emitter)
	a.broadcastStart(ctx, netCfg, bcastCfg, emitter)

	if stashEnabled {
		a.StashProber = middleware.NewStashProber("./bin/stash/stash-win.exe", "http://127.0.0.1:9999/", emitter)
		a.StashProber.Start(ctx)
	}

	// 起動時にwhitelist→accountsテーブルのgroup_name/alias_ofを一括同期
	if synced, err := a.Repo.SyncAllWhitelistGroups(); err == nil && synced > 0 {
		log.Printf("[Startup] Synced whitelist groups to %d accounts", synced)
	}

	a.ReadyOnce.Do(func() { close(a.Ready) })
	runtime.EventsEmit(ctx, "app:ready", true)
}

func (a *App) schedulerStart(ctx context.Context, schedCfg models.SchedulerConfig, emitter func(string, ...interface{})) {
	a.Scheduler = middleware.NewSchedulerService(schedCfg, a.Repo, a.JobOrchestrator, emitter)
	a.Scheduler.Start(ctx)
}

func (a *App) broadcastStart(ctx context.Context, netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig, emitter func(string, ...interface{})) {
	a.BroadcastService = middleware.NewBroadcastService(netCfg, bcastCfg, a.UnifiedHandler, a.TimelineService, emitter)
	if a.DistFS != nil {
		a.BroadcastService.SetDistFS(a.DistFS)
	}
	_ = a.BroadcastService.Start(ctx)
}

func (a *App) DomReady(ctx context.Context) {
	runtime.WindowShow(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if a.StashProber != nil { a.StashProber.Stop() }
		if a.BroadcastService != nil { _ = a.BroadcastService.Stop() }
		if a.Scheduler != nil { a.Scheduler.Stop() }
		if a.JobOrchestrator != nil { a.JobOrchestrator.Close() }
	}()

	select {
	case <-done:
	case <-time.After(300 * time.Millisecond):
	}
	os.Exit(0)
}
