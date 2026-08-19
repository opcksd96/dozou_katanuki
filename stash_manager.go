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

func NewStashManager() *StashManager {
	return &StashManager{}
}

// WaitForReady は Stash サーバー (http://127.0.0.1:9999) が応答するまで待機します
func (sm *StashManager) WaitForReady(ctx context.Context, timeout time.Duration) error {
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := client.Get("http://127.0.0.1:9999/")
		if err == nil {
			_ = resp.Body.Close()
			if sm.ctx != nil {
				runtime.EventsEmit(sm.ctx, "stash:ready", true)
			}
			return nil
		}

		time.Sleep(150 * time.Millisecond)
	}

	return fmt.Errorf("stash server readiness check timed out after %v", timeout)
}

// 起動前のゾンビプロセス強制クレンジング
func (sm *StashManager) PurgeZombies() {
	fmt.Println("[StashManager] ゾンビプロセスの事前スキャン＆パージ中...")
	_ = exec.Command("taskkill", "/F", "/IM", "stash-win.exe").Run()
}

// ヘッドレス起動とパイプ監視の開始
func (sm *StashManager) Start(ctx context.Context, stashPath string) error {
	sm.ctx = ctx
	sm.PurgeZombies()

	// 🌟 相対パスを完全な絶対パスに変換（これでCWD変更によるパス迷子を100%防止）
	absPath, err := filepath.Abs(stashPath)
	if err != nil {
		return fmt.Errorf("Stashパスの絶対パス解決失敗: %w", err)
	}

	cmd := exec.Command(absPath)
	// 作業ディレクトリをバイナリが存在するディレクトリ（./bin）に明示的に設定
	cmd.Dir = filepath.Dir(absPath)

	// CREATE_NO_WINDOW (0x08000000) フラグで黒いコンソール画面を完全非表示化
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("StdoutPipe 取得失敗: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("StderrPipe 取得失敗: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("stash-win.exe 起動失敗: %w", err)
	}

	sm.cmd = cmd
	sm.Running = true
	fmt.Printf("[StashManager] Stash をヘッドレス起動しました (PID: %d, Path: %s)\n", cmd.Process.Pid, absPath)

	// Stdout / Stderr の非同期ログスキャン Goroutine
	go sm.scanPipe(stdout, "STDOUT")
	go sm.scanPipe(stderr, "STDERR")

	return nil
}

// パイプを常時走査し、Wails Event でフロントエンドへストリーミング
func (sm *StashManager) scanPipe(pipe io.Reader, pipeType string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		// フロントエンドの Scraper / Console View へ一方向プッシュ
		if sm.ctx != nil {
			runtime.EventsEmit(sm.ctx, "stash:log", fmt.Sprintf("[%s] %s", pipeType, line))
		}
	}
}

// 道連れ終了（Lifeline Synchronization）
func (sm *StashManager) Stop() {
	if sm.cmd != nil && sm.cmd.Process != nil {
		fmt.Println("[StashManager] Wails 終了を検知。Stash を道連れ終了 (Kill) します...")
		_ = sm.cmd.Process.Kill()
		sm.Running = false
	}
}
