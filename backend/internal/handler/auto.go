package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type autoHandler struct {
	svc service.AutoService
}

// NewAutoHandler crea un handler REST para autos.
func NewAutoHandler(svc service.AutoService) Handler {
	return &autoHandler{svc: svc}
}

func (h *autoHandler) Register(r Router) {
	// Rutas anidadas por categoría
	r.Get("/api/categorias/{categoriaId}/autos", h.ListByCategoria)
	r.Post("/api/categorias/{categoriaId}/autos", h.Create)
	// Rutas directas por ID
	r.Get("/api/autos/{id}", h.GetByID)
	r.Put("/api/autos/{id}", h.Update)
	r.Delete("/api/autos/{id}", h.Delete)
}

type createAutoRequest struct {
	Numero  int    `json:"numero"`
	Nombre  string `json:"nombre"`
	Creador string `json:"creador"`
	Edad    int    `json:"edad"`
	FotoURL string `json:"foto_url"`
}

func (h *autoHandler) ListByCategoria(w http.ResponseWriter, r *http.Request) {
	categoriaID := chi.URLParam(r, "categoriaId")
	autos, err := h.svc.ListByCategoria(categoriaID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if autos == nil {
		autos = []domain.Auto{}
	}
	writeJSON(w, http.StatusOK, autos)
}

func (h *autoHandler) Create(w http.ResponseWriter, r *http.Request) {
	categoriaID := chi.URLParam(r, "categoriaId")

	var req createAutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	a, err := h.svc.Create(categoriaID, req.Numero, req.Nombre, req.Creador, req.Edad, req.FotoURL)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "categoria "+categoriaID+" no encontrada" {
			code = http.StatusNotFound
		}
		writeError(w, code, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

func (h *autoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if a == nil {
		writeError(w, http.StatusNotFound, "auto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *autoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req createAutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	a, err := h.svc.Update(id, req.Numero, req.Nombre, req.Creador, req.Edad, req.FotoURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, a)
}

func (h *autoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
