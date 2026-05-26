package events

import "testing"

func TestMockBusReturnsEvents(t *testing.T) {
	bus := NewMockBus()
	events := bus.Recent()
	if len(events) == 0 {
		t.Fatal("expected mock events")
	}
}
