// main.go (100行以下 - SPEC-PRINCIPLE-001)
package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"dozou_katanuki/app"
)

//go:embed all:frontend/dist
var assets embed.FS

const MainBuildRevision = "main-20260829-0800"

// UIApp struct for simple Wails bindings
type UIApp struct {
	ctx context.Context
}

func (u *UIApp) startup(ctx context.Context) {
	u.ctx = ctx
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		println("\n[Shutdown] 強制終了シグナルを受信しました。プロセスを終了します...")
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	uiApp := &UIApp{}
	
	// Create actual application logic core
	kApp := app.NewApp(nil, nil) // FIXME: middleware initialization needed if not injected later
	var appCtx context.Context

	appMenu := menu.NewMenu()
	toolsMenu := appMenu.AddSubmenu("設定・管理 (Admin)")
	toolsMenu.AddText("⚙️ 設定・Admin Board...", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		if appCtx != nil { runtime.EventsEmit(appCtx, "open:admin") }
	})

	viewMenu := appMenu.AddSubmenu("表示 (View)")
	viewMenu.AddText("再読み込み (Reload)", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		if appCtx != nil { runtime.WindowReload(appCtx) }
	})

	err := wails.Run(&options.App{
		Title:            "dozou_katanuki",
		Width:            1024,
		Height:           768,
		StartHidden:      true,
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 220},
		Menu:             appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			appCtx = ctx
			uiApp.startup(ctx)
			kApp.Startup(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			runtime.WindowShow(ctx)
			kApp.DomReady(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			kApp.Shutdown(ctx)
		},
		Windows: &windows.Options{
			Theme:                windows.Dark,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
		Bind: []interface{}{
			uiApp,
			kApp,
		},
	})

	if err != nil { println("Error:", err.Error()) }
}

