// app/app_rpc_thunder_escalate.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const defaultThunderPath = `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`

// AddTaskViaThunderCOM は Windows COM (ThunderAgent) を呼び出して迅雷のキューへ直接投入します
func AddTaskViaThunderCOM(downloadURL string, fileName string, destDir string) bool {
	if downloadURL == "" {
		return false
	}
	if fileName == "" {
		fileName = filepath.Base(downloadURL)
	}
	if destDir != "" {
		if abs, err := filepath.Abs(destDir); err == nil {
			destDir = abs
		}
	}
	cleanDest := strings.ReplaceAll(destDir, "'", "''")
	// COM AddTask(URL, FileName, Path, Comments, ReferUrl, StartMode, OnlyFromOrigin, OriginThreadCount)
	psCmd := fmt.Sprintf(
		`$ids=@('ThunderAgent.Agent64.1','ThunderAgent.Agent64','ThunderAgent.Agent.1','ThunderAgent.Agent'); foreach($id in $ids){ try{ $a=New-Object -ComObject $id; if($a){ $a.AddTask('%s','%s','%s','','',1,0,5); $a.CommitTasks(); exit 0 } }catch{} }; exit 1`,
		downloadURL, fileName, cleanDest,
	)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", psCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	return cmd.Run() == nil
}

// LaunchThunder は Thunder.exe 単体を起動します
func (a *App) LaunchThunder() (bool, error) {
	if _, err := os.Stat(defaultThunderPath); err != nil {
		return false, err
	}
	cmd := exec.Command(defaultThunderPath)
	return cmd.Start() == nil, nil
}

// EscalateToThunder は指定URLをThunderへ投入し、DBステータスをESCALATEDへ更新します
func (a *App) EscalateToThunder(mediaID string, downloadURL string) (bool, error) {
	if downloadURL == "" && mediaID != "" && a.Repo != nil {
		if m, err := a.Repo.GetMediaByID(mediaID); err == nil && m != nil {
			downloadURL = m.DownloadURL
		}
	}
	if downloadURL == "" {
		return false, fmt.Errorf("download url is required")
	}

	destDir := a.getMediaDownloadDir()
	started := AddTaskViaThunderCOM(downloadURL, mediaID, destDir)

	if !started {
		thURL := EncodeThunderURL(downloadURL)
		cmdScheme := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", thURL)
		cmdScheme.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
		started = (cmdScheme.Start() == nil)
	}

	if started && mediaID != "" && a.Repo != nil {
		_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", "Thunder P2SP エスカレーション投入")
	}
	return started, nil
}

func (a *App) getMediaDownloadDir() string {
	if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.LocalMediaDir != "" {
		return cfg.Storage.LocalMediaDir
	}
	return filepath.Join("blobs")
}

// GiveUpRetainedMedia はユーザーが明示的に諦めたタスクを DEAD_404 状態へ確定します
func (a *App) GiveUpRetainedMedia(mediaID string) (bool, error) {
	if mediaID == "" || a.Repo == nil {
		return false, fmt.Errorf("mediaID is required")
	}
	err := a.Repo.UpdateMediaMetadata(mediaID, "DEAD_404", "", "", "ユーザーによる探索諦め (GIVE_UP)")
	return err == nil, err
}
