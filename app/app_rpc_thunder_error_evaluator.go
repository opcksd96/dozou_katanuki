// app/app_rpc_thunder_error_evaluator.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"regexp"
	"strings"
	"time"
)

type ThunderErrorDecision string

const (
	DecisionNone     ThunderErrorDecision = "NONE"
	DecisionRetire   ThunderErrorDecision = "RETIRE"
	DecisionHold     ThunderErrorDecision = "HOLD"
	DecisionCooldown ThunderErrorDecision = "COOLDOWN"
)

type ThunderEvaluation struct {
	Decision    ThunderErrorDecision
	SummarySize string
	HasSummary  bool
	Reason      string
}

var sizeRegex = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(B|KB|MB|GB)`)

// EvaluateThunderTaskError はCDPタスクの生テキストからエラー文言とダウンロードサマリサイズ(>1B)を評価します
func EvaluateThunderTaskError(rawText string) ThunderEvaluation {
	summary, hasSummary := extractSummarySize(rawText)
	res := ThunderEvaluation{Decision: DecisionNone, SummarySize: summary, HasSummary: hasSummary}

	// 1. メタデータ（サマリサイズ > 1B）取得済みの場合はリソース存在確認済み ➔ HOLD維持
	if hasSummary {
		res.Decision = DecisionHold
		res.Reason = "メタデータ(サマリサイズ>1B)取得済みのためタスク維持 (ESCALATED)"
		return res
	}

	// 2. 429 等のネットワーク制限・一時異常 (10分クールダウン対象)
	if strings.Contains(rawText, "429") || strings.Contains(rawText, "Too Many Requests") ||
		strings.Contains(rawText, "网络异常") || strings.Contains(rawText, "连接超时") {
		res.Decision = DecisionCooldown
		res.Reason = "ネットワーク制限または429検知 (10分クールダウン)"
		return res
	}

	// 3. サマリ未取得(0B)かつリソース不存在・各種失敗エラーの場合 ➔ RETIRE（取り下げ）
	retirePatterns := []string{
		"原始资源不存在", "未找到候选资源", "无法继续下载", "暂无任何有效资源",
		"请更换下载链接", "下载失败", "任务出错", "资源不足", "下载遇到错误",
		"文件不存在", "链接失效", "404 Not Found", "403 Forbidden",
	}
	for _, p := range retirePatterns {
		if strings.Contains(rawText, p) {
			res.Decision = DecisionRetire
			res.Reason = p + " かつサマリ0B (RETIRED・取り下げ)"
			return res
		}
	}

	// 4. アクティブ接続中・ダウンロード中
	if strings.Contains(rawText, "正在连接") || strings.Contains(rawText, "正在下载") ||
		strings.Contains(rawText, "资源连接中") || strings.Contains(rawText, "排队") {
		return res
	}

	return res
}

func extractSummarySize(text string) (string, bool) {
	matches := sizeRegex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 { return "0B", false }
	for _, m := range matches {
		valStr, unit := m[1], strings.ToUpper(m[2])
		if valStr == "0" || valStr == "0.0" || valStr == "0.00" { continue }
		return valStr + unit, true
	}
	return "0B", false
}

// IsThunderCooldownActive は前回実行から10分経過しているかを判定します
func IsThunderCooldownActive(lastAttempt *time.Time) bool {
	if lastAttempt == nil { return false }
	return time.Since(*lastAttempt) < 10*time.Minute
}
