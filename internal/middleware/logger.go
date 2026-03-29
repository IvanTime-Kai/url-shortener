package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Wrap ResponseWriter để bắt status code
            wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
            next.ServeHTTP(wrapped, r)

            logger.Info("request",
                "method",  r.Method,
                "path",    r.URL.Path,
                "status",  wrapped.status,
                "latency", time.Since(start).String(),
                "ip",      r.RemoteAddr,
            )
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(status int) {
    rw.status = status
    rw.ResponseWriter.WriteHeader(status)
}