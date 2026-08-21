// middleware/job_state.go (100行以下)
package middleware

import (
	"fmt"
	"log"
	"time"

	"dozou_katanuki/models"
)

func (j *JobOrchestrator) finishJob(jobID string, status models.JobStatus, msg string, err error) {
	now := time.Now()
	j.mu.Lock()
	var snapshot *models.JobProgress
	if p, ok := j.jobs[jobID]; ok {
		p.Status, p.Message, p.FinishedAt = status, msg, &now
		if err != nil { p.Error = err.Error() }
		if status == models.JobStatusCompleted {
			p.Percentage = 100.0
			if p.Total > 0 && p.Current < p.Total { p.Current = p.Total }
		}
		s := *p
		snapshot = &s
	}
	if j.activeJobID == jobID { j.activeJobID, j.activeCancel, j.activeCmd = "", nil, nil }
	j.mu.Unlock()

	log.Printf("[JobOrchestrator] Finished job %s: status=%s, msg=%s", jobID, status, msg)
	if snapshot != nil {
		j.emitEvent("job:finished", snapshot)
		j.emitEvent("job:progress", snapshot)
	}
}

func (j *JobOrchestrator) GetStatus(jobID string) *models.JobProgress {
	j.mu.RLock(); defer j.mu.RUnlock()
	if p, ok := j.jobs[jobID]; ok { s := *p; return &s }
	return nil
}

func (j *JobOrchestrator) GetActiveJob() *models.JobProgress {
	j.mu.RLock(); defer j.mu.RUnlock()
	if j.activeJobID != "" {
		if p, ok := j.jobs[j.activeJobID]; ok { s := *p; return &s }
	}
	return nil
}

func (j *JobOrchestrator) ListJobs() []*models.JobProgress {
	j.mu.RLock(); defer j.mu.RUnlock()
	list := make([]*models.JobProgress, 0, len(j.jobs))
	for _, p := range j.jobs { s := *p; list = append(list, &s) }
	return list
}

func (j *JobOrchestrator) CancelJob(jobID string) error {
	j.mu.Lock(); defer j.mu.Unlock()
	p, ok := j.jobs[jobID]
	if !ok { return fmt.Errorf("job not found: %s", jobID) }
	if p.Status == models.JobStatusRunning && j.activeJobID == jobID {
		if j.activeCancel != nil { j.activeCancel() }
		if j.activeCmd != nil && j.activeCmd.Process != nil { _ = j.activeCmd.Process.Kill() }
		return j.setCancelled(p, "Cancelled by user")
	}
	if p.Status == models.JobStatusPending { return j.setCancelled(p, "Cancelled before execution") }
	return fmt.Errorf("job is not running or pending (current status: %s)", p.Status)
}

func (j *JobOrchestrator) setCancelled(p *models.JobProgress, msg string) error {
	now := time.Now()
	p.Status, p.Message, p.FinishedAt = models.JobStatusCancelled, msg, &now
	snapshot := *p
	go func() {
		j.emitEvent("job:finished", &snapshot)
		j.emitEvent("job:progress", &snapshot)
	}()
	return nil
}
