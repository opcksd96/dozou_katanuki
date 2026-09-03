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
	
	"dozou_katanuki/adapters/driven/sqlite"
	"dozou_katanuki/application/account"
	"dozou_katanuki/application/timeline"
	"dozou_katanuki/domain/ports"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Startup(ctx context.Context) {
	// Wailsアプリ(ローカルプロセス)からの呼び出しはすべて管理者権限として扱う
	a.Ctx = ports.WithScope(ctx, ports.ScopeAdmin)
	
	db, err := driver.InitDB("archive.db")
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	a.Repo = driver.NewRepository(db)
	
	// Initialize Hexagonal Architecture Ports/UseCases
	accountRepo := sqlite.NewAccountRepositoryImpl(db)
	articleRepo := sqlite.NewArticleRepositoryImpl(db)
	mediaRepo := sqlite.NewMediaRepositoryImpl(db)
	
	a.AccountUseCase = account.NewAccountUseCase(accountRepo)
	a.TimelineUseCase = timeline.NewTimelineUseCase(articleRepo, mediaRepo)
	
	a.TimelineService = middleware.NewTimelineService(a.Repo)
	emitter := func(event string, data ...interface{}) {
		if a.Ctx != nil && !a.IsHeadless {
			runtime.EventsEmit(a.Ctx, event, data...)
		}
		if a.BroadcastService != nil {
			a.BroadcastService.PushEvent(event, data...)
		}
	}
	a.JobOrchestrator = middleware.NewJobOrchestrator(ctx, emitter)
	if a.UnifiedHandler != nil {
		a.UnifiedHandler.SetJobOrchestrator(a.JobOrchestrator)
	}
	a.AuditService = middleware.NewAuditService(a.Repo, emitter)

	cfg, _ := a.GetConfig()
	schedCfg, netCfg, bcastCfg := models.SchedulerConfig{}, models.NetworkConfig{}, models.BroadcastConfig{}
	if cfg != nil {
		schedCfg, netCfg, bcastCfg = cfg.Scheduler, cfg.Network, cfg.Broadcast
		a.JobOrchestrator.SetStorageConfig(cfg.Storage)
		if a.UnifiedHandler != nil {
			a.UnifiedHandler.SetMediaDir(cfg.Storage.LocalMediaDir)
		}
	}

	a.schedulerStart(ctx, schedCfg, emitter)
	a.broadcastStart(ctx, netCfg, bcastCfg, emitter)

	a.BeaconService = middleware.NewBeaconService(emitter)
	a.BeaconService.Start(ctx)

	// 起動時にwhitelist→accountsテーブルのgroup_name/alias_ofを一括同期
	if synced, err := a.Repo.SyncAllWhitelistGroups(); err == nil && synced > 0 {
		log.Printf("[Startup] Synced whitelist groups to %d accounts", synced)
	}

	a.ReadyOnce.Do(func() { close(a.Ready) })
	a.EmitEvent("app:ready", true)
}

func (a *App) EmitEvent(eventName string, data ...interface{}) {
	if a.Ctx != nil && !a.IsHeadless {
		runtime.EventsEmit(a.Ctx, eventName, data...)
	}
	if a.BroadcastService != nil {
		a.BroadcastService.PushEvent(eventName, data...)
	}
}

func (a *App) schedulerStart(ctx context.Context, schedCfg models.SchedulerConfig, emitter func(string, ...interface{})) {
	a.Scheduler = middleware.NewSchedulerService(schedCfg, a.Repo, a.JobOrchestrator, emitter)
	a.Scheduler.Start(ctx)
}

type adminUseCaseWrapper struct {
	app *App
}

func (w *adminUseCaseWrapper) TogglePipelineAutoEngine(enable bool) (bool, error) {
	return w.app.TogglePipelineAutoEngine(enable)
}
func (w *adminUseCaseWrapper) IsPipelineAutoEngineRunning() bool {
	return w.app.IsPipelineAutoEngineRunning()
}
func (w *adminUseCaseWrapper) GetPipelineOverview() (interface{}, error) {
	return w.app.GetPipelineOverview()
}
func (w *adminUseCaseWrapper) GetPipelineLogs(stage string, limit int) (interface{}, error) {
	return w.app.GetPipelineLogs(stage, limit)
}
func (w *adminUseCaseWrapper) SyncThunderDownloads(req string) (interface{}, error) {
	return w.app.SyncThunderDownloads(req)
}
func (w *adminUseCaseWrapper) ResetAllToQueuedAndBootstrap() (interface{}, error) {
	return w.app.ResetAllToQueuedAndBootstrap()
}
func (w *adminUseCaseWrapper) IgnitePipeline() (interface{}, error) {
	return w.app.IgnitePipeline()
}

func (a *App) broadcastStart(ctx context.Context, netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig, emitter func(string, ...interface{})) {
	a.BroadcastService = middleware.NewBroadcastService(netCfg, bcastCfg, a.UnifiedHandler, a.TimelineService, emitter)
	a.BroadcastService.SetAdminUseCases(&adminUseCaseWrapper{app: a})
	if a.DistFS != nil {
		a.BroadcastService.SetDistFS(a.DistFS)
	}
	_ = a.BroadcastService.Start(ctx)
}

func (a *App) DomReady(ctx context.Context) {
	runtime.WindowShow(ctx)
	// CDPチェックは外部(PS1)からのビーコンに委譲したため、ここでの自動チェックと起動は廃止
}

func (a *App) Shutdown(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		// BeaconService relies on ctx cancellation, no Stop() needed.
		if a.BroadcastService != nil {
			_ = a.BroadcastService.Stop()
		}
		if a.Scheduler != nil {
			a.Scheduler.Stop()
		}
		if a.JobOrchestrator != nil {
			a.JobOrchestrator.Close()
		}
	}()

	select {
	case <-done:
	case <-time.After(300 * time.Millisecond):
	}
	os.Exit(0)
}
