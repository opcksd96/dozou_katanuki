// middleware/job_queue_test.go (100行以下)
package middleware

import (
	"context"
	"testing"
	"time"

	"dozou_katanuki/models"
)

func TestJobOrchestrator_QueueExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := NewJobOrchestrator(ctx, func(string, ...interface{}) {})
	defer orch.Close()

	req1 := &models.JobRequest{
		ID: "job_q1", Type: models.JobTypeCustom,
		Args: []string{"python", "-c", "print('PROGRESS: 1/1 | step 1', flush=True)"},
	}
	if _, err := orch.EnqueueJob(req1); err != nil { t.Fatalf("failed to enqueue: %v", err) }

	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout: t.Fatalf("timeout waiting for job_q1")
		default:
			if s1 := orch.GetStatus("job_q1"); s1 != nil && s1.Status == models.JobStatusCompleted { return }
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestJobOrchestrator_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := NewJobOrchestrator(ctx, nil)
	defer orch.Close()

	req := &models.JobRequest{
		ID: "job_cancel_test", Type: models.JobTypeCustom,
		Args: []string{"python", "-c", "import time; time.sleep(10)"},
	}
	if _, err := orch.EnqueueJob(req); err != nil { t.Fatalf("failed to enqueue: %v", err) }

	for i := 0; i < 20; i++ {
		if st := orch.GetStatus("job_cancel_test"); st != nil && st.Status == models.JobStatusRunning { break }
		time.Sleep(50 * time.Millisecond)
	}

	if err := orch.CancelJob("job_cancel_test"); err != nil { t.Fatalf("failed to cancel: %v", err) }
	for i := 0; i < 20; i++ {
		if st := orch.GetStatus("job_cancel_test"); st != nil && st.Status == models.JobStatusCancelled { return }
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("job was not cancelled in time")
}
