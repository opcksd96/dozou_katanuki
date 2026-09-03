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

	// 1. 429 等のネットワーク制限・一時異常 (10分クールダウン対象)
	if strings.Contains(rawText, "429") || strings.Contains(rawText, "Too Many Requests") ||
		strings.Contains(rawText, "网络异常") || strings.Contains(rawText, "连接超时") {
		res.Decision = DecisionCooldown
		res.Reason = "ネットワーク制限または429検知 (10分クールダウン)"
		return res
	}

	// 2. 「原始资源不存在，且未找到候选资源，无法继续下载」 かつ 0B (サマリ未取得)
	if strings.Contains(rawText, "原始资源不存在") || strings.Contains(rawText, "未找到候选资源") {
		if !hasSummary {
			res.Decision = DecisionRetire
			res.Reason = "原始リソース不存在かつサマリ0B (RETIRED・取り下げ)"
			return res
		}
	}

	// 3. 「暂无任何有效资源可连接，无法正常下载，请更换下载链接」 かつ >1B (サマリ取得済み)
	if strings.Contains(rawText, "暂无任何有效资源") || strings.Contains(rawText, "请更换下载链接") {
		if hasSummary {
			res.Decision = DecisionHold
			res.Reason = "有効リソース探索中だがサマリ>1B捕捉済み (タスク維持・ESCALATED)"
			return res
		}
		// サマリも取れていない場合は枯渇として取り下げ
		res.Decision = DecisionRetire
		res.Reason = "有効リソースなし・サマリ未取得 (RETIRED・取り下げ)"
		return res
	}

	return res
}

func extractSummarySize(text string) (string, bool) {
	matches := sizeRegex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 { return "0B", false }
	for _, m := range matches {
		valStr, unit := m[1], strings.ToUpper(m[2])
		sizeStr := valStr + unit
		if valStr == "0" || valStr == "0.0" || valStr == "0.00" {
			continue
		}
		return sizeStr, true // > 1B のサマリを捕捉
	}
	return "0B", false
}

// IsThunderCooldownActive は前回実行から10分経過しているかを判定します
func IsThunderCooldownActive(lastAttempt *time.Time) bool {
	if lastAttempt == nil { return false }
	return time.Since(*lastAttempt) < 10*time.Minute
}
