// app/app_rpc_thunder_orchestrator_reaper.go (100行以下 - SPEC-PRINCIPLE-001)
package app

// ReapCompletedDuplicates は 1つの候補が完了した際、thunder_tasks DB と連動して他候補を迅雷から即座に取り下げます
func (a *App) ReapCompletedDuplicates(mediaID, completedFileName string) {
	if mediaID == "" || a.Repo == nil { return }

	reaps, err := a.Repo.MarkThunderTaskCompletedAndReapOthers(mediaID, completedFileName)
	if err != nil || len(reaps) == 0 { return }

	for _, task := range reaps {
		_, _ = DeleteTaskViaDirectCDP(9222, task.FileName)
	}
}

// ReapDepletedTask は 枯渇タスクを取り下げ、全候補が全滅した時のみ RETAINED へ退避します
func (a *App) ReapDepletedTask(fileName string) {
	if fileName == "" { return }
	_, _ = DeleteTaskViaDirectCDP(9222, fileName)

	if a.Repo != nil {
		allDepleted, mediaID, err := a.Repo.MarkThunderTaskDepleted(fileName)
		if err == nil && allDepleted && mediaID != "" {
			_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "迅雷: 全候補タスク枯渇によりRETAINED退避")
		}
	}
}

// deleteTaskByFileNameSilent は 指定ファイル名のタスクを迅雷からサイレントに取り下げ（Direct CDP）
func (a *App) deleteTaskByFileNameSilent(wsURL, fileName string) {
	if fileName == "" { return }
	_, _ = DeleteTaskViaDirectCDP(9222, fileName)
}

// deleteTaskByFileNameAndTextSilent は Direct CDP 削除のエイリアス
func (a *App) deleteTaskByFileNameAndTextSilent(wsURL, fileName, requireText string) {
	if fileName == "" { return }
	_, _ = DeleteTaskViaDirectCDP(9222, fileName)
}
