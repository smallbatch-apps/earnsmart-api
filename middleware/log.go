package middleware

import (
	"log"
	"net/http"
	"time"
)

// type wrappedWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func (w *wrappedWriter) WriteHeader(statusCode int) {
// 	w.ResponseWriter.WriteHeader(statusCode)
// 	w.statusCode = statusCode
// }

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		start := time.Now()

		// wrapped := &wrappedWriter{
		// 	ResponseWriter: w,
		// 	statusCode:     http.StatusOK,
		// }

		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}
