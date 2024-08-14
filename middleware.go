package main

import (
    "net/http"
    "time"
    "big-john/logger"
)

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.ResponseWriter.WriteHeader(code)
    lrw.statusCode = code
}


func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        log := logger.Get()
        lrw := newLoggingResponseWriter(w)

        defer func() {
            panicVal := recover()
            if panicVal != nil {
                lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
                panic(panicVal) // continue panicking
            }
            log.
            Info().
            Str("method", r.Method).
            Str("url", r.URL.RequestURI()).
            Str("user_agent", r.UserAgent()).
            Dur("elapsed_ms", time.Since(start)).
            Int("status_code", lrw.statusCode).
            Msg("incoming request")
        }()
        next.ServeHTTP(lrw, r)
    }

}


func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "Bearer token" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    }
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
    return func(next http.Handler) http.HandlerFunc{
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next.ServeHTTP
    }
}
