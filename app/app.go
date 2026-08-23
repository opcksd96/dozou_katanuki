// app/app.go (100行以下)
package app

import (
	"context"
	"fmt"
	"io/fs"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
)

type App struct {
	Ctx              context.Context
	Repo             *driver.Repository
	TimelineService  *middleware.TimelineService
	StashProber      *middleware.StashProber
	JobOrchestrator  *middleware.JobOrchestrator
	Scheduler        *middleware.SchedulerService
	AuditService     *middleware.AuditService
	BroadcastService *middleware.BroadcastService
	UnifiedHandler   *middleware.UnifiedHandler
	DistFS           fs.FS
	Ready            chan struct{}
	ReadyOnce        sync.Once
}

func NewApp(handler *middleware.UnifiedHandler, distFS fs.FS) *App {
	return &App{UnifiedHandler: handler, DistFS: distFS, Ready: make(chan struct{})}
}

func (a *App) IsStashReady() bool {
	if a.StashProber != nil {
		return a.StashProber.IsConnected()
	}
	return false
}

func (a *App) WaitForReady() error {
	select {
	case <-a.Ready:
		if a.TimelineService == nil { return fmt.Errorf("timeline service not initialized") }
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for core initialization")
	}
}

func (a *App) GetTimelineService() *middleware.TimelineService {
	return a.TimelineService
}
