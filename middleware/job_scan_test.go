// middleware/job_scan_test.go (100行以下)
package middleware

import (
	"strings"
	"sync"
	"testing"

	"dozou_katanuki/models"
)

func TestJobOrchestrator_ProgressScan(t *testing.T) {
	orch := &JobOrchestrator{jobs: make(map[string]*models.JobProgress), maxLogs: 100}
	jobID := "test_job_1"
	orch.jobs[jobID] = &models.JobProgress{ID: jobID, Logs: make([]string, 0)}

	var emittedEvents []string
	var mu sync.Mutex
	orch.emitter = func(event string, data ...interface{}) {
		mu.Lock(); defer mu.Unlock()
		emittedEvents = append(emittedEvents, event)
	}

	dummy := "Starting task...\nPROGRESS: 1/10 | Initializing\nPROGRESS: 5/10 | Halfway\nPROGRESS: 10/10 | Done\n"
	orch.scanStdoutProgress(jobID, strings.NewReader(dummy))

	st := orch.GetStatus(jobID)
	if st == nil || st.Current != 10 || st.Percentage != 100.0 || st.Message != "Done" {
		t.Fatalf("Progress scan failed: %+v", st)
	}
}
