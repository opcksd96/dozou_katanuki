// middleware/proxy_job_api.go (100行以下)
package middleware

import (
	"encoding/json"
	"net/http"
)

func (h *UnifiedHandler) serveJobAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if h.jobOrch == nil {
		http.Error(w, `{"error":"Job orchestrator not initialized"}`, http.StatusServiceUnavailable)
		return
	}

	path := r.URL.Path
	switch {
	case path == "/api/jobs/salvage" && r.Method == http.MethodPost:
		var b struct { Platform, Account string; Limit int }
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil { http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest); return }
		p, err := h.jobOrch.EnqueueSalvage(b.Platform, b.Account, b.Limit)
		if err != nil { w.WriteHeader(http.StatusConflict); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p}); return }
		w.WriteHeader(http.StatusAccepted); _ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/import-manual" && r.Method == http.MethodPost:
		var b struct { WARCPath string `json:"warc_path"`; Offline bool `json:"offline"` }
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil { http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest); return }
		p, err := h.jobOrch.EnqueueManualImport(b.WARCPath, b.Offline)
		if err != nil { w.WriteHeader(http.StatusConflict); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p}); return }
		w.WriteHeader(http.StatusAccepted); _ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/restore" && r.Method == http.MethodPost:
		var b struct { DumpsDir string `json:"dumps_dir"` }
		_ = json.NewDecoder(r.Body).Decode(&b)
		p, err := h.jobOrch.EnqueueRestore(b.DumpsDir)
		if err != nil { w.WriteHeader(http.StatusConflict); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p}); return }
		w.WriteHeader(http.StatusAccepted); _ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/status" && r.Method == http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			if a := h.jobOrch.GetActiveJob(); a != nil { _ = json.NewEncoder(w).Encode(a) } else { _ = json.NewEncoder(w).Encode(h.jobOrch.ListJobs()) }
			return
		}
		st := h.jobOrch.GetStatus(id)
		if st == nil { http.Error(w, `{"error":"Job not found"}`, http.StatusNotFound); return }
		_ = json.NewEncoder(w).Encode(st)

	case path == "/api/jobs/cancel" && r.Method == http.MethodPost:
		var b struct { ID string `json:"id"` }
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.ID == "" { http.Error(w, `{"error":"Job ID is required"}`, http.StatusBadRequest); return }
		if err := h.jobOrch.CancelJob(b.ID); err != nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return }
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Cancel requested"})

	default:
		http.NotFound(w, r)
	}
}
