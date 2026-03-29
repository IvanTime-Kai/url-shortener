package middleware

import (
	"net"
	"net/http"

	"github.com/IvanTime-Kai/url-shortener/internal/cache"
)


func RateLimit(limiter *cache.RateLimit) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr

            if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
                ip = forwarded
            } else {
                // Tách IP ra khỏi "IP:port"
                host, _, err := net.SplitHostPort(r.RemoteAddr)
                if err == nil {
                    ip = host
                }
            }

            allowed, err := limiter.Allow(r.Context(), ip)
            if err != nil || !allowed {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusTooManyRequests)
                w.Write([]byte(`{"error":"too many requests"}`))
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}