// app/app_rpc_thunder_escalate.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ThunderCOMTask は迅雷COM投入用の単一タスク定義です
type ThunderCOMTask struct {
	URL      string
	FileName string
	DestDir  string
}

// AddTaskViaThunderCOM は 単一タスクをWindows COM (ThunderAgent) でサイレント投入します
func AddTaskViaThunderCOM(downloadURL string, fileName string, destDir string) bool {
	return AddBatchTasksViaThunderCOM([]ThunderCOMTask{{URL: downloadURL, FileName: fileName, DestDir: destDir}})
}

// AddBatchTasksViaThunderCOM は 複数タスクを1プロセスで迅雷COMへ一括投入します
func AddBatchTasksViaThunderCOM(tasks []ThunderCOMTask) bool {
	if len(tasks) == 0 {
		return false
	}
	var b strings.Builder
	b.WriteString(`$ids=@('ThunderAgent.Agent64.1','ThunderAgent.Agent64','ThunderAgent.Agent.1','ThunderAgent.Agent'); foreach($id in $ids){ try{ $a=New-Object -ComObject $id; if($a){ `)
	for _, t := range tasks {
		if t.URL == "" {
			continue
		}
		fn := t.FileName
		if fn == "" {
			fn = filepath.Base(t.URL)
		}
		dest := t.DestDir
		if dest != "" {
			if abs, err := filepath.Abs(dest); err == nil {
				dest = abs
			}
		}
		cleanURL := strings.ReplaceAll(t.URL, "'", "''")
		cleanFN := strings.ReplaceAll(fn, "'", "''")
		cleanDest := strings.ReplaceAll(dest, "'", "''")
		b.WriteString(fmt.Sprintf(`$a.AddTask('%s','%s','%s','','',1,0,5); `, cleanURL, cleanFN, cleanDest))
	}
	b.WriteString(`$a.CommitTasks(); exit 0 } }catch{} }; exit 1`)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", b.String())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd.Run() == nil
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

	check := a.CheckThunderEscalationEligibility(mediaID, downloadURL)
	if !check.ShouldEscalate {
		if check.ExistingStatus == "ALREADY_COMPLETED" && mediaID != "" && a.Repo != nil {
			_ = a.Repo.UpdateMediaMetadata(mediaID, "COMPLETED", "", "", "実体検知による完了同期")
		} else if (check.ExistingStatus == "ALREADY_IN_THUNDER" || check.ExistingStatus == "DOWNLOADING_XLTD") && mediaID != "" && a.Repo != nil {
			_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", check.Reason)
		}
		a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ 迅雷エスカレーション判定: スキップ (%s)", check.Reason))
		return true, nil
	}

	if mediaID != "" && a.Repo != nil {
		_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", "迅雷オーケストレーター投入待機 (3スロット制御)")
		a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ 迅雷キュー登録: %s (自律オーケストレーターへ委託)", mediaID))
		if !a.isThunderOrchestratorRunning() {
			_, _ = a.StartThunderOrchestrator(3, 4)
		}
	}
	return true, nil
}

func (a *App) getMediaDownloadDir() string {
	if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.LocalMediaDir != "" {
		return cfg.Storage.LocalMediaDir
	}
	return filepath.Join("blobs")
}
