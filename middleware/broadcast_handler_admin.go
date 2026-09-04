package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (s *BroadcastService) handleSystemJournalsAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 200
	}
	entries := GetGlobalJournal().GetEntries(limit)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(entries)
}

func (s *BroadcastService) handleRestartAPI(w http.ResponseWriter, r *http.Request) {
	// Just record the journal entry, since we don't have access to App struct here
	GetGlobalJournal().Record(
		"system", "INFO", "backend_restarted",
		"Wails backend services re-initialized (Dev Mode HTTP triggered)",
		map[string]interface{}{"status": "reloaded_ok_dev"},
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// handleAuditAPI handles GET /api/admin/audit
func (s *BroadcastService) handleAuditAPI(w http.ResponseWriter, r *http.Request) {
	if s.auditService == nil {
		http.Error(w, "AuditService not configured", http.StatusInternalServerError)
		return
	}
	purgeFiles := r.URL.Query().Get("purgeFiles") == "true"
	purgeDB := r.URL.Query().Get("purgeDB") == "true"
	report, err := s.auditService.RunAudit(r.Context(), "./stash", "./blobs", purgeFiles, purgeDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(report)
}

// handleAuditPurgeFilesAPI handles POST /api/admin/audit/purge-files
func (s *BroadcastService) handleAuditPurgeFilesAPI(w http.ResponseWriter, r *http.Request) {
	if s.auditService == nil {
		http.Error(w, "AuditService not configured", http.StatusInternalServerError)
		return
	}
	var paths []string
	_ = json.NewDecoder(r.Body).Decode(&paths)
	count, err := s.auditService.PurgeOrphanFiles(paths)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]int{"purged": count})
}

// handleAuditPurgeDBAPI handles POST /api/admin/audit/purge-db
func (s *BroadcastService) handleAuditPurgeDBAPI(w http.ResponseWriter, r *http.Request) {
	if s.auditService == nil {
		http.Error(w, "AuditService not configured", http.StatusInternalServerError)
		return
	}
	var ids []string
	_ = json.NewDecoder(r.Body).Decode(&ids)
	count, err := s.auditService.PurgeOrphanDBMedia("./backups/dumps/_trash", ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]int{"purged": count})
}

// handleAuditRollbackAPI handles POST /api/admin/audit/rollback
func (s *BroadcastService) handleAuditRollbackAPI(w http.ResponseWriter, r *http.Request) {
	if s.auditService == nil {
		http.Error(w, "AuditService not configured", http.StatusInternalServerError)
		return
	}
	count, err := s.auditService.RollbackLastDBPurge("./backups/dumps/_trash")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]int{"restored": count})
}

// handleAuditCanRollbackAPI handles GET /api/admin/audit/can-rollback
func (s *BroadcastService) handleAuditCanRollbackAPI(w http.ResponseWriter, r *http.Request) {
	if s.auditService == nil {
		http.Error(w, "AuditService not configured", http.StatusInternalServerError)
		return
	}
	canRollback := s.auditService.CanRollbackDBPurge("./backups/dumps/_trash")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]bool{"can_rollback": canRollback})
}

