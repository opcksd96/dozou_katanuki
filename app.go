package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	stashManager *StashManager
}

func NewApp() *App {
	return &App{
		stashManager: NewStashManager(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 1. Stash のヘッドレスキックスタート
	stashBin := filepath.Join(".", "bin", "stash-win.exe")
	go func() {
		time.Sleep(500 * time.Millisecond) // 初期化マージン
		if err := a.stashManager.Start(ctx, stashBin); err != nil {
			runtime.EventsEmit(ctx, "stash:error", err.Error())
			return
		}
		// 疎通シグナル送信
		runtime.EventsEmit(ctx, "stash:status", map[string]interface{}{
			"status": "ONLINE",
			"port":   9999,
		})
	}()
}

func (a *App) shutdown(ctx context.Context) {
	// 2. 親ウィンドウ終了時に Stash を完全道連れ終了
	a.stashManager.Stop()
}
