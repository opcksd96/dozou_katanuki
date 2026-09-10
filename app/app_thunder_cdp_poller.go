// app/app_thunder_cdp_poller.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type CDPEvalResult struct {
	Result struct {
		Result struct {
			Value []string `json:"value"`
		} `json:"result"`
	} `json:"result"`
}

var (
	cdpPollerMu     sync.Mutex
	isCDPPollerBusy bool
)

// StartThunderCDPAdaptivePoller は 軽量サイレントバックグラウンド同期を実行します (単一常駐)
func (a *App) StartThunderCDPAdaptivePoller() {
	cdpPollerMu.Lock()
	if isCDPPollerBusy {
		cdpPollerMu.Unlock()
		return
	}
	isCDPPollerBusy = true
	cdpPollerMu.Unlock()

	go func() {
		defer func() {
			cdpPollerMu.Lock()
			isCDPPollerBusy = false
			cdpPollerMu.Unlock()
		}()

		wsURL := ""
		interval := 3000 * time.Millisecond

		for {
			time.Sleep(interval)

			if !a.IsPipelineAutoEngineRunning() || !a.isThunderOrchestratorRunning() {
				interval = 3000 * time.Millisecond
				continue
			}

			if wsURL == "" {
				u, err := FetchThunderMainRendererWSUrl(9222)
				if err != nil {
					interval = 4000 * time.Millisecond
					continue
				}
				wsURL = u
				a.AppendPipelineLog("THUNDER", "INFO", "迅雷 CDP WebSocket 接続を確立しました")
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
