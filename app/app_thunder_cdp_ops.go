// app/app_thunder_cdp_ops.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"time"

	"dozou_katanuki/models"
)

// GetTasksViaDirectCDP は Direct CDP 経由で迅雷の全タスク情報（生JSON含む）を取得します
func GetTasksViaDirectCDP(port int) ([]models.ActiveThunderTask, error) {
	wsURL, err := FetchThunderMainRendererWSUrl(port)
	if err != nil { return nil, err }

	script := `(() => {
		const fmt = (b) => {
			if (!b || b <= 0) return '0B';
			const k = 1024, s = ['B', 'KB', 'MB', 'GB', 'TB'];
			const i = Math.floor(Math.log(b) / Math.log(k));
			return (b / Math.pow(k, i)).toFixed(1) + s[i];
		};
		const domMap = {};
		document.querySelectorAll('.xly-side-item').forEach(el => {
			const tEl = el.querySelector('.xly-file-name__ad, .xly-file-name, .xly-side-title');
			if (tEl && tEl.innerText) {
				domMap[tEl.innerText.trim()] = (el.innerText || '').split('\n').map(s => s.trim()).filter(Boolean);
			}
		});
		let list = [];
		const tb = document.querySelector('.xly-download-tab__operate');
		if (tb && tb.__vue__ && tb.__vue__.taskBaseMap) {
			const map = tb.__vue__.taskBaseMap;
			for (let k in map) {
				const t = map[k];
				if (!t || !t.taskName) continue;
				const domLines = domMap[t.taskName] || [];
				const domText = domLines.join(' | ');
				let status = '⚪ 待機中', detail = '';
				if (domText.includes('无法继续下载') || domText.includes('原始资源不存在') || domText.includes('暂无任何有效资源') || t.errorCode === 402) {
					status = '🔴 リソース枯渇';
					detail = domLines.find(l => l.includes('资源') || l.includes('无法')) || '原始资源不存在，且未找到候选资源，无法继续下载';
				} else if (domText.includes('连接资源') || (t.taskStatus === 0 && (!t.downloadSpeed || t.downloadSpeed === 0))) {
					status = '🟡 ピア探索中';
					detail = 'ピア接続待機中';
				} else if (t.downloadSpeed > 0 || domText.includes('正在下载') || domText.includes('下载中')) {
					status = '🟢 ダウンロード中';
				} else if (domText.includes('暂停') || t.taskStatus === 7 || t.taskStatus === 8) {
					status = '⏸️ 一時停止';
				} else if (domText.includes('已完成') || t.taskStatus === 11) {
					status = '✅ 完了';
				}
				const prog = (t.fileSize > 0) ? (fmt(t.downloadSize) + ' / ' + fmt(t.fileSize)) : fmt(t.downloadSize);
				const spd = t.downloadSpeed ? (fmt(t.downloadSpeed) + '/s') : '0B/s';
				list.push({
					file_name: t.taskName, target: t.url || t.taskName, status: status,
					progress_text: prog, speed_text: spd, detail_text: detail,
					task_id: t.taskId || 0, task_type: t.taskType || 0,
					task_status_code: t.taskStatus || 0, error_code: t.errorCode || 0,
					url: t.url || '', save_path: t.savePath || '',
					file_size: t.fileSize || 0, download_size: t.downloadSize || 0,
					download_speed: t.downloadSpeed || 0, vip_speed: t.vipSpeed || 0,
					create_time: t.createTime || 0, completion_time: t.completionTime || 0,
					src_total: t.srcTotal || 0, src_using: t.srcUsing || 0,
					peers: t.nTotalAvailablePeer || 0, cid: t.cid || '', gcid: t.gcid || '',
					origin: t.origin || '', raw_json: JSON.stringify(t), raw_text: domText
				});
			}
		}
		return list;
	})()`

	resJSON, err := EvaluateCDP(wsURL, script, 3*time.Second)
	if err != nil { return nil, err }

	var evalRes struct {
		Result struct {
			Result struct {
				Value []models.ActiveThunderTask `json:"value"`
			} `json:"result"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resJSON), &evalRes); err != nil {
		return nil, fmt.Errorf("failed to parse CDP tasks response: %w", err)
	}
	return evalRes.Result.Result.Value, nil
}
