package http

import (
	"log/slog"
	"net/http"
	"time"
)

func WriteAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		status := rw.status
		if status == 0 {
			status = http.StatusOK
		}

		slog.InfoContext(
			r.Context(),
			"http request",
			"http.request.method", r.Method,
			"url.path", r.URL.Path,
			"http.response.status_code", status,
			"duration", time.Since(start),
		)
	})
}
