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
