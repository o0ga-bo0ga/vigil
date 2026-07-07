package api

import (
	"log"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xApiKey := r.Header.Get("x-api-key")

		apiKey := os.Getenv("API_KEY")

		if xApiKey == "" {
			http.Error(w, "Missing api key", http.StatusUnauthorized)
			return
		}
		if apiKey == "" {
			http.Error(w, "Server Misconfigured: API_KEY not set", http.StatusInternalServerError)
			return
		}
		if xApiKey != apiKey {
			http.Error(w, "Invalid api key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}