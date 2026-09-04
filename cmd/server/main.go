// cmd/server/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"dozou_katanuki/app"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

func main() {
	log.Println("[Boot] Katanuki Core API Server is starting in headless mode...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		log.Println("\n[Shutdown] 強制終了シグナルを受信しました。プロセスを終了します...")
		cancel()
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	port := 9999
	if data, err := os.ReadFile("config.json"); err == nil {
		var cfg models.AppConfig
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.Network.StashPort > 0 {
			port = cfg.Network.StashPort
		}
	}

	stashURL, err := url.Parse("http://127.0.0.1:" + strconv.Itoa(port))
	if err != nil {
		log.Println("Stash URL Parse Error:", err.Error())
	}

	unifiedHandler := middleware.NewUnifiedHandler("./assets", stashURL)
	appInstance := app.NewApp(unifiedHandler, nil)
	appInstance.IsHeadless = true

	appInstance.Startup(ctx)
	log.Println("\n[Ready] Katanuki Core API Server is running in headless mode. (Press Ctrl+C to stop)")

	<-ctx.Done()
	appInstance.Shutdown(ctx)
}
