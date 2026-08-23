// middleware/stash_prober_process.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func (p *StashProber) isStashProcessAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "tasklist", "/FI", "IMAGENAME eq stash-win.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.Output()
	if err != nil { return false }
	return strings.Contains(strings.ToLower(string(out)), "stash-win.exe")
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
		if p.isStashProcessAlive() { return }
		absPath, err := filepath.Abs(p.stashPath)
		if err != nil { return }
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
