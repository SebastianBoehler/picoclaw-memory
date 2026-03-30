package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"picoclaw-memory/internal/config"
	"picoclaw-memory/internal/httpapi"
	"picoclaw-memory/internal/memory"
	sqlitestore "picoclaw-memory/internal/storage/sqlite"
)

type App struct {
	cfg     config.Config
	store   *sqlitestore.Store
	handler http.Handler
}

func New(cfg config.Config) (*App, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
		return nil, fmt.Errorf("create sqlite directory: %w", err)
	}

	store, err := sqlitestore.Open(cfg.SQLitePath)
	if err != nil {
		return nil, err
	}

	service := memory.NewService(store)
	handler := httpapi.NewHandler(service)

	return &App{
		cfg:     cfg,
		store:   store,
		handler: handler,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.store.Close()
	}()

	return httpapi.Run(ctx, a.cfg.ListenAddr, a.handler)
}
