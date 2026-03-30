package config

import (
	"fmt"
	"os"
)

type Config struct {
	ListenAddr string
	DataDir    string
	SQLitePath string
}

func Load() (Config, error) {
	cfg := Config{
		ListenAddr: envOrDefault("LISTEN_ADDR", ":8080"),
		DataDir:    envOrDefault("DATA_DIR", "./var"),
		SQLitePath: envOrDefault("SQLITE_PATH", "./var/memory.db"),
	}

	if cfg.ListenAddr == "" {
		return Config{}, fmt.Errorf("LISTEN_ADDR must not be empty")
	}
	if cfg.DataDir == "" {
		return Config{}, fmt.Errorf("DATA_DIR must not be empty")
	}
	if cfg.SQLitePath == "" {
		return Config{}, fmt.Errorf("SQLITE_PATH must not be empty")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
