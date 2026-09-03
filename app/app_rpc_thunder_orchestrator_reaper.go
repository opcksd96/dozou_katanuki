// app/app_rpc_thunder_orchestrator_reaper.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"time"
)

// ReapCompletedDuplicates は 1つの候補が完了した際、thunder_tasks DB と連動して他候補を迅雷から即座に取り下げます
func (a *App) ReapCompletedDuplicates(mediaID, completedFileName string) {
	if mediaID == "" || a.Repo == nil {
		return
	}
	wsURL, err := FetchThunderMainRendererWSUrl(9222)
	if err != nil || wsURL == "" {
		return
	}

	reaps, err := a.Repo.MarkThunderTaskCompletedAndReapOthers(mediaID, completedFileName)
	if err != nil || len(reaps) == 0 {
		return
	}

	for _, task := range reaps {
		a.deleteTaskByFileNameSilent(wsURL, task.FileName)
	}
}

// ReapDepletedTask は 枯渇タスクを取り下げ、全候補が全滅した時のみ RETAINED へ退避します
func (a *App) ReapDepletedTask(fileName string) {
	if fileName == "" {
		return
	}
	wsURL, _ := FetchThunderMainRendererWSUrl(9222)
	if wsURL != "" {
		a.deleteTaskByFileNameSilent(wsURL, fileName)
	}

	if a.Repo != nil {
		allDepleted, mediaID, err := a.Repo.MarkThunderTaskDepleted(fileName)
		if err == nil && allDepleted && mediaID != "" {
			_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "迅雷: 全候補タスク枯渇によりRETAINED退避")
		}
	}
}

// deleteTaskByFileNameSilent は 指定ファイル名のタスクを迅雷からサイレントに取り下げ（削除）します
func (a *App) deleteTaskByFileNameSilent(wsURL, fileName string) {
	a.deleteTaskByFileNameAndTextSilent(wsURL, fileName, "")
}

// deleteTaskByFileNameAndTextSilent は ファイル名と特定のエラー文言（オプショナル）の両方に一致するタスクを安全に取り下げます
func (a *App) deleteTaskByFileNameAndTextSilent(wsURL, fileName, requireText string) {
	if wsURL == "" || fileName == "" {
		return
	}
	script := `(() => {
		const items = Array.from(document.querySelectorAll('.td-draglist-item, .xly-side-item, .xly-side-content'));
		for (let it of items) {
			const text = it.innerText || '';
			if (!text.includes('` + fileName + `')) continue;
			if ('` + requireText + `' !== '' && !text.includes('` + requireText + `')) continue;
			
			try { it.click(); } catch(e) {}
			
			let delBtn = it.querySelector('button[title="删除"], a[title="删除任务记录"], [title="删除"], .td-button[title*="删除"], [title*="彻底删除"]');
			if (!delBtn) {
				delBtn = document.querySelector('.xly-download-tab__operate button[title="删除"], button[title="删除"], .td-button[title*="删除"]');
			}
			
			if (delBtn) {
				delBtn.click();
				setTimeout(() => {
					const confirmBtn = Array.from(document.querySelectorAll('.td-dialog button, .td-dialog .td-button, .el-button, .xly-modal button, button')).find(b => b.innerText && (b.innerText.includes('确定') || b.innerText.includes('删除')));
					if (confirmBtn) confirmBtn.click();
				}, 150);
				return true;
			}
		}
		return false;
	})()`
	_, _ = EvaluateCDPExpression(wsURL, script, 1200*time.Millisecond)
}
