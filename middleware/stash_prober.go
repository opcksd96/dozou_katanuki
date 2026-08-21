// middleware/stash_prober.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type StashProber struct {
	emitter    EventEmitter
	stashPath  string
	targetURL  string
	running    bool
	connected  bool
	cancelFunc context.CancelFunc
}

func NewStashProber(stashPath string, targetURL string, emitter EventEmitter) *StashProber {
	if targetURL == "" {
		targetURL = "http://127.0.0.1:9999/"
	}
	if stashPath == "" {
		stashPath = "./bin/stash-win.exe"
	}
	return &StashProber{
		stashPath: stashPath,
		targetURL: targetURL,
		emitter:   emitter,
	}
}

func (p *StashProber) IsConnected() bool {
	return p.connected
}

func (p *StashProber) Start(ctx context.Context) {
	proberCtx, cancel := context.WithCancel(ctx)
	p.cancelFunc = cancel
	p.running = true
	go p.probeLoop(proberCtx)
}

func (p *StashProber) probeLoop(ctx context.Context) {
	client := &http.Client{Timeout: 180 * time.Millisecond}
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	disconnectedSince := time.Now()
	lastKickAt := time.Time{}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := client.Get(p.targetURL)
			if err == nil {
				_ = resp.Body.Close()
				if !p.connected {
					p.connected = true
					if p.emitter != nil {
						p.emitter("stash:ready", true)
						p.emitter("toast:notify", map[string]string{
							"type":    "success",
							"message": "🟢 Stash 接続完了！",
						})
					}
					fmt.Println("[StashProber] 🟢 Stash 接続確認完了 (Ready)")
				}
			} else {
				if p.connected {
					p.connected = false
					disconnectedSince = time.Now()
					if p.emitter != nil {
						p.emitter("stash:ready", false)
					}
					fmt.Println("[StashProber] ⚠️ Stash 切断検知 (Reconnecting...)")
				}

				// 2秒以上切断状態が継続し、直近5秒以内にキックしていない場合
				if time.Since(disconnectedSince) > 2*time.Second && time.Since(lastKickAt) > 5*time.Second {
					if !p.isStashProcessAlive() {
						lastKickAt = time.Now()
						fmt.Println("[StashProber] 🔄 Stash プロセス不在を検知 ➔ 自動再キック実行")
						p.kickStashAsync()
					}
				}
			}
		}
	}
}

func (p *StashProber) isStashProcessAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "tasklist", "/FI", "IMAGENAME eq stash-win.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	s := strings.ToLower(string(out))
	return strings.Contains(s, "stash-win.exe")
}

func (p *StashProber) kickStashAsync() {
	if p.isStashProcessAlive() {
		fmt.Println("[StashProber] Stash プロセスが既に存在するため二重起動を抑止しました。")
		return
	}

	if p.emitter != nil {
		p.emitter("toast:notify", map[string]string{
			"type":    "info",
			"message": "📦 Stash をバックグラウンド起動中...",
		})
	}
	go func() {
		if p.isStashProcessAlive() {
			return
		}
		absPath, err := filepath.Abs(p.stashPath)
		if err != nil {
			return
		}
		cmd := exec.Command(absPath)
		cmd.Dir = filepath.Dir(absPath)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000 | 0x00000200}
		_ = cmd.Start()
		fmt.Printf("[StashProber] Stash プロセスをバックグラウンドキックしました: %s\n", absPath)
	}()
}

func (p *StashProber) Stop() {
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
	p.running = false
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "taskkill", "/F", "/IM", "stash-win.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = cmd.Run()
}
