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
	if len(statuses) == 0 {
		t.Fatal("expected at least one persisted service status")
	}
}
