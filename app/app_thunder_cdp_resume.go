// app/app_thunder_cdp_resume.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ResumeAllTasksViaDirectCDP は Direct CDP で迅雷の「全部开始」を実行して一時停止中タスクを一括再開します
func ResumeAllTasksViaDirectCDP(port int) (bool, error) {
	wsURL, err := FetchThunderMainRendererWSUrl(port)
	if err != nil { return false, err }

	script := `(() => {
		// 1. ツールバーの「全部开始」ボタンを検索してクリック
		const buttons = Array.from(document.querySelectorAll('button, .td-button, .xly-download-tab__operate button'));
		const startAllBtn = buttons.find(b => {
			const txt = (b.innerText || b.getAttribute('title') || '').trim();
			return txt.includes('全部开始') || txt === '开始' || txt.includes('全部下载');
		});
		if (startAllBtn && !startAllBtn.classList.contains('is-disabled')) {
			startAllBtn.click();
			return { success: true, action: "click_all_start" };
		}

		// 2. Vue インスタンスから一時停止中タスクを一括開始
		const toolbar = document.querySelector('.xly-download-tab__operate');
		const p = toolbar && toolbar.__vue__ && toolbar.__vue__.$parent;
		if (p && p.handleStartAll) {
			p.handleStartAll();
			return { success: true, action: "vue_handleStartAll" };
		}
		if (p && p.taskBaseMap) {
			const map = p.taskBaseMap;
			let pausedIds = [];
			for (const id in map) {
				if (map[id] && (map[id].taskStatus === 7 || map[id].taskStatus === 8)) {
					pausedIds.push(map[id].taskId);
				}
			}
			if (pausedIds.length > 0 && p.handleStart) {
				p.downloadingSelectedIds = pausedIds;
				p.currentSelectdTaskIds = pausedIds;
				p.handleStart();
				return { success: true, action: "vue_resume_paused", count: pausedIds.length };
			}
		}
		return { success: true, action: "no_paused_tasks" };
	})()`

	resJSON, err := EvaluateCDP(wsURL, script, 3*time.Second)
	if err != nil { return false, err }

	var evalRes struct {
		Result struct {
			Result struct {
				Value struct {
					Success bool   `json:"success"`
					Action  string `json:"action"`
				} `json:"value"`
			} `json:"result"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resJSON), &evalRes); err != nil {
		return false, fmt.Errorf("failed to parse resume response: %w", err)
	}
	return evalRes.Result.Result.Value.Success, nil
}

// ResumeTaskViaDirectCDP は Direct CDP で指定ファイル名のタスクを再開します
func ResumeTaskViaDirectCDP(port int, targetFileName string) (bool, error) {
	wsURL, err := FetchThunderMainRendererWSUrl(port)
	if err != nil { return false, err }
	clean := strings.TrimSpace(targetFileName)

	script := fmt.Sprintf(`((targetFileName) => {
		const clean = targetFileName.trim();
		const toolbar = document.querySelector('.xly-download-tab__operate');
		const p = toolbar && toolbar.__vue__ && toolbar.__vue__.$parent;
		if (p && p.taskBaseMap) {
			for (const id in p.taskBaseMap) {
				const t = p.taskBaseMap[id];
				if (t && t.taskName && (t.taskName.includes(clean) || clean.includes(t.taskName))) {
					p.downloadingSelectedIds = [t.taskId];
					p.currentSelectdTaskIds = [t.taskId];
					if (p.handleStart) { p.handleStart(); return { success: true }; }
				}
			}
		}
		return { success: false, error: "Task not found" };
	})('%s')`, clean)

	resJSON, err := EvaluateCDP(wsURL, script, 3*time.Second)
	if err != nil { return false, err }
	return strings.Contains(resJSON, `"success":true`), nil
}
