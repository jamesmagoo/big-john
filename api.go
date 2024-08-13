package main

import (
	"log"
	"net/http"
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

    router := http.NewServeMux()
    router.HandleFunc("GET /users/{uid}", func(w http.ResponseWriter, r *http.Request) {
        userID := r.PathValue("uid")
        w.Write([]byte("User ID:" + userID))
    })

    server := http.Server{
        Addr: s.addr,
        Handler: router,
    }

    log.Printf("Serving on PORT %v", s.addr)

    return server.ListenAndServe()
}
