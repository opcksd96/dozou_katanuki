// main.go (100行以下)
package main

import (
	"embed"
	"net/url"

	"dozou_katanuki/middleware"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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

	err = wails.Run(&options.App{
		Title:  "dozou_katanuki",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: middleware.NewUnifiedHandler("./assets", stashURL), // アバター解決 & Stashリバースプロキシ
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
