package handler

import (
	"log/slog"
	"net/http"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/middleware"
)

// Handler contains the dependencies for all handlers
type Handler struct {
	Cache  cache.Cache
	Logger *slog.Logger
}

// New creates a new Handler with the given dependencies
func New(cache cache.Cache, logger *slog.Logger) *Handler {
	return &Handler{
		Cache:  cache,
		Logger: logger,
	}
}

// RegisterRoutes registers all the routes for the cache API
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// String operations
	mux.Handle("GET /api/v1/string/{key}", h.wrapHandler(h.GetString))
	mux.Handle("POST /api/v1/string/{key}", h.wrapHandler(h.SetString))
	mux.Handle("PUT /api/v1/string/{key}", h.wrapHandler(h.UpdateString))

	// List operations
	mux.Handle("POST /api/v1/list/{key}/front", h.wrapHandler(h.PushFront))
	mux.Handle("POST /api/v1/list/{key}/back", h.wrapHandler(h.PushBack))
	mux.Handle("DELETE /api/v1/list/{key}/front", h.wrapHandler(h.PopFront))
	mux.Handle("DELETE /api/v1/list/{key}/back", h.wrapHandler(h.PopBack))
	mux.Handle("GET /api/v1/list/{key}/range", h.wrapHandler(h.ListRange))

	// TTL operations
	mux.Handle("PUT /api/v1/ttl/{key}", h.wrapHandler(h.SetTTL))
	mux.Handle("GET /api/v1/ttl/{key}", h.wrapHandler(h.GetTTL))
	mux.Handle("DELETE /api/v1/ttl/{key}", h.wrapHandler(h.RemoveTTL))

	// General operations
	mux.Handle("DELETE /api/v1/key/{key}", h.wrapHandler(h.Remove))
	mux.Handle("GET /api/v1/key/{key}/exists", h.wrapHandler(h.Exists))
	mux.Handle("GET /api/v1/key/{key}/type", h.wrapHandler(h.Type))
	mux.Handle("DELETE /api/v1/keys", h.wrapHandler(h.Clear))
}

func (h *Handler) wrapHandler(handler http.HandlerFunc) http.Handler {
	return middleware.LoggingMiddleware(
		h.Logger,
		middleware.RecoveryMiddleware(
			h.Logger,
			handler,
		),
	)
}
