package middleware

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"dozou_katanuki/models"
)

func TestJobOrchestrator_ProgressScan(t *testing.T) {
	orch := &JobOrchestrator{
		jobs:    make(map[string]*models.JobProgress),
		maxLogs: 100,
	}

	jobID := "test_job_1"
	orch.jobs[jobID] = &models.JobProgress{
		ID:   jobID,
		Logs: make([]string, 0),
	}

	var emittedEvents []string
	var mu sync.Mutex
	orch.emitter = func(event string, data ...interface{}) {
		mu.Lock()
		defer mu.Unlock()
		emittedEvents = append(emittedEvents, event)
	}

	dummyOutput := `Starting task...
PROGRESS: 1/10 | Initializing components
PROGRESS: 5/10 | Halfway done
Some debug message
PROGRESS: 10/10 | Finished all items
Done!
`

	orch.scanStdoutProgress(jobID, strings.NewReader(dummyOutput))

	status := orch.GetStatus(jobID)
	if status == nil {
		t.Fatalf("expected status to be non-nil")
	}

	if status.Current != 10 || status.Total != 10 {
		t.Errorf("expected 10/10, got %d/%d", status.Current, status.Total)
	}
	if status.Percentage != 100.0 {
		t.Errorf("expected 100%%, got %f", status.Percentage)
	}
	if status.Message != "Finished all items" {
		t.Errorf("expected message 'Finished all items', got '%s'", status.Message)
	}
	if len(status.Logs) != 6 {
		t.Errorf("expected 6 logs, got %d", len(status.Logs))
	}
}

func TestJobOrchestrator_QueueExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var events []string
	var mu sync.Mutex
	emitter := func(event string, data ...interface{}) {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, event)
	}

	orch := NewJobOrchestrator(ctx, emitter)
	defer orch.Close()

	// Python -c でテスト実行
	req1 := &models.JobRequest{
		ID:   "job_q1",
		Type: models.JobTypeCustom,
		Args: []string{"python", "-c", "import sys; print('PROGRESS: 1/2 | step 1', flush=True); print('PROGRESS: 2/2 | step 2', flush=True)"},
	}

	req2 := &models.JobRequest{
		ID:   "job_q2",
		Type: models.JobTypeCustom,
		Args: []string{"python", "-c", "import sys; print('PROGRESS: 5/5 | done', flush=True)"},
	}

	_, err := orch.EnqueueJob(req1)
	if err != nil {
		t.Fatalf("failed to enqueue req1: %v", err)
	}
	_, err = orch.EnqueueJob(req2)
	if err != nil {
		t.Fatalf("failed to enqueue req2: %v", err)
	}

	// 完了待機 (最大5秒)
	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for jobs to complete")
		default:
			s1 := orch.GetStatus("job_q1")
			s2 := orch.GetStatus("job_q2")
			if s1 != nil && s2 != nil &&
				(s1.Status == models.JobStatusCompleted || s1.Status == models.JobStatusFailed) &&
				(s2.Status == models.JobStatusCompleted || s2.Status == models.JobStatusFailed) {
				if s1.Status != models.JobStatusCompleted {
					t.Errorf("job_q1 status = %s, error = %s", s1.Status, s1.Error)
				}
				if s2.Status != models.JobStatusCompleted {
					t.Errorf("job_q2 status = %s, error = %s", s2.Status, s2.Error)
				}
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func TestJobOrchestrator_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := NewJobOrchestrator(ctx, nil)
	defer orch.Close()

	// 長時間かかるジョブ (python -c "import time; time.sleep(10)")
	req := &models.JobRequest{
		ID:   "job_cancel_test",
		Type: models.JobTypeCustom,
		Args: []string{"python", "-c", "import time; time.sleep(10)"},
	}

	_, err := orch.EnqueueJob(req)
	if err != nil {
		t.Fatalf("failed to enqueue: %v", err)
	}

	// 実行中になるのを待機
	for i := 0; i < 20; i++ {
		st := orch.GetStatus("job_cancel_test")
		if st != nil && st.Status == models.JobStatusRunning {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	err = orch.CancelJob("job_cancel_test")
	if err != nil {
		t.Fatalf("failed to cancel job: %v", err)
	}

	// キャンセル完了を待機
	for i := 0; i < 20; i++ {
		st := orch.GetStatus("job_cancel_test")
		if st != nil && st.Status == models.JobStatusCancelled {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	status := orch.GetStatus("job_cancel_test")
	if status.Status != models.JobStatusCancelled {
		t.Errorf("expected status to be cancelled, got %s (err: %s)", status.Status, status.Error)
	}
}
