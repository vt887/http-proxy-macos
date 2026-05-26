package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vt887/macnet-gateway/daemon/internal/api"
	"github.com/vt887/macnet-gateway/daemon/internal/config"
	"github.com/vt887/macnet-gateway/daemon/internal/db"
	"github.com/vt887/macnet-gateway/daemon/internal/events"
	"github.com/vt887/macnet-gateway/daemon/internal/services"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer store.Close()

	if err := db.Initialize(context.Background(), store); err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	if err := db.InsertSettingIfMissing(context.Background(), store, "ui.theme", "system"); err != nil {
		log.Fatalf("failed to seed settings: %v", err)
	}
	for _, status := range services.NewMockRegistry() {
		if err := db.UpsertServiceStatus(context.Background(), store, db.ServiceStatusRecord{
			Name:    status.Name,
			Status:  status.Status,
			Message: status.Message,
		}); err != nil {
			log.Fatalf("failed to seed service status: %v", err)
		}
	}

	eventBus := events.NewMockBus()
	server := api.NewServer(store, eventBus)

	httpServer := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: server.Routes(),
	}

	go func() {
		log.Printf("macnet-gatewayd listening on %s", cfg.ListenAddress)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
}
