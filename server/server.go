package server

import (
	"log/slog"
	"net/http"

	"github.com/kartverket/gcp-sts-proxy/token"
)

type Config struct {
	Port string
}

type Server struct {
	config Config

	handler       http.Handler
	tokenProvider *token.Provider
}

func New(config Config, tokenProvider *token.Provider) *Server {
	server := &Server{
		config:        config,
		tokenProvider: tokenProvider,
	}

	server.initRoutes()

	return server
}

func (server *Server) initRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.proxyHandler)

	server.handler = mux
}

func (s *Server) Start() error {
	slog.Info("starting server", "port", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, s.handler)
}
