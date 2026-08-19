// middleware/job_queue.go (100行以下)
package middleware

import (
	"fmt"
	"strconv"
	"time"

	"dozou_katanuki/models"
)

func (j *JobOrchestrator) EnqueueSalvage(platform, account string, limit int) (*models.JobProgress, error) {
	if limit <= 0 { limit = 50 }
	return j.EnqueueJob(&models.JobRequest{
		ID: fmt.Sprintf("job_salvage_%d", time.Now().UnixNano()), Type: models.JobTypeSalvage,
		Platform: platform, Account: account, Limit: limit,
		ScriptPath: fmt.Sprintf("plugins/%s/scraper/main.py", platform),
		Args: []string{"--mode", "auto", "--platform", platform, "--account", account, "--limit", strconv.Itoa(limit)},
		CreatedAt: time.Now(),
	})
}

func (j *JobOrchestrator) EnqueueManualImport(warcPath string, offline bool) (*models.JobProgress, error) {
	args := []string{"--mode", "manual", "--warc-path", warcPath}
	if offline { args = append(args, "--offline") }
	return j.EnqueueJob(&models.JobRequest{
		ID: fmt.Sprintf("job_import_%d", time.Now().UnixNano()), Type: models.JobTypeImportManual,
		WARCPath: warcPath, Offline: offline, ScriptPath: "plugins/twitter/scraper/main.py",
		Args: args, CreatedAt: time.Now(),
	})
}

func (j *JobOrchestrator) EnqueueMediaDownload(platform, mediaID string) (*models.JobProgress, error) {
	args := []string{"--mode", "download", "--platform", platform}
	if mediaID != "" { args = append(args, "--media-id", mediaID) }
	return j.EnqueueJob(&models.JobRequest{
		ID: fmt.Sprintf("job_dl_%d", time.Now().UnixNano()), Type: models.JobTypeMediaDownload,
		Platform: platform, ScriptPath: fmt.Sprintf("plugins/%s/scraper/main.py", platform),
		Args: args, CreatedAt: time.Now(),
	})
}

func (j *JobOrchestrator) EnqueueJob(req *models.JobRequest) (*models.JobProgress, error) {
	if req.ID == "" { req.ID = fmt.Sprintf("job_%d", time.Now().UnixNano()) }
	if req.CreatedAt.IsZero() { req.CreatedAt = time.Now() }

	progress := &models.JobProgress{
		ID: req.ID, Type: req.Type, Status: models.JobStatusPending,
		Current: 0, Total: req.Limit, Percentage: 0,
		Message: "Queued (Waiting for worker)", Logs: make([]string, 0),
	}

	j.mu.Lock()
	j.jobs[req.ID] = progress
	j.mu.Unlock()

	select {
	case j.queue <- req:
		j.emitEvent("job:queued", progress)
		return progress, nil
	default:
		j.mu.Lock()
		progress.Status = models.JobStatusFailed
		progress.Message = "Queue is full"
		progress.Error = "job queue overflow"
		j.mu.Unlock()
		return progress, fmt.Errorf("job queue is full")
	}
}
