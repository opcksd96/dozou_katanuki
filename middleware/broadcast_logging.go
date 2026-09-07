// middleware/broadcast_logging.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *BroadcastService) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if shouldSkipJournalLog(path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		clientIP := s.extractClientIP(r)
		level := "INFO"
		if lrw.statusCode >= 500 {
			level = "ERROR"
		} else if lrw.statusCode >= 400 {
			level = "WARN"
		}

		msg := fmt.Sprintf("%s %s -> %d (%v, %s)", r.Method, path, lrw.statusCode, duration.Round(time.Millisecond), clientIP)
		payload := map[string]interface{}{
			"method":   r.Method,
			"path":     path,
			"status":   lrw.statusCode,
			"duration": duration.String(),
			"client_ip": clientIP,
		}

		GetGlobalJournal().Record("system", level, "http_response", msg, payload)
	})
}

func shouldSkipJournalLog(path string) bool {
	if path == "/api/admin/system/journals" {
		return true
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".js", ".mjs", ".css", ".wasm", ".woff2", ".svg", ".png", ".jpg", ".jpeg", ".ico":
		return true
	}
	return strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/avatars/")
}
