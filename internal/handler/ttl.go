package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/responder"
)

// TTLRequest represents a request to set a TTL for a key
type TTLRequest struct {
	TTL time.Duration `json:"ttl"` // in seconds
}

// SetTTL handles PUT /api/v1/ttl/{key}
func (h *Handler) SetTTL(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/ttl/")

	var req TTLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.Cache.SetTTL(key, req.TTL); err != nil {
		h.HandleError(w, err)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "TTL set successfully", map[string]any{
		"key": key,
		"ttl": req.TTL.Seconds(),
	})
}

// GetTTL handles GET /api/v1/ttl/{key}
func (h *Handler) GetTTL(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/ttl/")

	ttl, found := h.Cache.GetTTL(key)
	if !found {
		responder.WriteError(w, http.StatusNotFound, cache.ErrKeyNotFound)
		return
	}

	var ttlSeconds float64
	if ttl < 0 {
		ttlSeconds = -1 // No expiration
	} else {
		ttlSeconds = ttl.Seconds()
	}

	responder.WriteSuccess(w, http.StatusOK, "TTL retrieved successfully", map[string]any{
		"key": key,
		"ttl": ttlSeconds,
	})
}

// RemoveTTL handles DELETE /api/v1/ttl/{key}
func (h *Handler) RemoveTTL(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/ttl/")

	if err := h.Cache.RemoveTTL(key); err != nil {
		h.HandleError(w, err)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "TTL removed successfully", map[string]string{
		"key": key,
	})
}
