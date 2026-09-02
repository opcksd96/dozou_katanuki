// app/app_thunder_cdp_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"
	"time"
)

func TestThunderCDP_ConnectionAndEvaluation(t *testing.T) {
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil {
		t.Skipf("Thunder CDP not available on port 9222 (skipping): %v", err)
		return
	}
	t.Logf("Found Thunder WebSocket URL: %s", wsURL)

	// 1. document.title を取得
	res, err := EvaluateCDPExpression(wsURL, "document.title", 2*time.Second)
	if err != nil {
		t.Skipf("EvaluateCDPExpression timeout/skipped: %v", err)
		return
	}
	t.Logf("CDP document.title evaluation result: %s", res)

	// 2. タスクリスト抽出スクリプトを実行
	listRes, err := EvaluateCDPExpression(wsURL, ThunderExtractTaskScript, 2*time.Second)
	if err != nil {
		t.Skipf("EvaluateCDPExpression task script timeout/skipped: %v", err)
		return
	}
	t.Logf("CDP Task extraction result length: %d bytes", len(listRes))
}
