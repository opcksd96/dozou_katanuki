// app/app_rpc_thunder_escalate_guard.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EscalationCheckResult は迅雷エスカレーション判定の結果を保持します
type EscalationCheckResult struct {
	ShouldEscalate bool
	Reason         string
	ExistingStatus string
}

// CheckThunderEscalationEligibility はローカルファイル・xltd一時ファイル・迅雷CDPリストを総合判定します
func (a *App) CheckThunderEscalationEligibility(mediaID string, downloadURL string) EscalationCheckResult {
	if mediaID == "" && downloadURL == "" {
		return EscalationCheckResult{ShouldEscalate: false, Reason: "URLまたはMediaIDが指定されていません"}
	}

	cleanID := strings.TrimSuffix(mediaID, filepath.Ext(mediaID))
	destDir := a.getMediaDownloadDir()

	// 1. ローカル実体ファイル判定
	for _, ext := range []string{".jpg", ".mp4", ".png"} {
		targetFile := filepath.Join(destDir, cleanID+ext)
		if fi, err := os.Stat(targetFile); err == nil && fi.Size() > 0 {
			return EscalationCheckResult{ShouldEscalate: false, Reason: "ローカル実体ファイルが既に存在します", ExistingStatus: "ALREADY_COMPLETED"}
		}
	}

	// 2. 迅雷一時ファイル (*.xltd) 判定
	matches, _ := filepath.Glob(filepath.Join(destDir, cleanID+"*.*.xltd"))
	if len(matches) > 0 {
		return EscalationCheckResult{ShouldEscalate: false, Reason: "迅雷ダウンロード中の一時ファイル(*.xltd)が存在します", ExistingStatus: "DOWNLOADING_XLTD"}
	}

	// 3. 迅雷 CDP ライブタスクリスト照合
	cdpStatus := a.GetThunderCDPStatus()
	if cdpStatus.IsConnected && len(cdpStatus.CapturedTasks) > 0 {
		for _, task := range cdpStatus.CapturedTasks {
			if cleanID != "" && strings.Contains(task.FileName, cleanID) {
				if strings.Contains(task.Status, "リソース枯渇") {
					return EscalationCheckResult{ShouldEscalate: false, Reason: fmt.Sprintf("迅雷内でリソース枯渇判定されています (%s)", task.FileName), ExistingStatus: "DEPLETED_RETAINED"}
				}
				if strings.Contains(task.Status, "完了") {
					return EscalationCheckResult{ShouldEscalate: false, Reason: fmt.Sprintf("迅雷内で既にダウンロード完了しています (%s)", task.FileName), ExistingStatus: "ALREADY_COMPLETED"}
				}
				return EscalationCheckResult{ShouldEscalate: false, Reason: fmt.Sprintf("迅雷のダウンロードリストに既に登録されています (%s)", task.FileName), ExistingStatus: "ALREADY_IN_THUNDER"}
			}
		}
	}

	return EscalationCheckResult{ShouldEscalate: true, Reason: "新規エスカレーション投入可能", ExistingStatus: "READY"}
}
