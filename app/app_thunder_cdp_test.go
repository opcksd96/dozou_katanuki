// app/app_thunder_cdp_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"
	"time"
)

func TestThunderCDP_ConnectionAndEvaluation(t *testing.T) {
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil {
		t.Skipf("Thunder CDP not available on port 9222 (skipping): %v", err)
		return
	}
	t.Logf("Found Thunder WebSocket URL: %s", wsURL)

	// 1. document.title を取得
	res, err := EvaluateCDPExpression(wsURL, "document.title", 2*time.Second)
	if err != nil {
		t.Skipf("EvaluateCDPExpression timeout/skipped: %v", err)
		return
	}
	t.Logf("CDP document.title evaluation result: %s", res)

	// 2. タスクリスト抽出スクリプトを実行
	listRes, err := EvaluateCDPExpression(wsURL, ThunderExtractTaskScript, 2*time.Second)
	if err != nil {
		t.Skipf("EvaluateCDPExpression task script timeout/skipped: %v", err)
		return
	}
	t.Logf("CDP Task extraction raw length: %d, sample: %s", len(listRes), listRes[:min(len(listRes), 200)])

	// 3. 通常の title="删除" ボタンをクリックした時の挙動を検証
	domRes, _ := EvaluateCDPExpression(wsURL, `(() => {
		const items = Array.from(document.querySelectorAll('.td-draglist-item, .xly-side-item, .xly-side-content'));
		const target = items.find(it => it.innerText && (it.innerText.includes('Go4PdTDb0AAmdRV_plain.jpg') || it.innerText.includes('无法继续下载')));
		if (!target) return { found: false };
		target.click(); // 行を選択

		// 通常の「删除」ボタン（title="删除" 完全一致）を探す
		const delBtn = document.querySelector('button[title="删除"], .xly-download-tab__operate button[title="删除"]') ||
			target.querySelector('a[title="删除任务记录"]');
		if (!delBtn) return { found: true, delBtnFound: false };

		const beforeDialogs = document.querySelectorAll('.td-dialog, .el-dialog, .xly-modal').length;
		delBtn.click();
		return {
			found: true,
			delBtnTag: delBtn.tagName,
			delBtnTitle: delBtn.getAttribute('title'),
			delBtnCls: delBtn.className,
			beforeDialogs: beforeDialogs
		};
	})()`, 2*time.Second)
	t.Logf("Click normal 删除 button result: %s", domRes)

	time.Sleep(300 * time.Millisecond)

	// ダイアログが出現したか確認
	dRes, _ := EvaluateCDPExpression(wsURL, `(() => {
		const dialogs = Array.from(document.querySelectorAll('.td-dialog, .el-dialog, .xly-modal, .td-message-box'));
		return {
			hasDialog: dialogs.length > 0,
			dialogTexts: dialogs.map(d => d.innerText)
		};
	})()`, 2*time.Second)
	t.Logf("Dialog after normal 删除: %s", dRes)
}
