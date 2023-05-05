package util

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func MonkeyPatchMiddleMan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lrw, r)

		end := time.Now()
		latency := end.Sub(start)

		log.Printf("Method: %s URI: %s Status: %d Latency: %s", r.Method, r.RequestURI, lrw.statusCode, latency.String())
	})
}
