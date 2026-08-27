// middleware/broadcast_handler.go (100行以下)
package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *BroadcastService) handleAccountsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	accounts, err := s.timelineService.ListRawAccounts()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(accounts)
}

func (s *BroadcastService) handleAccountDetailAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	id := r.URL.Query().Get("id")
	if id == "" { http.Error(w, `{"error":"id required"}`, 400); return }
	detail, err := s.timelineService.GetAccountDetail(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *BroadcastService) handleMediaAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit")); offset, _ := strconv.Atoi(q.Get("offset"))
	res, err := s.timelineService.SearchMediaDetails(q.Get("account_id"), q.Get("status"), q.Get("type"), limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *BroadcastService) handleMediaStatsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	stats, err := s.timelineService.FetchDownloadStatusStats(r.URL.Query().Get("account_id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(stats)
}


