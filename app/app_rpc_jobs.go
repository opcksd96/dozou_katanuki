// app/app_rpc_jobs.go (100行以下)
package app

import (
	"dozou_katanuki/models"
)

// StartSalvageJob は自動サルベージジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartSalvageJob(platform, account, source string, limit int) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.EnqueueSalvage(platform, account, source, limit, a.getTranslationEnv())
}

// StartManualImportJob は手動 WARC インポートジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartManualImportJob(warcPath string, offline bool) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.EnqueueManualImport(warcPath, offline)
}

// StartTranslateJob は未翻訳記事の一括翻訳ジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartTranslateJob(account string, overwrite bool) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.EnqueueTranslate(account, overwrite, a.getTranslationEnv())
}

// GetJobStatus は指定されたジョブのステータスを取得する Wails バインドメソッドです
func (a *App) GetJobStatus(jobID string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.GetStatus(jobID), nil
}

// GetActiveJob は現在実行中のジョブを取得する Wails バインドメソッドです
func (a *App) GetActiveJob() (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.GetActiveJob(), nil
}

// ListJobs は全ジョブの履歴ステータス一覧を取得する Wails バインドメソッドです
func (a *App) ListJobs() ([]*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.JobOrchestrator.ListJobs(), nil
}

// StartMediaDownloadJob はメディアダウンロードジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartMediaDownloadJob(platform, mediaID string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if platform == "" {
		platform = "twitter"
	}
	return a.JobOrchestrator.EnqueueMediaDownload(platform, mediaID)
}

// StartMediaPollJob は外部委託（Aria2/Motrix）メディア回収ジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartMediaPollJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if platform == "" {
		platform = "twitter"
	}
	return a.JobOrchestrator.EnqueueMediaPoll(platform)
}

// StartMediaEscalateJob は DEAD_404 メディアの Motrix/Aria2 外注ジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartMediaEscalateJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if platform == "" {
		platform = "twitter"
	}
	return a.JobOrchestrator.EnqueueMediaEscalate(platform)
}

// CancelJob はジョブの実行を中断する Wails バインドメソッドです
func (a *App) CancelJob(jobID string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	return a.JobOrchestrator.CancelJob(jobID)
}
