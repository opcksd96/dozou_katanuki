// app/app_thunder_cdp_poller.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
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
			}

			resJSON, err := EvaluateCDPExpression(wsURL, ThunderExtractTaskScript, 1500*time.Millisecond)
			if err != nil {
				wsURL = ""
				interval = 3000 * time.Millisecond
				continue
			}

			var evalRes CDPEvalResult
			if err := json.Unmarshal([]byte(resJSON), &evalRes); err != nil || len(evalRes.Result.Result.Value) == 0 {
				interval = 2000 * time.Millisecond
				continue
			}

			interval = 2000 * time.Millisecond
			_, _ = a.ReconcileThunderTasksWithDB()
		}
	}()
}
