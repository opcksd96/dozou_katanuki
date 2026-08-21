// main.go (100行以下)
package main

import (
	"embed"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dozou_katanuki/middleware"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Ctrl+C / SIGINT / SIGTERM を検知して即座に終了するセーフガード
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		println("\n[Shutdown] 強制終了シグナルを受信しました。プロセスを終了します...")
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	stashURL, err := url.Parse("http://127.0.0.1:9999")
	if err != nil {
		println("Stash URL Parse Error:", err.Error())
	}

	unifiedHandler := middleware.NewUnifiedHandler("./assets", stashURL)
	app := NewApp(unifiedHandler, assets)

	// OS ネイティブメニュー
	appMenu := menu.NewMenu()
	toolsMenu := appMenu.AddSubmenu("設定・管理 (Admin)")
	toolsMenu.AddText("⚙️ 設定・Admin Board...", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "open:admin")
		}
	})

	viewMenu := appMenu.AddSubmenu("表示 (View)")
	viewMenu.AddText("再読み込み (Reload)", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		if app.ctx != nil {
			runtime.WindowReload(app.ctx)
		}
	})

	err = wails.Run(&options.App{
		Title:            "dozou_katanuki",
		Width:            1024,
		Height:           768,
		StartHidden:      true,
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 255},
		Menu:             appMenu,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: unifiedHandler,
		},
		OnStartup:  app.startup,
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
		Windows: &windows.Options{
			Theme:                windows.Dark,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Bind: []interface{}{app},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
