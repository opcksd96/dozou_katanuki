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
	if port <= 0 {
		port = 9222
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/json", port)
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to connect to CDP port %d: %w", port, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var targets []CDPPageTarget
	if err := json.Unmarshal(body, &targets); err != nil {
		return "", err
	}

	for _, t := range targets {
		if strings.Contains(t.URL, "main-renderer") || t.Title == "迅雷" {
			if t.WebSocketDebuggerURL != "" {
				return t.WebSocketDebuggerURL, nil
			}
		}
	}
	if len(targets) > 0 && targets[0].WebSocketDebuggerURL != "" {
		return targets[0].WebSocketDebuggerURL, nil
	}
	return "", fmt.Errorf("no valid main-renderer target found")
}

// ThunderExtractTaskScript は画面上のダウンロードリストを瞬時に抜き取るステルススクリプトです
const ThunderExtractTaskScript = `(() => {
	const sel = '.td-draglist-item, .xly-side-item, .xly-side-content';
	let items = Array.from(document.querySelectorAll(sel));
	if (items.length === 0) {
		const list = document.querySelector('.xly-side-list, .td-draglist, .xly-main');
		if (list) items = Array.from(list.children);
	}
	return items.map(el => el.innerText).filter(t => t && (t.includes('.jpg') || t.includes('.mp4') || t.includes('.png') || t.includes('.webp')));
})()`
