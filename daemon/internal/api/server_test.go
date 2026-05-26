package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vt887/macnet-gateway/daemon/internal/db"
	"github.com/vt887/macnet-gateway/daemon/internal/events"
	"github.com/vt887/macnet-gateway/daemon/internal/services"
	"github.com/vt887/macnet-gateway/daemon/internal/services/dns"
	"github.com/vt887/macnet-gateway/daemon/internal/services/squid"
)

type failingValidateSquidService struct{}
type failingValidateDNSService struct{}
type invalidValidateDNSService struct{ failingValidateDNSService }

func (f failingValidateSquidService) Status(context.Context) (services.ServiceStatus, error) {
	return services.ServiceStatus{Name: "squid", Status: "mock", Message: "test"}, nil
}
func (f failingValidateSquidService) Settings(context.Context) (services.ProxySettings, error) {
	return services.ProxySettings{}, nil
}
func (f failingValidateSquidService) ConfigPreview(context.Context) (string, error) { return "", nil }
func (f failingValidateSquidService) RenderConfig(context.Context) error            { return nil }
func (f failingValidateSquidService) ValidateConfig(context.Context) error {
	return errors.New("io failure")
}
func (f failingValidateSquidService) Reload(context.Context) error  { return nil }
func (f failingValidateSquidService) Start(context.Context) error   { return nil }
func (f failingValidateSquidService) Stop(context.Context) error    { return nil }
func (f failingValidateSquidService) Restart(context.Context) error { return nil }
func (f failingValidateSquidService) TailAccessLog(context.Context) (<-chan services.ProxyEvent, error) {
	ch := make(chan services.ProxyEvent)
	close(ch)
	return ch, nil
}

func (f failingValidateDNSService) Status(context.Context) (services.ServiceStatus, error) {
	return services.ServiceStatus{Name: "dns", Status: "mock", Message: "test"}, nil
}
func (f failingValidateDNSService) Settings(context.Context) (services.DNSSettings, error) {
	return services.DNSSettings{}, nil
}
func (f failingValidateDNSService) ConfigPreview(context.Context) (string, error) { return "", nil }
func (f failingValidateDNSService) RenderConfig(context.Context) error            { return nil }
func (f failingValidateDNSService) ValidateConfig(context.Context) error {
	return errors.New("io failure")
}
func (f failingValidateDNSService) Reload(context.Context) error  { return nil }
func (f failingValidateDNSService) Start(context.Context) error   { return nil }
func (f failingValidateDNSService) Stop(context.Context) error    { return nil }
func (f failingValidateDNSService) Restart(context.Context) error { return nil }
func (f failingValidateDNSService) TailQueryLog(context.Context) (<-chan services.DNSEvent, error) {
	ch := make(chan services.DNSEvent)
	close(ch)
	return ch, nil
}

func (f invalidValidateDNSService) ValidateConfig(context.Context) error {
	return errors.Join(dns.ErrInvalidConfig, errors.New("missing listen-address"))
}

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
	dnsSvc := dns.NewMockService(filepath.Join(t.TempDir(), "generated", "dns"))
	if err := dnsSvc.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	return NewServer(store, events.NewMockBus(), squidSvc, dnsSvc)
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
	var statusPayload map[string]string
	if err := json.Unmarshal(statusRec.Body.Bytes(), &statusPayload); err != nil {
		t.Fatal(err)
	}
	if statusPayload["name"] != "squid" {
		t.Fatalf("expected squid status payload, got %#v", statusPayload)
	}

	settingsReq := httptest.NewRequest(http.MethodGet, "/api/proxy/settings", nil)
	settingsRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(settingsRec, settingsReq)
	if settingsRec.Code != http.StatusOK {
		t.Fatalf("unexpected settings endpoint code: %d", settingsRec.Code)
	}
	var settingsPayload map[string]string
	if err := json.Unmarshal(settingsRec.Body.Bytes(), &settingsPayload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(settingsPayload["config_preview"], "http_port 127.0.0.1:3128") {
		t.Fatalf("expected http_port in config preview, got %#v", settingsPayload)
	}

	validateReq := httptest.NewRequest(http.MethodPost, "/api/proxy/validate", nil)
	validateRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(validateRec, validateReq)
	if validateRec.Code != http.StatusOK {
		t.Fatalf("unexpected validate endpoint code: %d", validateRec.Code)
	}
	var validatePayload map[string]string
	if err := json.Unmarshal(validateRec.Body.Bytes(), &validatePayload); err != nil {
		t.Fatal(err)
	}
	if validatePayload["status"] != "valid" {
		t.Fatalf("expected valid status, got %#v", validatePayload)
	}

	reloadReq := httptest.NewRequest(http.MethodPost, "/api/proxy/reload", nil)
	reloadRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(reloadRec, reloadReq)
	if reloadRec.Code != http.StatusOK {
		t.Fatalf("unexpected reload endpoint code: %d", reloadRec.Code)
	}
	var reloadPayload map[string]string
	if err := json.Unmarshal(reloadRec.Body.Bytes(), &reloadPayload); err != nil {
		t.Fatal(err)
	}
	if reloadPayload["status"] != "reloaded" {
		t.Fatalf("expected reloaded status, got %#v", reloadPayload)
	}
}

func TestProxyValidateInternalErrorReturnsServerError(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.sqlite")
	store, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := db.Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	dnsSvc := dns.NewMockService(filepath.Join(t.TempDir(), "generated", "dns"))
	if err := dnsSvc.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	server := NewServer(store, events.NewMockBus(), failingValidateSquidService{}, dnsSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/proxy/validate", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for internal validate error, got %d", rec.Code)
	}
}

func TestDNSEndpoints(t *testing.T) {
	server := testServer(t)

	statusReq := httptest.NewRequest(http.MethodGet, "/api/dns/status", nil)
	statusRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(statusRec, statusReq)
	if statusRec.Code != http.StatusOK {
		t.Fatalf("unexpected dns status endpoint code: %d", statusRec.Code)
	}

	settingsReq := httptest.NewRequest(http.MethodGet, "/api/dns/settings", nil)
	settingsRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(settingsRec, settingsReq)
	if settingsRec.Code != http.StatusOK {
		t.Fatalf("unexpected dns settings endpoint code: %d", settingsRec.Code)
	}
	var settingsPayload map[string]string
	if err := json.Unmarshal(settingsRec.Body.Bytes(), &settingsPayload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(settingsPayload["config_preview"], "listen-address=127.0.0.1") {
		t.Fatalf("expected dns preview to contain listen-address, got %#v", settingsPayload)
	}

	validateReq := httptest.NewRequest(http.MethodPost, "/api/dns/validate", nil)
	validateRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(validateRec, validateReq)
	if validateRec.Code != http.StatusOK {
		t.Fatalf("unexpected dns validate endpoint code: %d", validateRec.Code)
	}

	reloadReq := httptest.NewRequest(http.MethodPost, "/api/dns/reload", nil)
	reloadRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(reloadRec, reloadReq)
	if reloadRec.Code != http.StatusOK {
		t.Fatalf("unexpected dns reload endpoint code: %d", reloadRec.Code)
	}
}

func TestDNSValidateInternalErrorReturnsServerError(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.sqlite")
	store, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := db.Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	squidSvc := squid.NewMockService(filepath.Join(t.TempDir(), "generated", "squid"))
	if err := squidSvc.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	server := NewServer(store, events.NewMockBus(), squidSvc, failingValidateDNSService{})
	req := httptest.NewRequest(http.MethodPost, "/api/dns/validate", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for internal dns validate error, got %d", rec.Code)
	}
}

func TestDNSValidateInvalidConfigReturnsBadRequest(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.sqlite")
	store, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := db.Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	squidSvc := squid.NewMockService(filepath.Join(t.TempDir(), "generated", "squid"))
	if err := squidSvc.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	server := NewServer(store, events.NewMockBus(), squidSvc, invalidValidateDNSService{})
	req := httptest.NewRequest(http.MethodPost, "/api/dns/validate", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid dns validate error, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["status"] != "invalid" {
		t.Fatalf("expected invalid status payload, got %#v", payload)
	}
}
