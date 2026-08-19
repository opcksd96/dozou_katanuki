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
		var body struct {
			Platform string `json:"platform"`
			Account  string `json:"account"`
			Limit    int    `json:"limit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
			return
		}
		p, err := h.jobOrch.EnqueueSalvage(body.Platform, body.Account, body.Limit)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/import-manual" && r.Method == http.MethodPost:
		var body struct {
			WARCPath string `json:"warc_path"`
			Offline  bool   `json:"offline"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
			return
		}
		p, err := h.jobOrch.EnqueueManualImport(body.WARCPath, body.Offline)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/restore" && r.Method == http.MethodPost:
		var body struct {
			DumpsDir string `json:"dumps_dir"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		p, err := h.jobOrch.EnqueueRestore(body.DumpsDir)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": p})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(p)

	case path == "/api/jobs/status" && r.Method == http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			active := h.jobOrch.GetActiveJob()
			if active != nil {
				_ = json.NewEncoder(w).Encode(active)
			} else {
				_ = json.NewEncoder(w).Encode(h.jobOrch.ListJobs())
			}
			return
		}
		status := h.jobOrch.GetStatus(id)
		if status == nil {
			http.Error(w, `{"error":"Job not found"}`, http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(status)

	case path == "/api/jobs/cancel" && r.Method == http.MethodPost:
		var body struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
			http.Error(w, `{"error":"Job ID is required"}`, http.StatusBadRequest)
			return
		}
		if err := h.jobOrch.CancelJob(body.ID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Cancel requested"})

	default:
		http.NotFound(w, r)
	}
}
