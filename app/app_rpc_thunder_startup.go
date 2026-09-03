// app/app_rpc_thunder_startup.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

// CheckThunderCDPOnStartup は dozou起動時に迅雷のCDP解放状態をテストしトースト通知・オーケストレーター起動を行います
func (a *App) CheckThunderCDPOnStartup() {
	time.Sleep(1200 * time.Millisecond) // UI描画完了を少し待機

	// ③ DB接続確認
	if a.Repo == nil || a.Repo.DB() == nil {
		a.emitToast("error", "❌ データベース未接続のため迅雷連携を待機します")
		return
	}

	// ① CDP解放状態テスト & トースト通知
	if isThunderProcessRunning() && isThunderCDPListening() {
		a.emitToast("success", "⚡ 迅雷 CDP (:9222) 接続確認完了 (Ready)")
		a.AppendPipelineLog("THUNDER", "SUCCESS", "⚡ 迅雷 CDP (:9222) 正常接続を確認")
	} else if isThunderProcessRunning() && !isThunderCDPListening() {
		a.emitToast("warning", "⚠️ 迅雷がCDP無効で起動中のため再起動します...")
		go func() {
			if ok, err := a.LaunchThunder(); ok {
				a.emitToast("success", "⚡ 迅雷 CDP (:9222) 有効化で再起動完了")
			} else {
				a.emitToast("error", fmt.Sprintf("❌ 迅雷再起動失敗: %v", err))
			}
		}()
	} else {
		a.AppendPipelineLog("THUNDER", "INFO", "ℹ️ 迅雷は現在未起動です (必要時に自律起動)")
	}

	// ② オーケストレータ起動判定 (ESCALATEDメディアが存在すれば自動点火)
	var escalatedCount int64
	_ = a.Repo.DB().Model(&models.Media{}).Where("download_status = 'ESCALATED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&escalatedCount).Error

	if escalatedCount > 0 && !a.isThunderOrchestratorRunning() {
		a.emitToast("info", fmt.Sprintf("⚡ 迅雷オーケストレーターを自動起動 (残 %d 件)", escalatedCount))
		_, _ = a.StartThunderOrchestrator(3, 4)
	}
}

func (a *App) emitToast(toastType, message string) {
	if a.Ctx != nil {
		a.EmitEvent("toast:notify", map[string]interface{}{
			"type":    toastType,
			"message": message,
		})
	}
}
