package main

import (
    "net/http"
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

    logger.PrintAsciiArt()

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


