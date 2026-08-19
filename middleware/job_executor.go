// middleware/job_executor.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"dozou_katanuki/models"
)

func (j *JobOrchestrator) worker() {
	for {
		select {
		case <-j.ctx.Done():
			return
		case req, ok := <-j.queue:
			if !ok {
				return
			}
			j.executeJob(req)
		}
	}
}

func (j *JobOrchestrator) executeJob(req *models.JobRequest) {
	jobCtx, jobCancel := context.WithCancel(j.ctx)

	j.mu.Lock()
	progress, exists := j.jobs[req.ID]
	if !exists {
		progress = &models.JobProgress{ID: req.ID, Type: req.Type, Logs: make([]string, 0)}
		j.jobs[req.ID] = progress
	}
	now := time.Now()
	progress.Status = models.JobStatusRunning
	progress.StartedAt = &now
	progress.Message = "Starting subprocess..."
	j.activeJobID = req.ID
	j.activeCancel = jobCancel
	pythonBin := j.pythonPath
	j.mu.Unlock()

	j.emitEvent("job:started", progress)
	j.emitEvent("job:progress", progress)

	var cmd *exec.Cmd
	if req.ScriptPath != "" {
		cmd = exec.CommandContext(jobCtx, pythonBin, append([]string{req.ScriptPath}, req.Args...)...)
	} else if len(req.Args) > 0 {
		cmd = exec.CommandContext(jobCtx, req.Args[0], req.Args[1:]...)
	} else {
		j.finishJob(req.ID, models.JobStatusFailed, "Empty script and args", nil)
		return
	}

	if len(req.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range req.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, "Failed stdout pipe", err)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, "Failed stderr pipe", err)
		return
	}

	if err := cmd.Start(); err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, "Subprocess start failed", err)
		return
	}

	j.mu.Lock()
	j.activeCmd = cmd
	j.mu.Unlock()
	log.Printf("[JobOrchestrator] Started job %s (PID: %d): %s %v", req.ID, cmd.Process.Pid, pythonBin, req.ScriptPath)

	go j.scanStdoutProgress(req.ID, stdoutPipe)
	go j.scanStderr(req.ID, stderrPipe)

	err = cmd.Wait()
	if jobCtx.Err() == context.Canceled {
		j.finishJob(req.ID, models.JobStatusCancelled, "Job cancelled by user", nil)
	} else if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, fmt.Sprintf("Subprocess error: %v", err), err)
	} else {
		j.finishJob(req.ID, models.JobStatusCompleted, "Completed successfully", nil)
	}
}
