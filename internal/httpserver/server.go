package httpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/hex/gatekeeper-go/internal/audit"
	"github.com/hex/gatekeeper-go/internal/policy"
	"github.com/hex/gatekeeper-go/internal/service/access"
)

type Server struct {
	httpServer *http.Server
}

func New(bindAddress string, engine *policy.Engine, auditStore *audit.Store) *Server {
	service := &access.Service{Engine: engine, Audit: auditStore}
	mux := http.NewServeMux()
	registerRoutes(mux, service, auditStore)

	return &Server{
		httpServer: &http.Server{
			Addr:              bindAddress,
			Handler:           requestLogger(mux),
			ReadHeaderTimeout: 3 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("gatekeeper listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
