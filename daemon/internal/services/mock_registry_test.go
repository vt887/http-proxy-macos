package services

import "testing"

func TestMockRegistry(t *testing.T) {
	all := NewMockRegistry()
	if len(all) < 2 {
		t.Fatal("expected seeded service statuses")
	}
	if all[0].Name != "macnet-gatewayd" {
		t.Fatalf("unexpected first service: %s", all[0].Name)
	}
}
