package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/vt887/macnet-gateway/daemon/internal/db"
	"github.com/vt887/macnet-gateway/daemon/internal/events"
	"github.com/vt887/macnet-gateway/daemon/internal/models"
)

type Server struct {
	db       *sql.DB
	eventBus events.Bus
}

func NewServer(dbConn *sql.DB, eventBus events.Bus) *Server {
	return &Server{
		db:       dbConn,
		eventBus: eventBus,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/dashboard", s.handleDashboard)
	mux.HandleFunc("GET /api/services", s.handleServices)
	mux.HandleFunc("GET /api/live-activity", s.handleLiveActivity)
	mux.HandleFunc("GET /api/settings", s.handleGetSettings)
	mux.HandleFunc("PATCH /api/settings", s.handlePatchSettings)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, models.Health{Status: "ok"})
}

func (s *Server) handleDashboard(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, models.Dashboard{
		ActiveDevices:          8,
		ProxyRequestsPerMinute: 143,
		DNSQueriesPerMinute:    322,
		BlockedRequests:        27,
		TrafficTodayMB:         1248,
		CacheHitRatio:          0.63,
	})
}

func (s *Server) handleServices(w http.ResponseWriter, _ *http.Request) {
	serviceStatuses, err := db.ListServiceStatuses(context.Background(), s.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, serviceStatuses)
}

func (s *Server) handleLiveActivity(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.eventBus.Recent())
}

func (s *Server) handleGetSettings(w http.ResponseWriter, _ *http.Request) {
	value, err := db.GetSetting(context.Background(), s.db, "ui.theme")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if value == "" {
		value = "system"
	}
	writeJSON(w, http.StatusOK, map[string]string{"ui.theme": value})
}

func (s *Server) handlePatchSettings(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for key, value := range body {
		if err := db.UpsertSetting(context.Background(), s.db, key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
