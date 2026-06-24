package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type categoriaHandler struct {
	svc service.CategoriaService
}

// NewCategoriaHandler crea un handler REST para categorías.
func NewCategoriaHandler(svc service.CategoriaService) Handler {
	return &categoriaHandler{svc: svc}
}

func (h *categoriaHandler) Register(r Router) {
	r.Get("/api/categorias", h.List)
	r.Post("/api/categorias", h.Create)
	r.Get("/api/categorias/{id}", h.GetByID)
	r.Put("/api/categorias/{id}", h.Update)
	r.Delete("/api/categorias/{id}", h.Delete)
}

type createCategoriaRequest struct {
	Nombre  string `json:"nombre"`
	EdadMin int    `json:"edad_min"`
	EdadMax int    `json:"edad_max"`
}

func (h *categoriaHandler) List(w http.ResponseWriter, r *http.Request) {
	categorias, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Ensure we return [] not nil
	if categorias == nil {
		categorias = []domain.Categoria{}
	}
	writeJSON(w, http.StatusOK, categorias)
}

func (h *categoriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	c, err := h.svc.Create(req.Nombre, req.EdadMin, req.EdadMax)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

func (h *categoriaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if c == nil {
		writeError(w, http.StatusNotFound, "categoría no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *categoriaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req createCategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	c, err := h.svc.Update(id, req.Nombre, req.EdadMin, req.EdadMax)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "categoria "+id+" not found" {
			code = http.StatusNotFound
		}
		writeError(w, code, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, c)
}

func (h *categoriaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
