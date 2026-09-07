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

// LaunchThunder は config 設定に従い Thunder.exe を起動（またはCDP再起動）します
func (a *App) LaunchThunder() (bool, error) {
	thunderPath, launchWithCDP, cdpPort := defaultThunderPath, true, 9222
	if cfg, err := a.GetConfig(); err == nil && cfg != nil {
		if cfg.Thunder.Path != "" { thunderPath = cfg.Thunder.Path }
		launchWithCDP = cfg.Thunder.LaunchWithCDP
		if cfg.Thunder.CDPPort > 0 { cdpPort = cfg.Thunder.CDPPort }
	}

	if _, err := os.Stat(thunderPath); err != nil {
		return false, fmt.Errorf("迅雷バイナリが見つかりません: %s", thunderPath)
	}

	// CDP付き起動が要求されているのに、CDP未疎通のThunderプロセスが居座っている場合は再起動
	if launchWithCDP && isThunderProcessRunning() && !isThunderCDPListening() {
		a.AppendPipelineLog("THUNDER", "WARN", "⚠️ 迅雷CDP未開通プロセスを検知。CDP有効化のため再起動します")
		KillThunderProcess()
	}

	var args []string
	args = append(args, "/c", "start", "", thunderPath)
	if launchWithCDP {
		args = append(args, fmt.Sprintf("--remote-debugging-port=%d", cdpPort))
		a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ 迅雷 (CDP port %d) 起動要求を送信", cdpPort))
	} else {
		a.AppendPipelineLog("THUNDER", "INFO", "⚡ 迅雷 (通常モード) 起動要求を送信")
	}

	cmd := exec.Command("cmd", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if err := cmd.Start(); err != nil {
		a.AppendPipelineLog("THUNDER", "ERROR", fmt.Sprintf("❌ 迅雷起動失敗: %v", err))
		return false, err
	}
	return true, nil
}

// EnsureThunderCDP は、CDPが未開通の場合にキック処理を呼び出します
func (a *App) EnsureThunderCDP() (bool, error) {
	if isThunderCDPListening() { return true, nil }
	return a.LaunchThunder()
}
