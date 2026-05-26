package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vt887/macnet-gateway/daemon/internal/db"
	"github.com/vt887/macnet-gateway/daemon/internal/events"
	"github.com/vt887/macnet-gateway/daemon/internal/services"
	"github.com/vt887/macnet-gateway/daemon/internal/services/squid"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "app.sqlite")
	store, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := db.Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}
	for _, service := range services.NewMockRegistry() {
		err := db.UpsertServiceStatus(context.Background(), store, db.ServiceStatusRecord{
			Name:    service.Name,
			Status:  service.Status,
			Message: service.Message,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	squidSvc := squid.NewMockService(filepath.Join(t.TempDir(), "generated", "squid"))
	if err := squidSvc.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	return NewServer(store, events.NewMockBus(), squidSvc)
}

func TestHealthEndpoint(t *testing.T) {
	server := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestPatchAndGetSettings(t *testing.T) {
	server := testServer(t)
	req := httptest.NewRequest(http.MethodPatch, "/api/settings", strings.NewReader(`{"ui.theme":"dark"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected patch status: %d", rec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	getRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get status: %d", getRec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(getRec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["ui.theme"] != "dark" {
		t.Fatalf("expected dark theme, got %q", payload["ui.theme"])
	}
}

func TestServicesEndpoint(t *testing.T) {
	server := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/services", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var payload []db.ServiceStatusRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload) == 0 {
		t.Fatal("expected seeded service statuses")
	}
}

func TestProxyEndpoints(t *testing.T) {
	server := testServer(t)

	statusReq := httptest.NewRequest(http.MethodGet, "/api/proxy/status", nil)
	statusRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(statusRec, statusReq)
	if statusRec.Code != http.StatusOK {
		t.Fatalf("unexpected status endpoint code: %d", statusRec.Code)
	}

	settingsReq := httptest.NewRequest(http.MethodGet, "/api/proxy/settings", nil)
	settingsRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(settingsRec, settingsReq)
	if settingsRec.Code != http.StatusOK {
		t.Fatalf("unexpected settings endpoint code: %d", settingsRec.Code)
	}

	validateReq := httptest.NewRequest(http.MethodPost, "/api/proxy/validate", nil)
	validateRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(validateRec, validateReq)
	if validateRec.Code != http.StatusOK {
		t.Fatalf("unexpected validate endpoint code: %d", validateRec.Code)
	}

	reloadReq := httptest.NewRequest(http.MethodPost, "/api/proxy/reload", nil)
	reloadRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(reloadRec, reloadReq)
	if reloadRec.Code != http.StatusOK {
		t.Fatalf("unexpected reload endpoint code: %d", reloadRec.Code)
	}
}
