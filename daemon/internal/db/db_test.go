package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeCreatesTables(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db", "app.sqlite")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestServiceStatusPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db", "app.sqlite")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	err = UpsertServiceStatus(context.Background(), store, ServiceStatusRecord{
		Name:    "macnet-gatewayd",
		Status:  "running",
		Message: "healthy",
	})
	if err != nil {
		t.Fatal(err)
	}

	statuses, err := ListServiceStatuses(context.Background(), store)
	if err != nil {
		t.Fatal(err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected exactly one persisted service status, got %d", len(statuses))
	}
	if statuses[0].Name != "macnet-gatewayd" || statuses[0].Status != "running" || statuses[0].Message != "healthy" {
		t.Fatalf("unexpected persisted status: %#v", statuses[0])
	}

	err = UpsertServiceStatus(context.Background(), store, ServiceStatusRecord{
		Name:    "macnet-gatewayd",
		Status:  "stopped",
		Message: "maintenance",
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := ListServiceStatuses(context.Background(), store)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated) != 1 {
		t.Fatalf("expected one row after update upsert, got %d", len(updated))
	}
	if updated[0].Status != "stopped" || updated[0].Message != "maintenance" {
		t.Fatalf("expected updated status/message, got %#v", updated[0])
	}
}

func TestInsertSettingIfMissingDoesNotOverrideExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db", "app.sqlite")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := Initialize(context.Background(), store); err != nil {
		t.Fatal(err)
	}

	if err := UpsertSetting(context.Background(), store, "ui.theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := InsertSettingIfMissing(context.Background(), store, "ui.theme", "system"); err != nil {
		t.Fatal(err)
	}

	value, err := GetSetting(context.Background(), store, "ui.theme")
	if err != nil {
		t.Fatal(err)
	}
	if value != "dark" {
		t.Fatalf("expected existing value to remain unchanged, got %q", value)
	}
}
