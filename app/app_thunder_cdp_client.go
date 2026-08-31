// app/app_thunder_cdp_client.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CDPPageTarget struct {
	Title                string `json:"title"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

// FetchThunderMainRendererWSUrl は 迅雷の main-renderer ページの WebSocket URL を取得します
func FetchThunderMainRendererWSUrl(port int) (string, error) {
	if port <= 0 { port = 9222 }
	url := fmt.Sprintf("http://127.0.0.1:%d/json", port)
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil { return "", fmt.Errorf("failed to connect to CDP port %d: %w", port, err) }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil { return "", err }

	var targets []CDPPageTarget
	if err := json.Unmarshal(body, &targets); err != nil { return "", err }

	for _, t := range targets {
		if strings.Contains(t.URL, "main-renderer") || t.Title == "迅雷" {
			if t.WebSocketDebuggerURL != "" { return t.WebSocketDebuggerURL, nil }
		}
	}
	if len(targets) > 0 && targets[0].WebSocketDebuggerURL != "" { return targets[0].WebSocketDebuggerURL, nil }
	return "", fmt.Errorf("no valid main-renderer target found")
}

// ThunderExtractTaskScript は画面静止ロック付きで裏でダウンロードリストを抜き取るステルススクリプトです
const ThunderExtractTaskScript = `(async () => {
	const navs = Array.from(document.querySelectorAll('.xly-nav__item, .xly-side-menu__item, [class*="nav__item"]'));
	const active = navs.find(it => it.classList.contains('is-active') || it.classList.contains('active') || it.classList.contains('selected'));
	const curTab = active ? (active.querySelector('.xly-nav__tab, .xly-nav__content') || active).innerText.split('\n')[0].trim() : '下载中';

	if (curTab === '下载中') {
		const list = document.querySelector('.xly-side-list');
		return list ? Array.from(list.querySelectorAll('.xly-side-content')).map(el => el.innerText) : [];
	}

	const area = document.querySelector('.xly-side-list, .xly-main, #app') || document.body;
	const overlay = area.cloneNode(true);
	overlay.id = '__stealth_freeze__';
	overlay.style.cssText = 'position:absolute;top:' + area.offsetTop + 'px;left:' + area.offsetLeft + 'px;width:' + area.offsetWidth + 'px;height:' + area.offsetHeight + 'px;z-index:999999;pointer-events:none;';
	document.body.appendChild(overlay);

	try {
		const allTabs = Array.from(document.querySelectorAll('.xly-nav__tab, span, a, .xly-nav__item'));
		const dlTab = allTabs.find(el => el.innerText && el.innerText.trim() === '下载中');
		const origTab = allTabs.find(el => el.innerText && el.innerText.trim() === curTab);

		if (dlTab) dlTab.click();
		await new Promise(r => setTimeout(r, 60));

		const list = document.querySelector('.xly-side-list');
		const tasks = list ? Array.from(list.querySelectorAll('.xly-side-content')).map(el => el.innerText) : [];

		if (origTab) origTab.click();
		await new Promise(r => setTimeout(r, 30));
		return tasks;
	} finally {
		const ov = document.getElementById('__stealth_freeze__');
		if (ov) ov.remove();
	}
})()`
