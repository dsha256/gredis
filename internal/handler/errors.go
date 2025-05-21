package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/responder"
)

func (h *Handler) HandleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError

	switch {
	case errors.Is(err, cache.ErrKeyNotFound):
		responder.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, cache.ErrTypeMismatch):
		responder.WriteError(w, http.StatusBadRequest, err)
	case errors.As(err, &syntaxErr) || errors.As(err, &unmarshalTypeErr):
		responder.WriteError(w, http.StatusBadRequest, errors.New("invalid request format"))
	default:
		h.Logger.Error("Internal server error", "error", err)
		responder.WriteError(w, http.StatusInternalServerError, err)
	}

	return true
}

func (h *Handler) DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.HandleError(w, err)
		return false
	}
	return true
}
