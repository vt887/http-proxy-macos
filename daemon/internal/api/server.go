package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/vt887/macnet-gateway/daemon/internal/db"
	"github.com/vt887/macnet-gateway/daemon/internal/events"
	"github.com/vt887/macnet-gateway/daemon/internal/models"
	"github.com/vt887/macnet-gateway/daemon/internal/services"
)

type Server struct {
	db       *sql.DB
	eventBus events.Bus
	squidSvc services.SquidService
}

func NewServer(dbConn *sql.DB, eventBus events.Bus, squidSvc services.SquidService) *Server {
	return &Server{
		db:       dbConn,
		eventBus: eventBus,
		squidSvc: squidSvc,
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
	mux.HandleFunc("GET /api/proxy/status", s.handleProxyStatus)
	mux.HandleFunc("GET /api/proxy/settings", s.handleProxySettings)
	mux.HandleFunc("POST /api/proxy/validate", s.handleProxyValidate)
	mux.HandleFunc("POST /api/proxy/reload", s.handleProxyReload)
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

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	serviceStatuses, err := db.ListServiceStatuses(r.Context(), s.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, serviceStatuses)
}

func (s *Server) handleLiveActivity(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.eventBus.Recent())
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	value, err := db.GetSetting(r.Context(), s.db, "ui.theme")
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
		if err := db.UpsertSetting(r.Context(), s.db, key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (s *Server) handleProxyStatus(w http.ResponseWriter, r *http.Request) {
	status, err := s.squidSvc.Status(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleProxySettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.squidSvc.Settings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	preview, err := s.squidSvc.ConfigPreview(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, models.ProxySettings{
		ListenAddress:       settings.ListenAddress,
		CacheDirectory:      settings.CacheDirectory,
		GeneratedConfigPath: settings.GeneratedConfigPath,
		ConfigPreview:       preview,
	})
}

func (s *Server) handleProxyValidate(w http.ResponseWriter, r *http.Request) {
	if err := s.squidSvc.ValidateConfig(r.Context()); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ActionResult{Status: "invalid", Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, models.ActionResult{Status: "valid", Message: "Squid config is valid"})
}

func (s *Server) handleProxyReload(w http.ResponseWriter, r *http.Request) {
	if err := s.squidSvc.Reload(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, models.ActionResult{Status: "reloaded", Message: "Squid reload request accepted"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
