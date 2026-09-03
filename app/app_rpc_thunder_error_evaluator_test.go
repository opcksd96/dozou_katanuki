// app/app_rpc_thunder_error_evaluator_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"
	"time"
)

func TestEvaluateThunderTaskError(t *testing.T) {
	// 1. 原始资源不存在 かつ サマリ0B ➔ RETIRE
	text1 := "test1.jpg\n0B\n原始资源不存在，且未找到候选资源，无法继续下载"
	res1 := EvaluateThunderTaskError(text1)
	if res1.Decision != DecisionRetire {
		t.Errorf("expected DecisionRetire for 0B + 原始资源不存在, got %v", res1.Decision)
	}

	// 2. 暂无任何有效资源 かつ サマリ>1B (363.46KB) ➔ HOLD
	text2 := "test2.jpg\n363.46KB\n暂无任何有效资源可连接，无法正常下载，请更换下载链接"
	res2 := EvaluateThunderTaskError(text2)
	if res2.Decision != DecisionHold {
		t.Errorf("expected DecisionHold for >1B + 暂无任何有效资源, got %v", res2.Decision)
	}
	if !res2.HasSummary || res2.SummarySize != "363.46KB" {
		t.Errorf("expected summary size 363.46KB, got %s (hasSummary=%v)", res2.SummarySize, res2.HasSummary)
	}

	// 3. 暂无任何有效资源 かつ サマリ0B ➔ RETIRE
	text3 := "test3.jpg\n0 B\n暂无任何有效资源可连接，无法正常下载，请更换下载链接"
	res3 := EvaluateThunderTaskError(text3)
	if res3.Decision != DecisionRetire {
		t.Errorf("expected DecisionRetire for 0B + 暂无任何有效资源, got %v", res3.Decision)
	}

	// 4. 429 エラー ➔ COOLDOWN
	text4 := "test4.jpg\nHTTP 429 Too Many Requests: 网络异常"
	res4 := EvaluateThunderTaskError(text4)
	if res4.Decision != DecisionCooldown {
		t.Errorf("expected DecisionCooldown for 429, got %v", res4.Decision)
	}

	// 5. 10分クールダウン判定
	now := time.Now()
	recent := now.Add(-5 * time.Minute) // 5分前
	if !IsThunderCooldownActive(&recent) {
		t.Errorf("expected cooldown to be active for 5 minutes ago")
	}

	old := now.Add(-11 * time.Minute) // 11分前
	if IsThunderCooldownActive(&old) {
		t.Errorf("expected cooldown to be inactive for 11 minutes ago")
	}
}
