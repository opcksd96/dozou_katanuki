// app/app_rpc_thunder_cdp_control.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CDPControlResult struct {
	Result struct {
		Result struct {
			Value struct {
				Success bool   `json:"success"`
				Action  string `json:"action"`
				Error   string `json:"error"`
			} `json:"value"`
		} `json:"result"`
	} `json:"result"`
}

// ControlThunderTaskViaCDP は CDP 経由で迅雷内部のタスクを操作します (pause / resume / delete / restore)
func (a *App) ControlThunderTaskViaCDP(fileName string, action string) (bool, error) {
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil || wsURL == "" { return false, fmt.Errorf("迅雷 CDP に接続できません (:9222)") }
	cleanFileName := strings.TrimSpace(fileName)
	if action == "restore" && cleanFileName == "" { return false, fmt.Errorf("restore action requires a specific target filename") }

	jsCode := fmt.Sprintf(`((targetFileName, action) => {
		const clean = targetFileName ? targetFileName.trim() : '';
		if (action === 'restore') {
			const recycleTab = Array.from(document.querySelectorAll('.xly-nav__tab, span, a')).find(el => el.innerText && el.innerText.trim() === '回收站');
			if (recycleTab) recycleTab.click();
		}

		const allElements = Array.from(document.querySelectorAll('*'));
		const targetEl = clean ? allElements.find(el => el.children.length === 0 && el.innerText && el.innerText.trim().includes(clean)) : null;
		let curr = targetEl, row = null;
		for (let i = 0; i < 6 && curr; i++) {
			try { curr.click(); } catch(e) {}
			if (curr.classList && (curr.classList.contains('td-media') || curr.classList.contains('td-draglist-item') || curr.classList.contains('xly-side-content__item'))) row = curr;
			curr = curr.parentElement;
		}

		if (action === 'restore') {
			const restoreBtn = (row && row.querySelector('[title*="还原"], .xly-side-operate__button[title*="还原"]')) || document.querySelector('[title*="还原"], .xly-side-operate__button[title*="还原"]');
			if (restoreBtn) { restoreBtn.click(); return { success: true, action: action, target: clean }; }
			return { success: false, error: "Restore button not found in recycle bin" };
		}

		if (action === 'delete') {
			const toolbar = document.querySelector('.xly-download-tab__operate');
			const p = toolbar && toolbar.__vue__ && toolbar.__vue__.$parent;
			if (p && p.taskBaseMap) {
				const map = p.taskBaseMap;
				let targetTask = null;
				for (const id in map) {
					if (map[id] && map[id].taskName && (map[id].taskName.includes(clean) || clean.includes(map[id].taskName))) {
						targetTask = map[id]; break;
					}
				}
				if (targetTask) {
					p.downloadingSelectedIds = [targetTask.taskId];
					p.currentSelectdTaskIds = [targetTask.taskId];
					try { p.handleDelete({ key: "Delete", keyCode: 46 }); return { success: true, action: action, target: clean }; } catch(e) {}
				}
			}
			return { success: false, error: "Task not found in taskBaseMap or delete failed" };
		}

		const titleMatch = (action === 'pause') ? '暂停' : '下载';
		const btn = (row && row.querySelector('[title*="' + titleMatch + '"]')) || document.querySelector('button[title*="' + titleMatch + '"], .td-button[title*="' + titleMatch + '"], [title="' + titleMatch + '"]');
		if (btn) { btn.click(); return { success: true, action: action, target: clean }; }
		return { success: false, error: "Button not found: " + titleMatch };
	})('%s', '%s')`, cleanFileName, action)

	resJSON, err := EvaluateCDPExpression(wsURL, jsCode, 2*time.Second)
	if err != nil { return false, err }
	var res CDPControlResult
	if err := json.Unmarshal([]byte(resJSON), &res); err != nil { return false, err }
	if !res.Result.Result.Value.Success { return false, fmt.Errorf("%s", res.Result.Result.Value.Error) }
	return true, nil
}
