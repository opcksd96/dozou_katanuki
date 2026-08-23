// middleware/stash_prober.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	if targetURL == "" { targetURL = "http://127.0.0.1:9999/" }
	candidates := []string{
		stashPath,
		filepath.Join("bin", "stash", "stash-win.exe"),
		filepath.Join("bin", "stash-win.exe"),
		"stash-win.exe",
	}
	resolvedPath := ""
	for _, c := range candidates {
		if c != "" {
			if _, err := os.Stat(c); err == nil {
				resolvedPath = c
				break
			}
		}
	}
	if resolvedPath == "" {
		if stashPath != "" { resolvedPath = stashPath } else { resolvedPath = filepath.Join("bin", "stash", "stash-win.exe") }
	}
	return &StashProber{stashPath: resolvedPath, targetURL: targetURL, emitter: emitter}
}

func (p *StashProber) IsConnected() bool { return p.connected }

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
						p.emitter("toast:notify", map[string]string{"type": "success", "message": "🟢 Stash 接続完了！"})
					}
					fmt.Println("[StashProber] 🟢 Stash 接続確認完了 (Ready)")
				}
			} else {
				if p.connected {
					p.connected = false
					disconnectedSince = time.Now()
					if p.emitter != nil { p.emitter("stash:ready", false) }
					fmt.Println("[StashProber] ⚠️ Stash 切断検知 (Reconnecting...)")
				}
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
