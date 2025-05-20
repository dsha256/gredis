package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/responder"
)

// Remove handles DELETE /api/key/{key}
func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/key/")

	if err := h.Cache.Remove(key); err != nil {
		h.HandleError(w, err)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Key removed successfully", map[string]string{
		"key": key,
	})
}

// Exists handles GET /api/key/{key}/exists
func (h *Handler) Exists(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/key/")
	key = strings.TrimSuffix(key, "/exists")

	exists := h.Cache.Exists(key)

	responder.WriteSuccess(w, http.StatusOK, "Key existence checked", map[string]any{
		"key":    key,
		"exists": exists,
	})
}

// Type handles GET /api/key/{key}/type
func (h *Handler) Type(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/key/")
	key = strings.TrimSuffix(key, "/type")

	dataType, found := h.Cache.Type(key)
	if !found {
		responder.WriteError(w, http.StatusNotFound, cache.ErrKeyNotFound)
		return
	}

	var typeStr string
	switch dataType {
	case cache.StringType:
		typeStr = "string"
	case cache.ListType:
		typeStr = "list"
	default:
		typeStr = "unknown"
	}

	responder.WriteSuccess(w, http.StatusOK, "Key type retrieved successfully", map[string]string{
		"key":  key,
		"type": typeStr,
	})
}

// Clear handles DELETE /api/keys
func (h *Handler) Clear(w http.ResponseWriter, _ *http.Request) {
	err := h.Cache.Clear()
	if err != nil {
		responder.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Cache cleared successfully", json.RawMessage{})
}
