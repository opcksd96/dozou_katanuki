// main.go (100行以下)
package main

import (
	"embed"
	"net/url"

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
	app := NewApp()

	// Stashapp サーバー URL (ローカル閉塞 :9999)
	stashURL, err := url.Parse("http://127.0.0.1:9999")
	if err != nil {
		println("Stash URL Parse Error:", err.Error())
	}

	// OS ネイティブメニュー（Ctrl+R / F5 アクセラレータ登録）
	appMenu := menu.NewMenu()
	viewMenu := appMenu.AddSubmenu("View")
	viewMenu.AddText("Reload", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		if app.ctx != nil {
			runtime.WindowReload(app.ctx)
		}
	})
	viewMenu.AddText("Refresh", keys.Key("f5"), func(_ *menu.CallbackData) {
		if app.ctx != nil {
			runtime.WindowReload(app.ctx)
		}
	})

	err = wails.Run(&options.App{
		Title:            "dozou_katanuki",
		Width:            1024,
		Height:           768,
		StartHidden:      true,
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 255}, // slate-950
		Menu:             appMenu,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: middleware.NewUnifiedHandler("./assets", stashURL), // アバター解決 & Stashリバースプロキシ
		},
		OnStartup:  app.startup,
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
		Windows: &windows.Options{
			Theme:                windows.Dark,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
