package events

import "github.com/vt887/macnet-gateway/daemon/internal/models"

type Bus interface {
	Recent() []models.Event
}

type MockBus struct{}

func NewMockBus() *MockBus {
	return &MockBus{}
}

func (m *MockBus) Recent() []models.Event {
	return []models.Event{
		{Time: "2026-01-01T12:00:00Z", Type: "SERVICE_STARTED", Target: "macnet-gatewayd", Action: "allowed"},
		{Time: "2026-01-01T12:02:00Z", Type: "DNS_QUERY", Target: "example.com", Action: "cached"},
	}
}
