package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowed := os.Getenv("CORS_ORIGIN")
		if allowed == "" {
			allowed = "http://localhost:5173,http://localhost:80,http://localhost"
		}

		allowedList := strings.Split(allowed, ",")
		for _, a := range allowedList {
			if strings.TrimSpace(a) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		// Fallback for dev: allow any localhost
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			if strings.Contains(origin, "localhost") || origin == "" {
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT, PATCH")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
