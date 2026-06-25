package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type fixtureHandler struct {
	svc service.FixtureService
}

func NewFixtureHandler(svc service.FixtureService) Handler {
	return &fixtureHandler{svc: svc}
}

func (h *fixtureHandler) Register(r Router) {
	r.Post("/api/categorias/{categoriaId}/fixture", h.Generar)
	r.Get("/api/categorias/{categoriaId}/fixture", h.Obtener)
	r.Get("/api/categorias/{categoriaId}/posiciones", h.Posiciones)
	r.Post("/api/carreras/{heatId}/resultado", h.RegistrarResultado)
	r.Post("/api/categorias/{categoriaId}/final", h.GenerarFinal)
	r.Post("/api/categorias/{categoriaId}/archivar", h.Archivar)
	r.Get("/api/categorias/{categoriaId}/archivos", h.GetArchivos)
}

func (h *fixtureHandler) Generar(w http.ResponseWriter, r *http.Request) {
	categoriaID := chi.URLParam(r, "categoriaId")
	rondas := 3
	if r.URL.Query().Get("rondas") != "" {
		if v, err := strconv.Atoi(r.URL.Query().Get("rondas")); err == nil && v > 0 {
			rondas = v
		}
	}

	f, err := h.svc.Generar(categoriaID, rondas)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

func (h *fixtureHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	f, err := h.svc.ObtenerFixture(chi.URLParam(r, "categoriaId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		writeError(w, http.StatusNotFound, "no hay fixture para esta categoría")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (h *fixtureHandler) Posiciones(w http.ResponseWriter, r *http.Request) {
	standings, err := h.svc.ObtenerPosiciones(chi.URLParam(r, "categoriaId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if standings == nil {
		standings = []domain.Standing{}
	}
	writeJSON(w, http.StatusOK, standings)
}

func (h *fixtureHandler) RegistrarResultado(w http.ResponseWriter, r *http.Request) {
	heatID := chi.URLParam(r, "heatId")

	var req struct {
		OrdenLlegada []string `json:"orden_llegada"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.svc.RegistrarResultado(heatID, req.OrdenLlegada); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *fixtureHandler) GenerarFinal(w http.ResponseWriter, r *http.Request) {
	final, err := h.svc.GenerarFinal(chi.URLParam(r, "categoriaId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, final)
}

func (h *fixtureHandler) Archivar(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Archivar(chi.URLParam(r, "categoriaId")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *fixtureHandler) GetArchivos(w http.ResponseWriter, r *http.Request) {
	archivos, err := h.svc.GetArchivos(chi.URLParam(r, "categoriaId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if archivos == nil {
		archivos = []map[string]any{}
	}
	writeJSON(w, http.StatusOK, archivos)
}
