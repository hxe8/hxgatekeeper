package main

import (
	"log"

	"github.com/hex/gatekeeper-go/internal/audit"
	"github.com/hex/gatekeeper-go/internal/config"
	"github.com/hex/gatekeeper-go/internal/httpserver"
	"github.com/hex/gatekeeper-go/internal/policy"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	engine, err := policy.LoadFromFile(cfg.PolicyFile, cfg.DefaultMode)
	if err != nil {
		log.Fatalf("load policy engine: %v", err)
	}

	auditStore := audit.NewStore(cfg.AuditCapacity)
	server := httpserver.New(cfg.BindAddress, engine, auditStore)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
