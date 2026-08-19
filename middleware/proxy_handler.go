// middleware/proxy_handler.go (100行以下)
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type UnifiedHandler struct {
	avatarDir  string
	stashProxy *httputil.ReverseProxy
	jobOrch    *JobOrchestrator
}

func NewUnifiedHandler(avatarDir string, stashURL *url.URL) *UnifiedHandler {
	_ = os.MkdirAll(avatarDir, 0755)

	var proxy *httputil.ReverseProxy
	if stashURL != nil {
		proxy = httputil.NewSingleHostReverseProxy(stashURL)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = stashURL.Host
		}
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			resp.Header.Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST")
			return nil
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "image/svg+xml")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250" viewBox="0 0 400 250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#ef4444" font-family="sans-serif" font-size="14">Stash Media Offline (502)</text></svg>`))
		}
	}

	return &UnifiedHandler{avatarDir: avatarDir, stashProxy: proxy}
}

// SetJobOrchestrator sets the JobOrchestrator instance for the unified handler
func (h *UnifiedHandler) SetJobOrchestrator(orch *JobOrchestrator) {
	h.jobOrch = orch
}

func (h *UnifiedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// 0. CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. Job Orchestrator API (/api/jobs/...)
	if strings.HasPrefix(path, "/api/jobs/") {
		h.serveJobAPI(w, r)
		return
	}

	// 2. アバター・ローカルアセット解決 (/avatars/ または /assets/)
	if strings.HasPrefix(path, "/avatars/") || strings.HasPrefix(path, "/assets/") {
		rel := strings.TrimPrefix(path, "/avatars/")
		if strings.HasPrefix(path, "/assets/") {
			rel = strings.TrimPrefix(path, "/assets/")
		}
		filePath := filepath.Join(h.avatarDir, filepath.Clean(rel))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}
		h.serveDefaultAvatar(w)
		return
	}

	// 3. Stash サーバーへのインメモリ・リバースプロキシ (/stash-proxy/...)
	if strings.HasPrefix(path, "/stash-proxy/") {
		if h.stashProxy != nil {
			r.URL.Path = strings.TrimPrefix(path, "/stash-proxy")
			if r.URL.RawPath != "" {
				r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "/stash-proxy")
			}
			h.stashProxy.ServeHTTP(w, r)
			return
		}
		h.serveMediaPlaceholder(w)
		return
	}

	http.NotFound(w, r)
}

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
			http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		progress, err := h.jobOrch.EnqueueSalvage(body.Platform, body.Account, body.Limit)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": progress})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(progress)

	case path == "/api/jobs/import-manual" && r.Method == http.MethodPost:
		var body struct {
			WARCPath string `json:"warc_path"`
			Offline  bool   `json:"offline"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		progress, err := h.jobOrch.EnqueueManualImport(body.WARCPath, body.Offline)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "job": progress})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(progress)

	case path == "/api/jobs/status" && r.Method == http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			// ID 未指定時はアクティブジョブまたは全件一覧
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

func (h *UnifiedHandler) serveDefaultAvatar(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64"><rect width="100%" height="100%" fill="#334155" rx="32"/><circle cx="32" cy="24" r="12" fill="#94a3b8"/><path d="M16 54c0-8.837 7.163-16 16-16s16 7.163 16 16" fill="#94a3b8"/></svg>`))
}

func (h *UnifiedHandler) serveMediaPlaceholder(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250" viewBox="0 0 400 250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#64748b" font-family="sans-serif" font-size="14">Attached Media Preview</text></svg>`))
}
