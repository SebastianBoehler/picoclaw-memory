package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"picoclaw-memory/internal/memory"
)

type Handler struct {
	memories *memory.Service
}

func NewHandler(memories *memory.Service) http.Handler {
	handler := &Handler{memories: memories}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.withMethod(http.MethodGet, handler.handleHealth))
	mux.HandleFunc("/v1/memories", handler.handleMemories)
	mux.HandleFunc("/v1/recall", handler.withMethod(http.MethodPost, handler.handleRecall))
	return mux
}

func (h *Handler) handleMemories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListMemories(w, r)
	case http.MethodPost:
		h.handleSaveMemory(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"service": "picoclaw-memory",
	})
}

func (h *Handler) handleSaveMemory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerTag string `json:"containerTag"`
		Content      string `json:"content"`
		Source       string `json:"source"`
		ExpiresAt    string `json:"expiresAt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	var expiresAt *time.Time
	if strings.TrimSpace(req.ExpiresAt) != "" {
		parsed, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "expiresAt must be RFC3339")
			return
		}
		expiresAt = &parsed
	}

	item, err := h.memories.Save(r.Context(), memory.SaveRequest{
		ContainerTag: req.ContainerTag,
		Content:      req.Content,
		Source:       req.Source,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"memory": item})
}

func (h *Handler) handleListMemories(w http.ResponseWriter, r *http.Request) {
	containerTag := r.URL.Query().Get("containerTag")
	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 20)

	items, err := h.memories.ListRecent(r.Context(), containerTag, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"memories": items})
}

func (h *Handler) handleRecall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerTag string `json:"containerTag"`
		Query        string `json:"query"`
		Limit        int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	results, err := h.memories.Recall(r.Context(), memory.RecallRequest{
		ContainerTag: req.ContainerTag,
		Query:        req.Query,
		Limit:        req.Limit,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func (h *Handler) withMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		next(w, r)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseIntOrDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
