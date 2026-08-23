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
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable); return
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
