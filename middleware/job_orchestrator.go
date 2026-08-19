package middleware

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"dozou_katanuki/models"
)

var progressRegex = regexp.MustCompile(`^PROGRESS:\s*(\d+)/(\d+)\s*\|\s*(.*)$`)

// EventEmitter は進捗イベント通知用のコールバック関数型です
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
	maxLogs      int
}

// NewJobOrchestrator creates a new JobOrchestrator instance
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

	// Python 実行パスの検出（環境変数 PYTHON_BIN または PATH）
	if p := os.Getenv("PYTHON_BIN"); p != "" {
		orch.pythonPath = p
	}

	// ワーカーゴルーチン起動（並行数1の厳格な排他処理）
	go orch.worker()

	return orch
}

// SetPythonPath は Python のバイナリパスを明示的に設定します
func (j *JobOrchestrator) SetPythonPath(path string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.pythonPath = path
}

// EnqueueSalvage は自動サルベージジョブをキューに追加します
func (j *JobOrchestrator) EnqueueSalvage(platform, account string, limit int) (*models.JobProgress, error) {
	if limit <= 0 {
		limit = 50
	}
	jobID := fmt.Sprintf("job_salvage_%d", time.Now().UnixNano())
	scriptPath := fmt.Sprintf("plugins/%s/scraper/main.py", platform)

	req := &models.JobRequest{
		ID:         jobID,
		Type:       models.JobTypeSalvage,
		Platform:   platform,
		Account:    account,
		Limit:      limit,
		ScriptPath: scriptPath,
		Args: []string{
			"--account", account,
			"--limit", strconv.Itoa(limit),
		},
		CreatedAt: time.Now(),
	}

	return j.EnqueueJob(req)
}

// EnqueueManualImport は手動 WARC インポートジョブをキューに追加します
func (j *JobOrchestrator) EnqueueManualImport(warcPath string, offline bool) (*models.JobProgress, error) {
	jobID := fmt.Sprintf("job_import_%d", time.Now().UnixNano())
	scriptPath := "cmd/warc_importer/main.py"

	args := []string{"--warc", warcPath}
	if offline {
		args = append(args, "--offline")
	}

	req := &models.JobRequest{
		ID:         jobID,
		Type:       models.JobTypeImportManual,
		WARCPath:   warcPath,
		Offline:    offline,
		ScriptPath: scriptPath,
		Args:       args,
		CreatedAt:  time.Now(),
	}

	return j.EnqueueJob(req)
}

// EnqueueJob は任意のジョブリクエストをキューに追加します
func (j *JobOrchestrator) EnqueueJob(req *models.JobRequest) (*models.JobProgress, error) {
	if req.ID == "" {
		req.ID = fmt.Sprintf("job_%d", time.Now().UnixNano())
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}

	progress := &models.JobProgress{
		ID:         req.ID,
		Type:       req.Type,
		Status:     models.JobStatusPending,
		Current:    0,
		Total:      req.Limit,
		Percentage: 0,
		Message:    "Queued (Waiting for worker)",
		Logs:       make([]string, 0),
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

// worker は単一ワーカーとしてキューから順次ジョブを取り出して実行します
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

// executeJob はサブプロセスを起動し、stdout をスキャンして完了まで監視します
func (j *JobOrchestrator) executeJob(req *models.JobRequest) {
	// ジョブ実行のキャンセル用 Context を作成
	jobCtx, jobCancel := context.WithCancel(j.ctx)

	j.mu.Lock()
	progress, exists := j.jobs[req.ID]
	if !exists {
		progress = &models.JobProgress{
			ID:     req.ID,
			Type:   req.Type,
			Logs:   make([]string, 0),
		}
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

	// コマンドの組み立て
	var cmd *exec.Cmd
	if req.ScriptPath != "" {
		cmdArgs := append([]string{req.ScriptPath}, req.Args...)
		cmd = exec.CommandContext(jobCtx, pythonBin, cmdArgs...)
	} else if len(req.Args) > 0 {
		cmd = exec.CommandContext(jobCtx, req.Args[0], req.Args[1:]...)
	} else {
		j.finishJob(req.ID, models.JobStatusFailed, "Invalid job command: empty script and args", nil)
		return
	}

	// 環境変数の注入
	if len(req.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range req.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// パイプ接続 (Stdout / Stderr)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, fmt.Sprintf("Failed to get stdout pipe: %v", err), err)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, fmt.Sprintf("Failed to get stderr pipe: %v", err), err)
		return
	}

	if err := cmd.Start(); err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, fmt.Sprintf("Subprocess start failed: %v", err), err)
		return
	}

	j.mu.Lock()
	j.activeCmd = cmd
	j.mu.Unlock()

	log.Printf("[JobOrchestrator] Started job %s (PID: %d): %s %v", req.ID, cmd.Process.Pid, pythonBin, req.ScriptPath)

	// Stdout / Stderr の非同期スキャン
	go j.scanStdoutProgress(req.ID, stdoutPipe)
	go j.scanStderr(req.ID, stderrPipe)

	err = cmd.Wait()

	if jobCtx.Err() == context.Canceled {
		j.finishJob(req.ID, models.JobStatusCancelled, "Job cancelled by user", nil)
	} else if err != nil {
		j.finishJob(req.ID, models.JobStatusFailed, fmt.Sprintf("Subprocess exited with error: %v", err), err)
	} else {
		j.finishJob(req.ID, models.JobStatusCompleted, "Completed successfully", nil)
	}
}

// scanStdoutProgress は stdout から PROGRESS: {cur}/{total} | {msg} をスキャンして進捗率をオンメモリ更新します
func (j *JobOrchestrator) scanStdoutProgress(jobID string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		j.appendLog(jobID, line)

		// PROGRESS: (\d+)/(\d+) | (.*) の判定
		matches := progressRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			cur, _ := strconv.Atoi(matches[1])
			tot, _ := strconv.Atoi(matches[2])
			msg := matches[3]

			var pct float64
			if tot > 0 {
				pct = (float64(cur) / float64(tot)) * 100.0
				if pct > 100.0 {
					pct = 100.0
				}
			}

			j.mu.Lock()
			if p, ok := j.jobs[jobID]; ok {
				p.Current = cur
				p.Total = tot
				p.Percentage = pct
				p.Message = msg
				// コピーを作成してイベント送信
				snapshot := *p
				j.mu.Unlock()
				j.emitEvent("job:progress", &snapshot)
			} else {
				j.mu.Unlock()
			}
		}
	}
}

// scanStderr は stderr 出力をログに記録します
func (j *JobOrchestrator) scanStderr(jobID string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		j.appendLog(jobID, "[STDERR] "+line)
	}
}

// appendLog はオンメモリログバッファに行を追加します
func (j *JobOrchestrator) appendLog(jobID, line string) {
	j.mu.Lock()
	if p, ok := j.jobs[jobID]; ok {
		p.Logs = append(p.Logs, line)
		if len(p.Logs) > j.maxLogs {
			p.Logs = p.Logs[len(p.Logs)-j.maxLogs:]
		}
	}
	j.mu.Unlock()
	j.emitEvent("job:log", map[string]string{"id": jobID, "line": line})
}

// finishJob はジョブの完了ステータスを確定し、アクティブ状態を解除します
func (j *JobOrchestrator) finishJob(jobID string, status models.JobStatus, message string, err error) {
	now := time.Now()
	j.mu.Lock()
	var snapshot *models.JobProgress
	if p, ok := j.jobs[jobID]; ok {
		p.Status = status
		p.Message = message
		p.FinishedAt = &now
		if err != nil {
			p.Error = err.Error()
		}
		if status == models.JobStatusCompleted {
			p.Percentage = 100.0
			if p.Total > 0 && p.Current < p.Total {
				p.Current = p.Total
			}
		}
		s := *p
		snapshot = &s
	}
	if j.activeJobID == jobID {
		j.activeJobID = ""
		j.activeCancel = nil
		j.activeCmd = nil
	}
	j.mu.Unlock()

	log.Printf("[JobOrchestrator] Finished job %s: status=%s, msg=%s", jobID, status, message)
	if snapshot != nil {
		j.emitEvent("job:finished", snapshot)
		j.emitEvent("job:progress", snapshot)
	}
}

// GetStatus は指定されたジョブの進捗ステータスをオンメモリから即座に取得します
func (j *JobOrchestrator) GetStatus(jobID string) *models.JobProgress {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if p, ok := j.jobs[jobID]; ok {
		s := *p
		return &s
	}
	return nil
}

// GetActiveJob は現在実行中のジョブステータスを取得します
func (j *JobOrchestrator) GetActiveJob() *models.JobProgress {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if j.activeJobID == "" {
		return nil
	}
	if p, ok := j.jobs[j.activeJobID]; ok {
		s := *p
		return &s
	}
	return nil
}

// ListJobs は全ジョブの履歴ステータス一覧を取得します
func (j *JobOrchestrator) ListJobs() []*models.JobProgress {
	j.mu.RLock()
	defer j.mu.RUnlock()
	list := make([]*models.JobProgress, 0, len(j.jobs))
	for _, p := range j.jobs {
		s := *p
		list = append(list, &s)
	}
	return list
}

// CancelJob は指定されたジョブの実行を中断します
func (j *JobOrchestrator) CancelJob(jobID string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	p, ok := j.jobs[jobID]
	if !ok {
		return fmt.Errorf("job not found: %s", jobID)
	}

	if p.Status == models.JobStatusRunning && j.activeJobID == jobID {
		if j.activeCancel != nil {
			j.activeCancel()
		}
		if j.activeCmd != nil && j.activeCmd.Process != nil {
			_ = j.activeCmd.Process.Kill()
		}
		p.Status = models.JobStatusCancelled
		p.Message = "Cancelled by user"
		now := time.Now()
		p.FinishedAt = &now
		snapshot := *p
		go func() {
			j.emitEvent("job:finished", &snapshot)
			j.emitEvent("job:progress", &snapshot)
		}()
		return nil
	}

	if p.Status == models.JobStatusPending {
		p.Status = models.JobStatusCancelled
		p.Message = "Cancelled before execution"
		now := time.Now()
		p.FinishedAt = &now
		snapshot := *p
		go func() {
			j.emitEvent("job:finished", &snapshot)
			j.emitEvent("job:progress", &snapshot)
		}()
		return nil
	}

	return fmt.Errorf("job is not running or pending (current status: %s)", p.Status)
}

// Close は Orchestrator を停止し、全リソースをクリーンアップします
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

func (j *JobOrchestrator) emitEvent(eventName string, data interface{}) {
	if j.emitter != nil {
		j.emitter(eventName, data)
	}
}
