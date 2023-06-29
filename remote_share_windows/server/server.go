package server

import (
	"fmt"
	"log"
	"net/http"

	"remote_share_windows/config"
	"remote_share_windows/lib/mux"
	"remote_share_windows/server/handler"
)

type Server struct {
	httpRouter *mux.Http
	c          *config.AppConfig
}

func NewServer(appConfig *config.AppConfig) *Server {
	r := mux.NewHttp(mux.WithAppConfig(&config.AppConfig{}))

	r.Get("/ws", handler.ApiWsGuacamole)

	return &Server{httpRouter: r, c: appConfig}
}

func (s *Server) Run() {
	fmt.Printf("server run on %s\n", s.c.Addr)
	log.Fatal(http.ListenAndServe(s.c.Addr, s.httpRouter))
}
