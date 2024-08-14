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

func pingHandler(w http.ResponseWriter, r *http.Request){
    pathParams := r.PathValue("name")
    w.Write([]byte("pong" + "_" + pathParams))
}


func (s *APIServer) Run() error {

    l := logger.Get()

    router := http.NewServeMux()

    router.HandleFunc("GET /users/{uid}", func(w http.ResponseWriter, r *http.Request) {
        userID := r.PathValue("uid")
        w.Write([]byte("User ID:" + userID))
    })
    router.HandleFunc("POST /ping/{name}", pingHandler)

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



