package middleware

import (
	"net/http"
	"os"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminPassword := r.URL.Query().Get("admin_password")
		if adminPassword != os.Getenv("ADMIN_PASSWORD") {
			http.Error(w, "Invalid admin password", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
