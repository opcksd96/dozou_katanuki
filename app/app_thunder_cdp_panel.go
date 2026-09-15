// app/app_thunder_cdp_panel.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// FindNewTaskPanelWSUrl は CDP ターゲットから「新建任务面板」の WebSocket URL を検索します
func FindNewTaskPanelWSUrl(port int) (string, error) {
	if port <= 0 { port = 9222 }
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/json", port))
	if err != nil { return "", err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var targets []CDPPageTarget
	_ = json.Unmarshal(body, &targets)

	for _, t := range targets {
		if strings.Contains(t.URL, "modeless-renderer") && (strings.Contains(t.Title, "新建任务") || strings.Contains(t.URL, "boxId=")) {
			return t.WebSocketDebuggerURL, nil
		}
	}
	return "", fmt.Errorf("新建任务面板 が見つかりません")
}

// EnsureCleanState は タスク投入前に開いている不要ダイアログを閉じ、安全な初期状態を整えます
func EnsureCleanState(port int) error {
	if port <= 0 { port = 9222 }
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/json", port))
	if err != nil { return fmt.Errorf("CDP 接続失敗: %w", err) }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var targets []CDPPageTarget
	_ = json.Unmarshal(body, &targets)

	for _, t := range targets {
		if strings.Contains(t.URL, "modeless-renderer") || strings.Contains(t.Title, "新建任务") {
			_, _ = EvaluateCDP(t.WebSocketDebuggerURL, `(() => {
				const c = document.querySelector('.xly-dialog-header__close, .xly-icon-close');
				if (c) { c.click(); return 'closed'; }
				window.close(); return 'window.close';
			})()`, 2*time.Second)
			time.Sleep(200 * time.Millisecond)
		}
		if strings.Contains(t.URL, "message-box") || strings.Contains(t.Title, "消息提示框") {
			_, _ = EvaluateCDP(t.WebSocketDebuggerURL, `(() => {
				const b = document.querySelector('button.td-button--primary, button.td-button');
				if (b) { b.click(); }
				window.close(); return 'closed';
			})()`, 2*time.Second)
			time.Sleep(200 * time.Millisecond)
		}
	}

	mainWS, err := FetchThunderMainRendererWSUrl(port)
	if err == nil {
		tabScript := `(() => {
			const dlTab = Array.from(document.querySelectorAll('.xly-sidebar-item, .xly-nav-item, div, span'))
				.find(el => (el.innerText || '').trim().startsWith('下载') || (el.innerText || '').trim().startsWith('我的下载'));
			if (dlTab) { dlTab.click(); }
			return true;
		})()`
		_, _ = EvaluateCDP(mainWS, tabScript, 2*time.Second)
	}
	return nil
}

// InjectGuardLock は 操作中の誤操作を防ぐためダイアログ全面にロックオーバーレイを注入します
func InjectGuardLock(wsURL string) error {
	script := `(() => {
		let old = document.getElementById('__cdp_guard_lock__');
		if (old) old.remove();
		const guard = document.createElement('div');
		guard.id = '__cdp_guard_lock__';
		guard.style.cssText = 'position:fixed;top:0;left:0;width:100vw;height:100vh;background:rgba(15,15,25,0.75);backdrop-filter:blur(3px);z-index:99999999;display:flex;flex-direction:column;align-items:center;justify-content:center;color:#fff;user-select:none;pointer-events:all;cursor:wait !important;';
		guard.innerHTML = '<div style="background:#1e1e2e;padding:22px 32px;border-radius:12px;box-shadow:0 12px 32px rgba(0,0,0,0.7);text-align:center;border:1px solid #313244;"><div style="font-size:18px;font-weight:bold;margin-bottom:8px;color:#89b4fa;">⚡ dozou 迅雷 Direct CDP 投入中</div><div style="font-size:13px;color:#cdd6f4;">URL入力・ファイル名正規化・保存先設定を実行しています...</div></div>';
		document.body.appendChild(guard);
		return true;
	})()`
	_, err := EvaluateCDP(wsURL, script, 3*time.Second)
	return err
}

// RemoveGuardLock は ロックオーバーレイを解除します
func RemoveGuardLock(wsURL string) {
	_, _ = EvaluateCDP(wsURL, `(() => { const g = document.getElementById('__cdp_guard_lock__'); if (g) g.remove(); return true; })()`, 2*time.Second)
}
