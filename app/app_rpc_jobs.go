// app/app_rpc_jobs.go (100行以下)
package app

import (
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

func (a *App) StartSalvageJob(platform, account, source string, limit int) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueSalvage(platform, account, source, limit)
}
func (a *App) StartManualImportJob(warcPath string, offline bool) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return a.JobOrchestrator.EnqueueManualImport(warcPath, offline)
}
func (a *App) StartTranslateJob(account string, overwrite bool) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return a.JobOrchestrator.EnqueueTranslate(account, overwrite)
}
func (a *App) GetActiveJob() (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return a.JobOrchestrator.GetActiveJob(), nil
}
func (a *App) ListJobs() ([]*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return a.JobOrchestrator.ListJobs(), nil
}
func (a *App) StartMediaDownloadJob(platform, mediaID string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueMediaDownload(platform, mediaID)
}
func (a *App) StartMediaPollJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueMediaPoll(platform)
}
func (a *App) StartMediaEscalateJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueMediaEscalate(platform)
}
func (a *App) StartSmartRecoveryJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueSmartRecovery(platform)
}
func (a *App) StartThunderEscalateJob(platform string) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }; if platform == "" { platform = "twitter" }
	return a.JobOrchestrator.EnqueueThunder(platform)
}
func (a *App) CancelJob(jobID string) error {
	if err := a.WaitForReady(); err != nil { return err }
	return a.JobOrchestrator.CancelJob(jobID)
}

// GetLatestMissionReport は直近のミッション完了テキストレポート（5W1H）を取得します
func (a *App) GetLatestMissionReport() (*models.MissionReport, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return middleware.GetLatestReport()
}

// GetSystemJournals はインフラ常駐処理のオンメモリ構造化JSONジャーナルを取得します
func (a *App) GetSystemJournals(limit int) ([]models.SystemJournalEntry, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return middleware.GetGlobalJournal().GetEntries(limit), nil
}
