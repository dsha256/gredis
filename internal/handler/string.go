package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/responder"
)

// StringRequest represents a request to set a string value
type StringRequest struct {
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl,omitempty"` // in seconds
}

// GetString handles GET /api/v1/string/{key}
func (h *Handler) GetString(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/string/")

	value, found := h.Cache.Get(key)
	if !found {
		h.HandleError(w, cache.ErrKeyNotFound)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Value retrieved successfully", map[string]string{
		"key":   key,
		"value": value,
	})
}

// SetString handles POST /api/v1/string/{key}
func (h *Handler) SetString(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/string/")

	var req StringRequest
	if !h.DecodeJSON(w, r, &req) {
		return
	}

	var err error
	if req.TTL > 0 {
		err = h.Cache.SetWithTTL(key, req.Value, req.TTL*time.Second)
	} else {
		err = h.Cache.Set(key, req.Value)
	}

	if h.HandleError(w, err) {
		return
	}

	responder.WriteSuccess(w, http.StatusCreated, "Value set successfully", map[string]string{
		"key":   key,
		"value": req.Value,
	})
}

// UpdateString handles PUT /api/v1/string/{key}
func (h *Handler) UpdateString(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/v1/string/")

	var req StringRequest
	if !h.DecodeJSON(w, r, &req) {
		return
	}

	if err := h.Cache.Update(key, req.Value); err != nil {
		h.HandleError(w, err)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Value updated successfully", map[string]string{
		"key":   key,
		"value": req.Value,
	})
}
