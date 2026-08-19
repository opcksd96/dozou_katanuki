package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dozou_katanuki/models"
)

func TestUnifiedHandler_JobAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := NewJobOrchestrator(ctx, nil)
	defer orch.Close()

	handler := NewUnifiedHandler("./assets", nil)
	handler.SetJobOrchestrator(orch)

	// 1. POST /api/jobs/salvage
	salvagePayload := []byte(`{"platform":"twitter","account":"test_user","limit":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/salvage", bytes.NewBuffer(salvagePayload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202 Accepted, got %d (body: %s)", w.Code, w.Body.String())
	}

	var salvageResp models.JobProgress
	if err := json.NewDecoder(w.Body).Decode(&salvageResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if salvageResp.ID == "" {
		t.Errorf("expected job ID to be non-empty")
	}
	if salvageResp.Status != models.JobStatusPending && salvageResp.Status != models.JobStatusRunning {
		t.Errorf("expected pending or running status, got %s", salvageResp.Status)
	}

	// 2. GET /api/jobs/status?id={job_id}
	reqStatus := httptest.NewRequest(http.MethodGet, "/api/jobs/status?id="+salvageResp.ID, nil)
	wStatus := httptest.NewRecorder()
	handler.ServeHTTP(wStatus, reqStatus)

	if wStatus.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", wStatus.Code)
	}

	var statusResp models.JobProgress
	if err := json.NewDecoder(wStatus.Body).Decode(&statusResp); err != nil {
		t.Fatalf("failed to decode status: %v", err)
	}
	if statusResp.ID != salvageResp.ID {
		t.Errorf("expected ID %s, got %s", salvageResp.ID, statusResp.ID)
	}

	// 3. POST /api/jobs/cancel
	cancelPayload := []byte(`{"id":"` + salvageResp.ID + `"}`)
	reqCancel := httptest.NewRequest(http.MethodPost, "/api/jobs/cancel", bytes.NewBuffer(cancelPayload))
	wCancel := httptest.NewRecorder()
	handler.ServeHTTP(wCancel, reqCancel)

	if wCancel.Code != http.StatusOK && wCancel.Code != http.StatusBadRequest {
		t.Fatalf("unexpected cancel status: %d", wCancel.Code)
	}
}
