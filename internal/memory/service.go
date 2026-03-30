package memory

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Save(ctx context.Context, req SaveRequest) (Memory, error) {
	if strings.TrimSpace(req.ContainerTag) == "" {
		return Memory{}, fmt.Errorf("containerTag is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return Memory{}, fmt.Errorf("content is required")
	}

	item := Memory{
		ID:           newID(),
		ContainerTag: strings.TrimSpace(req.ContainerTag),
		Content:      strings.TrimSpace(req.Content),
		Source:       strings.TrimSpace(req.Source),
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    req.ExpiresAt,
	}

	if err := s.store.Save(ctx, item); err != nil {
		return Memory{}, err
	}

	return item, nil
}

func (s *Service) Recall(ctx context.Context, req RecallRequest) ([]RecallResult, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, fmt.Errorf("query is required")
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	req.Query = strings.TrimSpace(req.Query)
	req.ContainerTag = strings.TrimSpace(req.ContainerTag)

	return s.store.Recall(ctx, req)
}

func (s *Service) ListRecent(ctx context.Context, containerTag string, limit int) ([]Memory, error) {
	if strings.TrimSpace(containerTag) == "" {
		return nil, fmt.Errorf("containerTag is required")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.store.ListRecent(ctx, strings.TrimSpace(containerTag), limit)
}

func newID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("mem_%d", time.Now().UTC().UnixNano())
	}
	return "mem_" + hex.EncodeToString(buf)
}
