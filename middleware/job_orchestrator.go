// middleware/job_orchestrator.go (100行以下)
package middleware

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"dozou_katanuki/models"
)

type EventEmitter func(eventName string, optionalData ...interface{})

// JobOrchestrator は Python サイドカープロセスを排他キュー（最大並行数1）で管理・実行するエンジンです
type JobOrchestrator struct {
	ctx          context.Context
	cancel       context.CancelFunc
	emitter      EventEmitter
	queue        chan *models.JobRequest
	mu           sync.RWMutex
	jobs         map[string]*models.JobProgress
	activeCancel context.CancelFunc
	activeCmd    *exec.Cmd
	activeJobID  string
	pythonPath   string
	storageArgs  []string
	maxLogs      int
}

func NewJobOrchestrator(ctx context.Context, emitter EventEmitter) *JobOrchestrator {
	cCtx, cCancel := context.WithCancel(ctx)
	orch := &JobOrchestrator{
		ctx:        cCtx,
		cancel:     cCancel,
		emitter:    emitter,
		queue:      make(chan *models.JobRequest, 100),
		jobs:       make(map[string]*models.JobProgress),
		pythonPath: "python",
		maxLogs:    500,
	}

	if p := os.Getenv("PYTHON_BIN"); p != "" {
		orch.pythonPath = p
	}

	go orch.worker()
	return orch
}

func (j *JobOrchestrator) SetPythonPath(path string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.pythonPath = path
}

func (j *JobOrchestrator) SetStorageConfig(cfg models.StorageConfig) {
	j.mu.Lock()
	defer j.mu.Unlock()
	var args []string
	if cfg.DBPath != "" {
		args = append(args, "--db-path", cfg.DBPath)
	}
	if cfg.LocalMediaDir != "" {
		args = append(args, "--storage-dir", cfg.LocalMediaDir)
	} else if cfg.StashDir != "" {
		args = append(args, "--storage-dir", cfg.StashDir)
	}
	j.storageArgs = args
}

func (j *JobOrchestrator) Close() {
	j.cancel()
	j.mu.Lock()
	if j.activeCancel != nil {
		j.activeCancel()
	}
	if j.activeCmd != nil && j.activeCmd.Process != nil {
		_ = j.activeCmd.Process.Kill()
	}
	j.mu.Unlock()
}
