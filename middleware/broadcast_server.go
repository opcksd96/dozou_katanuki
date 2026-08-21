// middleware/broadcast_server.go (100行以下)
package middleware

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func (s *BroadcastService) startServerLocked() error {
	addr := fmt.Sprintf("%s:%d", s.netCfg.PublicBindAddress, s.netCfg.MiddlewarePort)
	rawListener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[Broadcast] Listen failed on %s: %v", addr, err)
		return err
	}

	var listener net.Listener = rawListener
	localIPs := GetLocalIPv4s()
	if cert, err := GenerateSelfSignedCert(localIPs); err == nil {
		listener = tls.NewListener(rawListener, &tls.Config{Certificates: []tls.Certificate{cert}})
		s.useTLS = true
		fmt.Printf("\n[Broadcast] 📡 LAN HTTPS 配信サーバー開通 (0.0.0.0:%d)\n", s.netCfg.MiddlewarePort)
		for _, ip := range localIPs {
			fmt.Printf("[Broadcast] 👉 https://%s:%d/\n", ip, s.netCfg.MiddlewarePort)
		}
	} else {
		s.useTLS = false
		log.Printf("[Broadcast TLS] Cert gen error (%v), falling back to HTTP", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/timeline", s.handleTimelineAPI)
	mux.HandleFunc("/api/search", s.handleSearchAPI)
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
		_ = server.Serve(listener)
	}()
	return nil
}

func (s *BroadcastService) handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(s.GetStatus())
}

func (s *BroadcastService) handleTimelineAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable); return
	}
	q := r.URL.Query()
	p, a, f := q.Get("platform"), q.Get("account_id"), q.Get("filter")
	if p == "" { p = "twitter" }
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	trees, err := s.timelineService.FetchTimeline(p, a, f, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError); _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trees)
}

func (s *BroadcastService) handleSearchAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable); return
	}
	q := r.URL.Query()
	queryText, a, f := q.Get("q"), q.Get("account_id"), q.Get("filter")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	res, err := s.timelineService.SearchArticles(queryText, a, f, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError); _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *BroadcastService) handleArticleAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable); return
	}
	p, id := r.URL.Query().Get("platform"), r.URL.Query().Get("id")
	if p == "" { p = "twitter" }
	if id == "" { http.Error(w, `{"error":"Article id query parameter is required"}`, http.StatusBadRequest); return }

	detail, err := s.timelineService.GetArticleDetail(p, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError); _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *BroadcastService) handleAccountsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable); return
	}
	p := r.URL.Query().Get("platform")
	if p == "" { p = "twitter" }
	accounts, err := s.timelineService.GetAccounts(p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError); _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(accounts)
}
