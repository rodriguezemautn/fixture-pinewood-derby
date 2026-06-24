package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

var UploadDir = "uploads"

type autoHandler struct {
	svc service.AutoService
}

// NewAutoHandler crea un handler REST para autos.
func NewAutoHandler(svc service.AutoService) Handler {
	return &autoHandler{svc: svc}
}

func (h *autoHandler) Register(r Router) {
	r.Get("/api/categorias/{categoriaId}/autos", h.ListByCategoria)
	r.Post("/api/categorias/{categoriaId}/autos", h.Create)
	r.Get("/api/autos", h.ListAll)
	r.Get("/api/autos/{id}", h.GetByID)
	r.Put("/api/autos/{id}", h.Update)
	r.Delete("/api/autos/{id}", h.Delete)
	r.Post("/api/autos/{id}/foto", h.UploadFoto)
}

// UploadFoto recibe un archivo de foto y lo guarda en el FS local.
func (h *autoHandler) UploadFoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Verificar que el auto existe
	a, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if a == nil {
		writeError(w, http.StatusNotFound, "auto no encontrado")
		return
	}

	// Parsear multipart form (max 10MB)
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("foto")
	if err != nil {
		writeError(w, http.StatusBadRequest, "archivo 'foto' requerido")
		return
	}
	defer file.Close()

	// Validar extensión
	ext := filepath.Ext(header.Filename)
	extValida := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !extValida[ext] {
		writeError(w, http.StatusBadRequest, "formato no soportado: use jpg, png, gif o webp")
		return
	}

	// Crear directorio uploads si no existe
	os.MkdirAll(UploadDir, 0755)

	// Nombre único
	nombre := fmt.Sprintf("%s_%d%s", id, time.Now().Unix(), ext)
	ruta := filepath.Join(UploadDir, nombre)

	// Guardar archivo
	dst, err := os.Create(ruta)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al guardar archivo")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "error al escribir archivo")
		return
	}

	// Actualizar FotoURL en la DB
	fotoURL := "/uploads/" + nombre
	if _, err := h.svc.Update(id, a.Numero, a.Nombre, a.Creador, a.Edad, fotoURL); err != nil {
		// Si falla la actualización, intentamos borrar el archivo
		os.Remove(ruta)
		writeError(w, http.StatusInternalServerError, "error al actualizar auto")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"foto_url": fotoURL})
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

func (h *autoHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	// Necesitamos acceso a todas las categorías para agrupar
	// Por ahora listamos todos los autos planos
	autos, err := h.svc.ListAll()
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
