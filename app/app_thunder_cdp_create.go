// app/app_thunder_cdp_create.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dozou_katanuki/models"
)

// AddTasksViaDirectCDP は Direct CDP を介して迅雷へタスクを一括登録します
func AddTasksViaDirectCDP(port int, items []models.ThunderTaskInputItem, saveDir string) error {
	if len(items) == 0 { return fmt.Errorf("タスクが空です") }
	if err := EnsureCleanState(port); err != nil { return fmt.Errorf("事前検証に失敗: %w", err) }
	time.Sleep(300 * time.Millisecond)

	mainWS, err := FetchThunderMainRendererWSUrl(port)
	if err != nil { return fmt.Errorf("メインウィンドウ接続失敗: %w", err) }

	clickScript := `(() => {
		const b = Array.from(document.querySelectorAll('button.td-button')).find(el => (el.innerText || '').trim() === '新建');
		if (b) { b.click(); return true; }
		return false;
	})()`
	if _, err := EvaluateCDP(mainWS, clickScript, 3*time.Second); err != nil {
		return fmt.Errorf("新建ボタンのクリックに失敗: %w", err)
	}

	time.Sleep(600 * time.Millisecond)
	panelWS, err := FindNewTaskPanelWSUrl(port)
	if err != nil { return fmt.Errorf("新規タスクパネル検出失敗: %w", err) }

	_ = InjectGuardLock(panelWS)
	defer RemoveGuardLock(panelWS)

	var urls []string
	nameMap := make(map[string]string)
	for _, it := range items {
		urls = append(urls, it.URL)
		if it.FileName != "" { nameMap[it.URL] = it.FileName }
	}
	inputScript := fmt.Sprintf(`(() => {
		const ta = document.querySelector('textarea.td-textarea__inner');
		if (!ta) return false;
		ta.value = "%s";
		ta.dispatchEvent(new Event('input', { bubbles: true }));
		ta.dispatchEvent(new Event('change', { bubbles: true }));
		return true;
	})()`, strings.Join(urls, "\\n"))
	if _, err := EvaluateCDP(panelWS, inputScript, 3*time.Second); err != nil {
		return fmt.Errorf("URL 入力失敗: %w", err)
	}

	time.Sleep(800 * time.Millisecond)
	nameMapJSON, _ := json.Marshal(nameMap)
	cleanDir := ""
	if saveDir != "" && !strings.Contains(saveDir, "?") {
		cleanDir = strings.ReplaceAll(strings.ReplaceAll(saveDir, "/", `\`), `\`, `\\`)
	}
	configScript := fmt.Sprintf(`(() => {
		const nm = %s, dir = "%s";
		const norm = document.querySelector('.xly-dialog-normal');
		if (norm && norm.__vue__) {
			const v = norm.__vue__;
			if (v.dataMap) {
				for (let k in v.dataMap) { if (nm[k] && v.dataMap[k].data) { v.dataMap[k].data.fileName = nm[k]; } }
			}
			if (v.selectedDownlist) {
				v.selectedDownlist.forEach(it => { if (it.data && nm[it.data.url]) { it.data.fileName = nm[it.data.url]; } });
			}
		}
		if (dir) {
			const dInp = Array.from(document.querySelectorAll('input.td-input__inner')).find(el => el.placeholder === '请选择下载目录');
			if (dInp) {
				const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value").set;
				setter.call(dInp, dir);
				dInp.dispatchEvent(new Event('input', { bubbles: true }));
				dInp.dispatchEvent(new Event('change', { bubbles: true }));
			}
		}
		return true;
	})()`, string(nameMapJSON), cleanDir)
	_, _ = EvaluateCDP(panelWS, configScript, 3*time.Second)

	btnScript := `(() => {
		const b = document.querySelector('button.xly-button-down');
		if (b && !b.classList.contains('is-disabled')) { b.click(); return 'clicked'; }
		return 'disabled';
	})()`
	for i := 0; i < 15; i++ {
		time.Sleep(300 * time.Millisecond)
		res, err := EvaluateCDP(panelWS, btnScript, 2*time.Second)
		if err == nil && strings.Contains(res, "clicked") { return nil }
	}
	return fmt.Errorf("「立即下载」の待機タイムアウト")
}
