// app/app_thunder_cdp_poller.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"time"
)

type CDPEvalResult struct {
	Result struct {
		Result struct {
			Value []string `json:"value"`
		} `json:"result"`
	} `json:"result"`
}

// StartThunderCDPAdaptivePoller は 2000ms 間隔の軽量サイレントバックグラウンド同期を実行します
func (a *App) StartThunderCDPAdaptivePoller() {
	go func() {
		wsURL := ""
		interval := 2000 * time.Millisecond

		for {
			time.Sleep(interval)

			if wsURL == "" {
				u, err := FetchThunderMainRendererWSUrl(9222)
				if err != nil {
					interval = 3000 * time.Millisecond
					continue
				}
				wsURL = u
				a.AppendPipelineLog("THUNDER", "INFO", "迅雷 CDP WebSocket 接続を確立しました")
			}

			if !a.isThunderOrchestratorRunning() {
				interval = 3000 * time.Millisecond
				continue
			}

			resJSON, err := EvaluateCDPExpression(wsURL, ThunderExtractTaskScript, 1500*time.Millisecond)
			if err != nil {
				wsURL = ""
				interval = 3000 * time.Millisecond
				a.AppendPipelineLog("THUNDER", "WARN", fmt.Sprintf("迅雷 CDP 通信途絶: %v", err))
				continue
			}

			var evalRes CDPEvalResult
			if err := json.Unmarshal([]byte(resJSON), &evalRes); err != nil || len(evalRes.Result.Result.Value) == 0 {
				interval = 2000 * time.Millisecond
				continue
			}

			interval = 2000 * time.Millisecond
			reconciled, _ := a.ReconcileThunderTasksWithDB()
			if reconciled > 0 {
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("迅雷タスクとDBの突合同期完了 (%d 件更新)", reconciled))
			}
		}
	}()
}
