// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dozou_katanuki/app"
	"dozou_katanuki/middleware"
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

	stashURL, err := url.Parse("http://127.0.0.1:9999")
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
