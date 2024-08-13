package main

import (
	"net/http"
    "time"
    "big-john/logger"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}


func (s *APIServer) Run() error {

    l := logger.Get()


	router := http.NewServeMux()
	router.HandleFunc("GET /users/{uid}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("uid")
		w.Write([]byte("User ID:" + userID))
	})

    v1 := http.NewServeMux()
    v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

    // order matters here...
    middlewareChain := MiddlewareChain(
        RequestLoggerMiddleware,
        RequireAuthMiddleware,
    )

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(v1),
	}

    l.Info().Str("port", server.Addr).Msg("BIG JOHN serving...")

	return server.ListenAndServe()
}


type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}



func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        l := logger.Get()
        lrw := newLoggingResponseWriter(w)

        defer func() {
            panicVal := recover()
            if panicVal != nil {
                lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
                panic(panicVal) // continue panicking
            }
            l.
            Info().
            Str("method", r.Method).
            Str("url", r.URL.RequestURI()).
            Str("user_agent", r.UserAgent()).
            Dur("elapsed_ms", time.Since(start)).
            Int("status_code", lrw.statusCode).

            Msg("incoming request")
        }()
        next.ServeHTTP(w, r)
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

type Middlware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middlware) Middlware {
    return func(next http.Handler) http.HandlerFunc{
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next.ServeHTTP
    }
}
