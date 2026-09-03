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
	if err != nil || wsURL == "" {
		return false, fmt.Errorf("迅雷 CDP に接続できません (:9222)")
	}

	cleanFileName := strings.TrimSpace(fileName)
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
			if (curr.classList && (curr.classList.contains('td-media') || curr.classList.contains('xly-side-content__item'))) row = curr;
			curr = curr.parentElement;
		}

		if (action === 'restore') {
			const restoreBtn = (row && row.querySelector('[title*="还原"], .xly-side-operate__button[title*="还原"]')) || document.querySelector('[title*="还原"], .xly-side-operate__button[title*="还原"]');
			if (restoreBtn) { restoreBtn.click(); return { success: true, action: action, target: clean }; }
			return { success: false, error: "Restore button not found in recycle bin" };
		}

		if (action === 'delete') {
			let delBtn = (row && row.querySelector('[title*="删除"], .xly-side-operate__button[title*="删除"]')) || document.querySelector('[title*="删除任务记录"], [title="删除"]');
			if (delBtn) {
				delBtn.click();
				setTimeout(() => {
					const confirmBtn = Array.from(document.querySelectorAll('.td-dialog button, .td-dialog .td-button, .xly-modal button')).find(b => b.innerText.includes('确定') || b.innerText.includes('删除'));
					if (confirmBtn) confirmBtn.click();
				}, 120);
				return { success: true, action: action, target: clean };
			}
			return { success: false, error: "Delete button not found" };
		}

		const titleMatch = (action === 'pause') ? '暂停' : '下载';
		const btn = (row && row.querySelector('[title*="' + titleMatch + '"]')) || document.querySelector('button[title*="' + titleMatch + '"], .td-button[title*="' + titleMatch + '"], [title="' + titleMatch + '"]');
		if (btn) { btn.click(); return { success: true, action: action, target: clean }; }
		return { success: false, error: "Button not found: " + titleMatch };
	})('%s', '%s')`, cleanFileName, action)

	resJSON, err := EvaluateCDPExpression(wsURL, jsCode, 2*time.Second)
	if err != nil {
		return false, err
	}

	var res CDPControlResult
	if err := json.Unmarshal([]byte(resJSON), &res); err != nil {
		return false, err
	}
	if !res.Result.Result.Value.Success {
		return false, fmt.Errorf("%s", res.Result.Result.Value.Error)
	}
	return true, nil
}
