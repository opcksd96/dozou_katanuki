// middleware/stash_prober.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type StashProber struct {
	emitter        EventEmitter
	stashPath      string
	targetURL      string
	running        bool
	connected      bool
	lastNotifiedAt time.Time
	cancelFunc     context.CancelFunc
}

func NewStashProber(stashPath string, targetURL string, emitter EventEmitter) *StashProber {
	if targetURL == "" { targetURL = "http://127.0.0.1:9999/" }
	candidates := []string{stashPath, filepath.Join("bin", "stash", "stash-win.exe"), filepath.Join("bin", "stash-win.exe"), "stash-win.exe"}
	resolvedPath := ""
	for _, c := range candidates {
		if c != "" {
			if _, err := os.Stat(c); err == nil { resolvedPath = c; break }
		}
	}
	if resolvedPath == "" {
		if stashPath != "" { resolvedPath = stashPath } else { resolvedPath = filepath.Join("bin", "stash", "stash-win.exe") }
	}
	return &StashProber{stashPath: resolvedPath, targetURL: targetURL, emitter: emitter}
}

func (p *StashProber) IsConnected() bool { return p.connected }

func (p *StashProber) Start(ctx context.Context) {
	_, cancel := context.WithCancel(ctx)
	p.cancelFunc = cancel
	p.running = true
	// 実際のプロービングは外部(PS1ビーコン)からの UpdateStatus 呼び出しに委譲
	fmt.Println("[StashProber] Start: ビーコン受信待機モードで起動しました。")
}

func (p *StashProber) UpdateStatus(isConnected bool) {
	if p.connected == isConnected {
		return // No change
	}

	p.connected = isConnected
	if isConnected {
		if p.emitter != nil {
			p.emitter("stash:ready", true)
			if time.Since(p.lastNotifiedAt) > 30*time.Second {
				p.lastNotifiedAt = time.Now()
				p.emitter("toast:notify", map[string]string{"type": "success", "message": "🟢 Stash 接続完了 (Beacon)！"})
			}
		}
		GetGlobalJournal().Record("stash", "INFO", "stash_connected", "Stash connection established via Beacon", nil)
		fmt.Println("[StashProber] 🟢 Stash 接続確認完了 (Beacon Ready)")
	} else {
		if p.emitter != nil {
			p.emitter("stash:ready", false)
		}
		GetGlobalJournal().Record("stash", "WARN", "stash_disconnected", "Stash connection lost via Beacon", nil)
		fmt.Println("[StashProber] ⚠️ Stash 未接続 / 停止 (Beacon Waiting...)")
	}
}
