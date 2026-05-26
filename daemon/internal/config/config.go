package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultDevRoot = ".macnet-gateway-dev"

type Config struct {
	ListenAddress           string `json:"listen_address"`
	DatabasePath            string `json:"database_path"`
	SquidGeneratedConfigDir string `json:"squid_generated_config_dir"`
}

func Load(path string) (Config, error) {
	cfg := defaults()
	if path == "" {
		return cfg, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, errors.New("config must be valid JSON for PR-1 scaffold")
	}
	return cfg, nil
}

func defaults() Config {
	home, _ := os.UserHomeDir()
	devRoot := filepath.Join(home, defaultDevRoot)
	return Config{
		ListenAddress:           "127.0.0.1:18080",
		DatabasePath:            filepath.Join(devRoot, "db", "app.sqlite"),
		SquidGeneratedConfigDir: filepath.Join(devRoot, "generated", "squid"),
	}
}
