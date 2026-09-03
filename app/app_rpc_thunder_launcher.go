// app/app_rpc_thunder_launcher.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const defaultThunderPath = `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`

// isThunderProcessRunning は Windows OS 上で Thunder.exe が稼働しているか判定します
func isThunderProcessRunning() bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq Thunder.exe", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.Output()
	return err == nil && strings.Contains(strings.ToLower(string(out)), "thunder.exe")
}

// isThunderCDPListening は ポート 9222 で迅雷 CDP が疎通可能か判定します
func isThunderCDPListening() bool {
	u, err := FetchThunderMainRendererWSUrl(9222)
	return err == nil && u != ""
}

// KillThunderProcess は 稼働中の Thunder.exe プロセスを強制終了します
func KillThunderProcess() bool {
	cmd := exec.Command("taskkill", "/F", "/IM", "Thunder.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = cmd.Run()
	time.Sleep(500 * time.Millisecond)
	return !isThunderProcessRunning()
}

// LaunchThunder は Thunder.exe を CDP デバッグポート 9222 を有効にして非同期で起動（キック）します。
// 死活監視やプロセス管理は行わず、単にシェルに展開するだけの責務を持ちます。
func (a *App) LaunchThunder() (bool, error) {
	if _, err := os.Stat(defaultThunderPath); err != nil {
		return false, fmt.Errorf("迅雷バイナリが見つかりません: %s", defaultThunderPath)
	}

	a.AppendPipelineLog("THUNDER", "INFO", "⚡ 迅雷プロセスの起動要求を送信します (Fire and Forget)")

	// cmd.exe /c start を使って、Goプロセスとは完全に切り離された非同期プロセスとしてキックする
	cmd := exec.Command("cmd", "/c", "start", "", defaultThunderPath, "--remote-debugging-port=9222")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	
	if err := cmd.Start(); err != nil {
		a.AppendPipelineLog("THUNDER", "ERROR", fmt.Sprintf("❌ 迅雷のキックに失敗しました: %v", err))
		return false, err
	}

	// 起動の成否やCDPの開通確認は BeaconService (ポート9222のポーリング) に任せる
	return true, nil
}

// EnsureThunderCDP は、CDPが未開通の場合にキック処理を呼び出します
func (a *App) EnsureThunderCDP() (bool, error) {
	if isThunderCDPListening() {
		return true, nil
	}
	return a.LaunchThunder()
}
