package middleware

import (
    "time"
    "net/http"
    "log"
)

type responseWriterWrapper struct {
    http.ResponseWriter
    statusCode  int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(wrappedWriter, r)
        duration := time.Since(start)

        log.Printf("[HTTP] %s %s %d - %s\n", r.Method, r.RequestURI, wrappedWriter.statusCode, duration)
    })
}
