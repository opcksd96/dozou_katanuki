// app/app_rpc_thunder_cdp_api.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"strings"
	"time"
)

type ThunderCDPTaskItem struct {
	FileName string `json:"file_name"`
	Status   string `json:"status"`
	RawText  string `json:"raw_text"`
}

type ThunderCDPStatus struct {
	IsConnected      bool                 `json:"is_connected"`
	Port             int                  `json:"port"`
	IntervalMs       int                  `json:"interval_ms"`
	ActiveWSUrl      string               `json:"active_ws_url"`
	LastPolledAt     string               `json:"last_polled_at"`
	ActiveTab        string               `json:"active_tab"`
	IsDownloadingTab bool                 `json:"is_downloading_tab"`
	CapturedTasks    []ThunderCDPTaskItem `json:"captured_tasks"`
}

// GetThunderCDPStatus は 現在の迅雷 CDP 接続状態・アクティブタブ・タスク一覧を取得します
func (a *App) GetThunderCDPStatus() ThunderCDPStatus {
	st := ThunderCDPStatus{Port: 9222, IntervalMs: 200, ActiveTab: "下载中", IsDownloadingTab: true}
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil || wsURL == "" { st.IsConnected = false; return st }
	st.IsConnected, st.ActiveWSUrl, st.LastPolledAt = true, wsURL, time.Now().Format("15:04:05")

	tabRes, _ := EvaluateCDPExpression(wsURL, `(() => {
		const items = Array.from(document.querySelectorAll('.xly-nav__item, .xly-side-menu__item, [class*="nav__item"]'));
		for (let it of items) {
			if (it.classList.contains('is-active') || it.classList.contains('active') || it.classList.contains('selected')) {
				const txt = (it.querySelector('.xly-nav__tab, .xly-nav__content') || it).innerText.split('\n')[0].trim();
				return { tab: txt, isDl: txt === '下载中' };
			}
		}
		return { tab: '下载中', isDl: true };
	})()`, 500*time.Millisecond)

	var tabData struct{ Result struct{ Result struct{ Value struct{ Tab string `json:"tab"`; IsDl bool `json:"isDl"` } `json:"value"` } `json:"result"` } `json:"result"` }
	if json.Unmarshal([]byte(tabRes), &tabData) == nil && tabData.Result.Result.Value.Tab != "" {
		st.ActiveTab = tabData.Result.Result.Value.Tab
		st.IsDownloadingTab = tabData.Result.Result.Value.IsDl
	}

	resJSON, err := EvaluateCDPExpression(wsURL, ThunderExtractTaskScript, 1000*time.Millisecond)
	if err == nil {
		var evalRes CDPEvalResult
		if err := json.Unmarshal([]byte(resJSON), &evalRes); err == nil {
			st.CapturedTasks = parseCDPTargetItems(evalRes.Result.Result.Value)
		}
	}
	return st
}

// SwitchThunderTabViaCDP は迅雷のタブを指定したもの（"下载中" など）に切り替えます
func (a *App) SwitchThunderTabViaCDP(targetTab string) bool {
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil || wsURL == "" { return false }
	script := `(() => {
		const tabs = Array.from(document.querySelectorAll('.xly-nav__tab, span, a, .xly-nav__item'));
		const target = tabs.find(el => el.innerText && el.innerText.trim() === '` + targetTab + `');
		if (target) { target.click(); return true; }
		return false;
	})()`
	EvaluateCDPExpression(wsURL, script, 1000*time.Millisecond)
	return true
}

func parseCDPTargetItems(rawBlocks []string) []ThunderCDPTaskItem {
	var items []ThunderCDPTaskItem
	seen := make(map[string]bool)
	for _, block := range rawBlocks {
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if (strings.Contains(line, ".jpg") || strings.Contains(line, ".mp4") || strings.Contains(line, ".png")) && !seen[line] {
				seen[line] = true
				status := "排队等待 / 探索中"
				if strings.Contains(block, "无法继续下载") || strings.Contains(block, "暂无任何有效资源") {
					status = "リソース枯渇 (RETAINED対象)"
				} else if strings.Contains(block, "连接资源") {
					status = "ピア探索中 (ESCALATED)"
				} else if strings.Contains(block, "KB") || strings.Contains(block, "MB") {
					status = "ダウンロード完了"
				}
				items = append(items, ThunderCDPTaskItem{FileName: line, Status: status, RawText: line})
			}
		}
	}
	return items
}
