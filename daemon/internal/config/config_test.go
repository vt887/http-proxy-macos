package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddress != "127.0.0.1:18080" {
		t.Fatalf("unexpected listen address: %s", cfg.ListenAddress)
	}
	if cfg.SquidGeneratedConfigDir == "" {
		t.Fatal("expected default squid generated config dir")
	}
}

func TestLoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "daemon.json")
	content := `{"listen_address":"127.0.0.1:19090","database_path":"/tmp/test.sqlite","squid_generated_config_dir":"/tmp/generated/squid"}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddress != "127.0.0.1:19090" {
		t.Fatalf("unexpected listen address: %s", cfg.ListenAddress)
	}
	if cfg.SquidGeneratedConfigDir != "/tmp/generated/squid" {
		t.Fatalf("unexpected squid config dir: %s", cfg.SquidGeneratedConfigDir)
	}
}
