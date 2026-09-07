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
	if len(targets) > 0 && targets[0].WebSocketDebuggerURL != "" {
		return targets[0].WebSocketDebuggerURL, nil
	}
	return "", fmt.Errorf("no valid main-renderer target found")
}

// ThunderExtractTaskScript は内部Vueデータ(taskBaseMap)と.xly-side-itemから正確なタスク行を抽出します
const ThunderExtractTaskScript = `(() => {
	const fmt = (b) => {
		if (!b || b <= 0) return '0B';
		const k = 1024, s = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(b) / Math.log(k));
		return (b / Math.pow(k, i)).toFixed(2) + s[i];
	};
	let sizeMap = {};
	const tb = document.querySelector('.xly-download-tab__operate');
	if (tb && tb.__vue__ && tb.__vue__.taskBaseMap) {
		const m = tb.__vue__.taskBaseMap;
		for (let k in m) { if (m[k] && m[k].taskName) sizeMap[m[k].taskName] = fmt(m[k].fileSize); }
	}
	const items = Array.from(document.querySelectorAll('.xly-side-item'));
	return items.filter(el => !el.querySelector('.xly-icon-restore')).map(el => {
		let text = el.innerText || '';
		const titleEl = el.querySelector('.xly-file-name__ad, .xly-file-name, .xly-side-title');
		const title = titleEl ? titleEl.innerText.trim() : '';
		if (title && sizeMap[title]) { text = text + '\n' + sizeMap[title]; }
		return text;
	}).filter(t => t && (t.includes('.jpg') || t.includes('.mp4') || t.includes('.png') || t.includes('.webp')));
})()`
