package handler

import (
	"encoding/json"
	"net/http"
)

type healthHandler struct{}

// NewHealthHandler crea un handler para health check.
func NewHealthHandler() Handler {
	return &healthHandler{}
}

func (h *healthHandler) Register(r Router) {
	r.Get("/health", h.Check)
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
