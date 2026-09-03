// app/app.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"context"
	"fmt"
	"io/fs"
	"sync"
	"time"

	"dozou_katanuki/application/account"
	"dozou_katanuki/application/timeline"
	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
)

const (
	AppVersion    = "v2.5.0"
	BuildRevision = "rev-20260829-0800"
)

type App struct {
	Ctx              context.Context
	Repo             *driver.Repository
	TimelineService  *middleware.TimelineService
	BeaconService    *middleware.BeaconService
	JobOrchestrator  *middleware.JobOrchestrator
	Scheduler        *middleware.SchedulerService
	AuditService     *middleware.AuditService
	BroadcastService *middleware.BroadcastService
	UnifiedHandler   *middleware.UnifiedHandler
	AccountUseCase   account.AccountUseCase
	TimelineUseCase  timeline.TimelineUseCase
	DistFS           fs.FS
	Ready            chan struct{}
	ReadyOnce        sync.Once
}

func NewApp(handler *middleware.UnifiedHandler, distFS fs.FS) *App {
	return &App{UnifiedHandler: handler, DistFS: distFS, Ready: make(chan struct{})}
}

func (a *App) GetAppVersion() string {
	return fmt.Sprintf("%s (%s)", AppVersion, BuildRevision)
}

func (a *App) IsStashReady() bool {
	if a.BeaconService != nil {
		return a.BeaconService.GetState().StashReady
	}
	return false
}

func (a *App) WaitForReady() error {
	select {
	case <-a.Ready:
		if a.TimelineService == nil {
			return fmt.Errorf("timeline service not initialized")
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for core initialization")
	}
}

func (a *App) GetTimelineService() *middleware.TimelineService {
	return a.TimelineService
}
