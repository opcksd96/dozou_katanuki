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

// LaunchThunder は Thunder.exe を CDP デバッグポート 9222 を有効にして起動します
// 既にCDP無効で起動中の場合は、プロセスを再起動してCDPを有効化します
func (a *App) LaunchThunder() (bool, error) {
	if _, err := os.Stat(defaultThunderPath); err != nil {
		return false, fmt.Errorf("迅雷バイナリが見つかりません: %s", defaultThunderPath)
	}

	// 1. 既に CDP が有効で稼働中であれば何もしない
	if isThunderProcessRunning() && isThunderCDPListening() {
		a.AppendPipelineLog("THUNDER", "INFO", "⚡ 迅雷は既に CDP(:9222) 有効で稼働中です")
		return true, nil
	}

	// 2. 迅雷が起動中だが CDP が無効な場合は強制終了して再起動
	if isThunderProcessRunning() && !isThunderCDPListening() {
		a.AppendPipelineLog("THUNDER", "INFO", "⚠️ 迅雷がCDP無効で起動中のため、CDP(:9222)有効化のため再起動します...")
		KillThunderProcess()
	}

	// 3. CDP ポート 9222 を有効にして起動
	cmd := exec.Command(defaultThunderPath, "--remote-debugging-port=9222")
	if err := cmd.Start(); err != nil {
		a.AppendPipelineLog("THUNDER", "ERROR", fmt.Sprintf("❌ 迅雷の起動に失敗しました: %v", err))
		return false, err
	}

	// 4. CDP リッスン開始を最大3秒間待機
	for i := 0; i < 6; i++ {
		time.Sleep(500 * time.Millisecond)
		if isThunderCDPListening() {
			a.AppendPipelineLog("THUNDER", "SUCCESS", "⚡ 迅雷を CDP(:9222) 有効化で起動・開通しました")
			return true, nil
		}
	}

	a.AppendPipelineLog("THUNDER", "INFO", "⚡ 迅雷プロセスを起動しました (CDP 待機中)")
	return true, nil
}

// EnsureThunderCDP は CDP 疎通を確認し、未開通なら自動でCDP有効起動します
func (a *App) EnsureThunderCDP() (bool, error) {
	if isThunderProcessRunning() && isThunderCDPListening() {
		return true, nil
	}
	return a.LaunchThunder()
}
