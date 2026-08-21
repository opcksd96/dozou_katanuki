// middleware/broadcast_server.go (100行以下)
package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

func (s *BroadcastService) startServerLocked() error {
	addr := fmt.Sprintf("%s:%d", s.netCfg.PublicBindAddress, s.netCfg.MiddlewarePort)
	listener, err := net.Listen("tcp", addr)
	if err != nil { return err }

	mux := http.NewServeMux()
	mux.HandleFunc("/api/timeline", s.handleTimelineAPI)
	mux.HandleFunc("/api/article", s.handleArticleAPI)
	mux.HandleFunc("/api/accounts", s.handleAccountsAPI)
	mux.HandleFunc("/api/broadcast/status", s.handleStatusAPI)
	mux.HandleFunc("/", s.handleRoot)

	server := &http.Server{
		Handler: s.securityMiddleware(s.corsMiddleware(mux)),
		ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second,
	}
	s.server, s.listener, s.running = server, listener, true

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.mu.Lock(); s.running = false; s.mu.Unlock()
		}
	}()
	return nil
}

func (s *BroadcastService) handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(s.GetStatus())
}

func (s *BroadcastService) handleTimelineAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}
	q := r.URL.Query()
	p, a, f := q.Get("platform"), q.Get("account_id"), q.Get("filter")
	if p == "" { p = "twitter" }
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	trees, err := s.timelineService.FetchTimeline(p, a, f, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trees)
}

func (s *BroadcastService) handleArticleAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}
	p, id := r.URL.Query().Get("platform"), r.URL.Query().Get("id")
	if p == "" { p = "twitter" }
	if id == "" { http.Error(w, `{"error":"Article id query parameter is required"}`, http.StatusBadRequest); return }

	detail, err := s.timelineService.GetArticleDetail(p, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *BroadcastService) handleAccountsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}
	p := r.URL.Query().Get("platform")
	if p == "" { p = "twitter" }
	accounts, err := s.timelineService.GetAccounts(p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(accounts)
}
