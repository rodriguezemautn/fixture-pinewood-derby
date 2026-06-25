package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ema/fixture/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

func parseQueryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}

type competenciaHandler struct {
	svc service.FixtureService
}

func NewCompetenciaHandler(svc service.FixtureService) Handler {
	return &competenciaHandler{svc: svc}
}

func (h *competenciaHandler) Register(r Router) {
	r.Get("/api/categorias/{categoriaId}/competencias", h.Listar)
	r.Post("/api/categorias/{categoriaId}/competencias", h.Crear)
	r.Post("/api/competencias/{id}/finalizar", h.Finalizar)
	r.Post("/api/competencias/{id}/desempate", h.Desempate)
}

func (h *competenciaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	comps, err := h.svc.ListarCompetencias(chi.URLParam(r, "categoriaId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comps)
}

func (h *competenciaHandler) Crear(w http.ResponseWriter, r *http.Request) {
	rondas := parseQueryInt(r.URL.Query().Get("rondas"), 3)
	comp, f, err := h.svc.CrearCompetencia(chi.URLParam(r, "categoriaId"), rondas)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"competencia": comp,
		"fixture":     f,
	})
}

func (h *competenciaHandler) Finalizar(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.FinalizarCompetencia(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *competenciaHandler) Desempate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AutoIDs []string `json:"auto_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	heat, err := h.svc.AgregarDesempate(chi.URLParam(r, "id"), req.AutoIDs)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, heat)
}
