// app/app_rpc_thunder_cdp_detail.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ThunderTaskDetailInfo struct {
	Success      bool   `json:"success"`
	FileName     string `json:"file_name"`
	DownloadURL  string `json:"download_url"`
	SavePath     string `json:"save_path"`
	ProgressText string `json:"progress_text"`
	SpeedText    string `json:"speed_text"`
	PeersText    string `json:"peers_text"`
	CreatedAt    string `json:"created_at"`
	ElapsedTime  string `json:"elapsed_time"`
	RawDetail    string `json:"raw_detail"`
	Error        string `json:"error,omitempty"`
}

type CDPDetailEvalResult struct {
	Result struct {
		Result struct {
			Value ThunderTaskDetailInfo `json:"value"`
		} `json:"result"`
	} `json:"result"`
}

// GetThunderTaskDetailViaCDP は CDP 経由で迅雷の「文件详情（詳細パネル）」を開いてタスク詳細を取得します
func (a *App) GetThunderTaskDetailViaCDP(fileName string) (*ThunderTaskDetailInfo, error) {
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil || wsURL == "" {
		return nil, fmt.Errorf("迅雷 CDP に接続できません (:9222)")
	}

	cleanFileName := strings.TrimSpace(fileName)
	jsCode := fmt.Sprintf(`((targetFileName) => {
		const clean = targetFileName.trim();
		const allElements = Array.from(document.querySelectorAll('*'));
		const targetEl = allElements.find(el => el.children.length === 0 && el.innerText && el.innerText.trim().includes(clean));
		if (!targetEl) return { success: false, error: "Task element not found: " + clean };

		let curr = targetEl;
		for (let i = 0; i < 5 && curr; i++) {
			try {
				curr.click();
				curr.dispatchEvent(new MouseEvent('dblclick', { bubbles: true, cancelable: true, view: window }));
			} catch(e) {}
			curr = curr.parentElement;
		}

		const detailPanel = document.querySelector('.xly-detail, .task-detail');
		const rawText = detailPanel ? detailPanel.innerText : '';
		
		const details = Array.from(document.querySelectorAll('.detail')).map(d => d.innerText.trim());
		let dlUrl = '', savePath = '', created = '', elapsed = '';
		
		const urlMatch = rawText.match(/https?:\/\/[^\s]+/);
		if (urlMatch) dlUrl = urlMatch[0];

		for (let d of details) {
			if (d.includes(':\\') || d.startsWith('/')) savePath = d;
			else if (d.match(/^\d{4}-\d{2}-\d{2}/)) created = d;
			else if (d.match(/^\d{2}:\d{2}:\d{2}$/)) elapsed = d;
		}

		return {
			success: true,
			file_name: clean,
			download_url: dlUrl,
			save_path: savePath || 'D:\\迅雷下载',
			progress_text: (rawText.match(/下载进度\s*\d+%%/) || [''])[0],
			speed_text: (rawText.match(/\d+(?:\.\d+)?(?:B|KB|MB)\/s/) || ['0B/s'])[0],
			peers_text: (rawText.match(/资源数\s*\d+\/\d+/) || ['0/0'])[0],
			created_at: created,
			elapsed_time: elapsed,
			raw_detail: rawText.slice(0, 500)
		};
	})('%s')`, cleanFileName)

	resJSON, err := EvaluateCDPExpression(wsURL, jsCode, 2*time.Second)
	if err != nil { return nil, err }

	var res CDPDetailEvalResult
	if err := json.Unmarshal([]byte(resJSON), &res); err != nil { return nil, err }
	val := res.Result.Result.Value
	if !val.Success { return nil, fmt.Errorf("%s", val.Error) }
	return &val, nil
}
