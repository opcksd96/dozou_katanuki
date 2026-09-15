// app/app_thunder_cdp_delete.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DeleteTaskViaDirectCDP は Direct CDP 経由で指定タスクを安全に停止させてからゴミ箱（回收站）へ移動します
func DeleteTaskViaDirectCDP(port int, targetFileName string) (bool, error) {
	clean := strings.TrimSpace(targetFileName)
	if clean == "" { return false, fmt.Errorf("targetFileName is empty") }

	wsURL, err := FetchThunderMainRendererWSUrl(port)
	if err != nil { return false, err }

	jsCode := fmt.Sprintf(`((targetFileName) => {
		const clean = targetFileName.trim();
		if (!clean) return { success: false, error: "Empty targetFileName" };
		const stripParen = (s) => (s || '').replace(/\s*\(\d+\)(\.[a-zA-Z0-9]+)$/, '$1').trim();
		const targetNorm = stripParen(clean);

		const toolbar = document.querySelector('.xly-download-tab__operate');
		const p = toolbar && toolbar.__vue__ && toolbar.__vue__.$parent;
		if (p && p.taskBaseMap) {
			const map = p.taskBaseMap;
			let targetTask = null;
			for (const id in map) {
				const t = map[id];
				if (!t || !t.taskName) continue;
				const tName = t.taskName.trim();
				if (tName === clean || stripParen(tName) === targetNorm) {
					targetTask = t; break;
				}
			}
			if (targetTask) {
				// ⚡ 排队等待・待機中ガード: taskStatus 0, 5 は待機中ジョブのため絶対に削除しない！
				if (targetTask.taskStatus === 0 || targetTask.taskStatus === 5) {
					return { success: false, error: "Task is queued or waiting (taskStatus=" + targetTask.taskStatus + "). Delete aborted." };
				}
				// ⚡ 通信中ガード: ダウンロード速度が出ているタスクは誤削除防止
				if (targetTask.downloadSpeed > 0 || targetTask.vipSpeed > 0) {
					return { success: false, error: "Task is actively downloading (speed > 0). Delete aborted." };
				}
				p.downloadingSelectedIds = [targetTask.taskId];
				p.currentSelectdTaskIds = [targetTask.taskId];
				try {
					p.handleDelete({ key: "Delete", keyCode: 46 });
					// ⚡ 万一「正在下载」「排队中」確認ダイアログが出た場合は即座にキャンセルして閉じる
					setTimeout(() => {
						const dlg = document.querySelector('.td-dialog, .xly-modal, .el-dialog');
						if (dlg && (dlg.innerText.includes('正在下载') || dlg.innerText.includes('排队') || dlg.innerText.includes('等待'))) {
							const cancel = Array.from(dlg.querySelectorAll('button')).find(b => (b.innerText || '').includes('取消'));
							if (cancel) cancel.click();
						}
					}, 100);
					return { success: true, target: clean };
				} catch(e) {
					return { success: false, error: e.toString() };
				}
			}
		}
		return { success: false, error: "Task not found in taskBaseMap: " + clean };
	})('%s')`, clean)

	resJSON, err := EvaluateCDP(wsURL, jsCode, 3*time.Second)
	if err != nil { return false, err }

	var evalRes struct {
		Result struct {
			Result struct {
				Value struct {
					Success bool   `json:"success"`
					Error   string `json:"error"`
				} `json:"value"`
			} `json:"result"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resJSON), &evalRes); err != nil {
		return false, fmt.Errorf("failed to unmarshal delete result: %w", err)
	}
	if !evalRes.Result.Result.Value.Success {
		return false, fmt.Errorf("delete failed: %s", evalRes.Result.Result.Value.Error)
	}
	return true, nil
}
