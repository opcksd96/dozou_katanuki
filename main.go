package main

import (
	"embed"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Stash (:9999) へのインメモリ・リバースプロキシターゲット
	stashURL, _ := url.Parse("http://127.0.0.1:9999")
	proxy := httputil.NewSingleHostReverseProxy(stashURL)

	err := wails.Run(&options.App{
		Title:  "dozou_katanuki (Chimera PoC)",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
			// インメモリ・リバースプロキシ（外部TCPポートを全廃し、メモリ内でStashへ中継）
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/stash-proxy/") {
					// パス書き換え: /stash-proxy/xxx -> /xxx
					r.URL.Path = strings.TrimPrefix(r.URL.Path, "/stash-proxy")
					proxy.ServeHTTP(w, r)
					return
				}
				http.NotFound(w, r)
			}),
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
