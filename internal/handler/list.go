package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/responder"
)

// ListRequest represents a request to add a value to a list
type ListRequest struct {
	Value string `json:"value"`
}

// ListRangeRequest represents a request to get a range of values from a list
type ListRangeRequest struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// PushFront handles POST /api/list/{key}/front
func (h *Handler) PushFront(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/list/")
	key = strings.TrimSuffix(key, "/front")

	var req ListRequest
	if !h.DecodeJSON(w, r, &req) {
		return
	}

	err := h.Cache.PushFront(key, req.Value)
	if h.HandleError(w, err) {
		return
	}

	responder.WriteSuccess(w, http.StatusCreated, "Value pushed to front of list successfully", map[string]string{
		"key":   key,
		"value": req.Value,
	})
}

// PushBack handles POST /api/list/{key}/back
func (h *Handler) PushBack(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/list/")
	key = strings.TrimSuffix(key, "/back")

	var req ListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.Cache.PushBack(key, req.Value); err != nil {
		h.HandleError(w, err)
		return
	}

	responder.WriteSuccess(w, http.StatusCreated, "Value pushed to back of list successfully", map[string]string{
		"key":   key,
		"value": req.Value,
	})
}

// PopFront handles DELETE /api/list/{key}/front
func (h *Handler) PopFront(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/list/")
	key = strings.TrimSuffix(key, "/front")

	value, found := h.Cache.PopFront(key)
	if !found {
		responder.WriteError(w, http.StatusNotFound, cache.ErrKeyNotFound)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Value popped from front of list successfully", map[string]string{
		"key":   key,
		"value": value,
	})
}

// PopBack handles DELETE /api/list/{key}/back
func (h *Handler) PopBack(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/list/")
	key = strings.TrimSuffix(key, "/back")

	value, found := h.Cache.PopBack(key)
	if !found {
		responder.WriteError(w, http.StatusNotFound, cache.ErrKeyNotFound)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Value popped from back of list successfully", map[string]string{
		"key":   key,
		"value": value,
	})
}

// ListRange handles GET /api/list/{key}/range
func (h *Handler) ListRange(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/list/")
	key = strings.TrimSuffix(key, "/range")

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := strconv.Atoi(startStr)
	if err != nil {
		responder.WriteError(w, http.StatusBadRequest, err)
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		responder.WriteError(w, http.StatusBadRequest, err)
		return
	}

	values, err := h.Cache.ListRange(key, start, end)
	if h.HandleError(w, err) {
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "List range retrieved successfully", map[string]any{
		"key":    key,
		"start":  start,
		"end":    end,
		"values": values,
	})
}
