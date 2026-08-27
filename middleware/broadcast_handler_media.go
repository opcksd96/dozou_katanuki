// middleware/broadcast_handler_media.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
)

func (s *BroadcastService) handleMediaBookmarkAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct{ MediaID string `json:"media_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MediaID == "" {
		http.Error(w, `{"error":"media_id required"}`, 400); return
	}
	res, err := s.timelineService.ToggleMediaBookmark(req.MediaID)
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]any{"is_bookmarked": res})
}

func (s *BroadcastService) handleMediaUpdateAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct {
		MediaID string `json:"media_id"`; DownloadStatus string `json:"download_status"`
		StashSceneID string `json:"stash_scene_id"`; StashImageID string `json:"stash_image_id"`; FailedReason string `json:"failed_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MediaID == "" {
		http.Error(w, `{"error":"media_id required"}`, 400); return
	}
	err := s.timelineService.UpdateMediaMetadata(req.MediaID, req.DownloadStatus, req.StashSceneID, req.StashImageID, req.FailedReason)
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *BroadcastService) handleMediaPurgeAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct{ MediaID string `json:"media_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MediaID == "" {
		http.Error(w, `{"error":"media_id required"}`, 400); return
	}
	if err := s.timelineService.PurgeMedia(req.MediaID); err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *BroadcastService) handleMediaPurgeByStatusAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct{ Status string `json:"status"`; AccountID string `json:"account_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		http.Error(w, `{"error":"status required"}`, 400); return
	}
	count, err := s.timelineService.PurgeMediaByStatus(req.Status, req.AccountID)
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]any{"purged_count": count})
}

func (s *BroadcastService) handleMediaRequeueAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct{ Status string `json:"status"`; AccountID string `json:"account_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, `{"error":"invalid json"}`, 400); return }
	count, err := s.timelineService.RequeueMediaByStatus(req.Status, req.AccountID)
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]any{"requeued_count": count})
}

func (s *BroadcastService) handleMediaOpenActionAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Service unavailable"}`, 503); return }
	var req struct{ MediaID string `json:"media_id"`; Action string `json:"action"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MediaID == "" {
		http.Error(w, `{"error":"media_id required"}`, 400); return
	}
	p, err := s.timelineService.ResolveMediaFilePath(req.MediaID)
	if err != nil { http.Error(w, err.Error(), 404); return }
	absPath, _ := filepath.Abs(p)
	if req.Action == "explorer" {
		_ = exec.Command("explorer", "/select,"+absPath).Start()
	} else if req.Action == "default" {
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", absPath).Start()
	}
	w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(map[string]any{"success": true, "path": absPath})
}
