package middleware

import (
	"encoding/json"
	"net/http"
	"dozou_katanuki/domain/ports"
)

func (s *BroadcastService) handleResetAllAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed); return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError); return
	}
	res, err := s.adminUseCases.ResetAllToQueuedAndBootstrap()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *BroadcastService) handleIgniteAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed); return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError); return
	}
	res, err := s.adminUseCases.IgnitePipeline()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}

// handlePipelineToggleAPI handles POST /api/admin/pipeline/toggle
func (s *BroadcastService) handlePipelineToggleAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Enable bool `json:"enable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError)
		return
	}

	isRunning, err := s.adminUseCases.TogglePipelineAutoEngine(req.Enable)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]bool{"isRunning": isRunning})
}

func (s *BroadcastService) handlePipelineOverviewAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError); return
	}
	overview, err := s.adminUseCases.GetPipelineOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(overview)
}

func (s *BroadcastService) handlePipelineLogsAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError); return
	}
	stage := r.URL.Query().Get("stage")
	limit := 50
	logs, err := s.adminUseCases.GetPipelineLogs(stage, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(logs)
}

func (s *BroadcastService) handleSyncThunderAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed); return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError); return
	}
	res, err := s.adminUseCases.SyncThunderDownloads("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}
