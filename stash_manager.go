// stash_manager.go (100行以下)
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type StashManager struct {
	ctx     context.Context
	cmd     *exec.Cmd
	Running bool
}

func NewStashManager() *StashManager { return &StashManager{} }

func (sm *StashManager) WaitForReady(ctx context.Context, timeout time.Duration) error {
	client, deadline := &http.Client{Timeout: 500 * time.Millisecond}, time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done(): return ctx.Err()
		default:
		}
		if resp, err := client.Get("http://127.0.0.1:9999/"); err == nil {
			_ = resp.Body.Close()
			if sm.ctx != nil && sm.ctx.Err() == nil { runtime.EventsEmit(sm.ctx, "stash:ready", true) }
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Errorf("stash server readiness check timed out after %v", timeout)
}

func (sm *StashManager) PurgeZombies() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "taskkill", "/F", "/IM", "stash-win.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = cmd.Run()
}

func (sm *StashManager) Start(ctx context.Context, stashPath string) error {
	sm.ctx = ctx
	sm.PurgeZombies()
	absPath, err := filepath.Abs(stashPath)
	if err != nil { return fmt.Errorf("Stashパスの絶対パス解決失敗: %w", err) }
	cmd := exec.Command(absPath)
	cmd.Dir = filepath.Dir(absPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000 | 0x00000200}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil { return fmt.Errorf("stash-win.exe 起動失敗: %w", err) }
	sm.cmd, sm.Running = cmd, true
	fmt.Printf("[StashManager] Stash ヘッドレス起動開始 (PID: %d)\n", cmd.Process.Pid)

	if stdout != nil { go sm.scanPipe(stdout, "STDOUT") }
	if stderr != nil { go sm.scanPipe(stderr, "STDERR") }
	go func() {
		if err := sm.WaitForReady(ctx, 15*time.Second); err == nil {
			fmt.Println("[StashManager] Stash HTTPサーバー疎通完了 (Ready)")
		} else {
			fmt.Printf("[StashManager] Stash 疎通確認タイムアウト: %v\n", err)
		}
	}()
	return nil
}

func (sm *StashManager) scanPipe(pipe io.Reader, pipeType string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		if sm.ctx != nil && sm.ctx.Err() == nil {
			runtime.EventsEmit(sm.ctx, "stash:log", fmt.Sprintf("[%s] %s", pipeType, scanner.Text()))
		} else { break }
	}
}

func (sm *StashManager) Stop() {
	if sm.cmd != nil && sm.cmd.Process != nil {
		pid := sm.cmd.Process.Pid
		_ = sm.cmd.Process.Kill()
		sm.Running = false
		go func(p int) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_ = exec.CommandContext(ctx, "taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", p)).Run()
		}(pid)
	}
	go sm.PurgeZombies()
}
