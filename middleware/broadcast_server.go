// middleware/broadcast_server.go (100行以下)
package middleware

import (
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
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[Broadcast] Listen failed on %s: %v", addr, err)
		return err
	}

	s.useTLS = false
	localIPs := GetLocalIPv4s()
	fmt.Printf("\n[Broadcast] 📡 LAN HTTP 配信サーバー開通 (0.0.0.0:%d)\n", s.netCfg.MiddlewarePort)
	for _, ip := range localIPs {
		fmt.Printf("[Broadcast] 👉 http://%s:%d/\n", ip, s.netCfg.MiddlewarePort)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/timeline", s.handleTimelineAPI)
	mux.HandleFunc("/api/search", s.handleSearchAPI)
	// Admin endpoints are now protected by ports.IsAdmin check in handlers
	mux.HandleFunc("/api/admin/pipeline/toggle", s.handlePipelineToggleAPI)
	mux.HandleFunc("/api/admin/pipeline/overview", s.handlePipelineOverviewAPI)
	mux.HandleFunc("/api/admin/pipeline/logs", s.handlePipelineLogsAPI)
	mux.HandleFunc("/api/admin/pipeline/sync-thunder", s.handleSyncThunderAPI)
	mux.HandleFunc("/api/admin/pipeline/reset-all", s.handleResetAllAPI)
	mux.HandleFunc("/api/admin/pipeline/ignite", s.handleIgniteAPI)
	mux.HandleFunc("/api/admin/system/journals", s.handleSystemJournalsAPI)
	mux.HandleFunc("/api/admin/system/restart", s.handleRestartAPI)
	mux.HandleFunc("/api/article/trash", s.handleArticleTrashAPI)
	mux.HandleFunc("/api/article/restore", s.handleArticleRestoreAPI)
	mux.HandleFunc("/api/article/batch-trash", s.handleArticleBatchTrashAPI)
	mux.HandleFunc("/api/article/batch-restore", s.handleArticleBatchRestoreAPI)
	mux.HandleFunc("/api/article/batch-reset-translations", s.handleArticleBatchResetTranslationsAPI)
	mux.HandleFunc("/api/article", s.handleArticleAPI)
	mux.HandleFunc("/api/accounts", s.handleAccountsAPI)
	mux.HandleFunc("/api/account/detail", s.handleAccountDetailAPI)
	mux.HandleFunc("/api/media/stats", s.handleMediaStatsAPI)
	mux.HandleFunc("/api/media/bookmark", s.handleMediaBookmarkAPI)
	mux.HandleFunc("/api/media/update", s.handleMediaUpdateAPI)
	mux.HandleFunc("/api/media/purge-status", s.handleMediaPurgeByStatusAPI)
	mux.HandleFunc("/api/media/purge", s.handleMediaPurgeAPI)
	mux.HandleFunc("/api/media/requeue", s.handleMediaRequeueAPI)
	mux.HandleFunc("/api/media/open", s.handleMediaOpenActionAPI)
	mux.HandleFunc("/api/media", s.handleMediaAPI)
	mux.HandleFunc("/api/broadcast/status", s.handleStatusAPI)
	mux.HandleFunc("/api/events", s.handleEventsAPI)
	mux.HandleFunc("/", s.handleRoot)

	server := &http.Server{
		Handler:     s.corsMiddleware(s.securityMiddleware(mux)),
		ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second,
	}
	s.server, s.listener, s.running = server, listener, true

	go func() { _ = server.Serve(listener) }()
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
	if a == "" { a = "all" }
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
