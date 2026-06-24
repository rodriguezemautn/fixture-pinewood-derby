package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type mockAutoService struct {
	autos map[string]*domain.Auto
}

func (m *mockAutoService) ListAll() ([]domain.Auto, error) {
	var res []domain.Auto
	for _, a := range m.autos {
		res = append(res, *a)
	}
	return res, nil
}

func (m *mockAutoService) ListByCategoria(categoriaID string) ([]domain.Auto, error) {
	var res []domain.Auto
	for _, a := range m.autos {
		if a.CategoriaID == categoriaID {
			res = append(res, *a)
		}
	}
	return res, nil
}
func (m *mockAutoService) GetByID(id string) (*domain.Auto, error) {
	a, ok := m.autos[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}
func (m *mockAutoService) Create(categoriaID string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error) {
	a := &domain.Auto{ID: "new-auto", CategoriaID: categoriaID, Numero: numero, Nombre: nombre, Creador: creador, Edad: edad, Peso: peso, FotoURL: fotoURL}
	m.autos[a.ID] = a
	return a, nil
}
func (m *mockAutoService) Update(id string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error) {
	a, ok := m.autos[id]
	if !ok {
		return nil, nil
	}
	a.Numero = numero
	a.Nombre = nombre
	a.Creador = creador
	a.Edad = edad
	a.Peso = peso
	a.FotoURL = fotoURL
	return a, nil
}
func (m *mockAutoService) Delete(id string) error {
	delete(m.autos, id)
	return nil
}

func setupAutoTest(t *testing.T) (*chi.Mux, *mockAutoService) {
	t.Helper()
	r := chi.NewRouter()
	svc := &mockAutoService{autos: make(map[string]*domain.Auto)}
	h := NewAutoHandler(svc)
	h.Register(r)
	return r, svc
}

func TestAutoHandler_ListByCategoria(t *testing.T) {
	r, svc := setupAutoTest(t)
	svc.autos["a1"] = &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "A", Creador: "J", Edad: 10}

	req := httptest.NewRequest(http.MethodGet, "/api/categorias/cat-1/autos", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var autos []domain.Auto
	json.NewDecoder(rec.Body).Decode(&autos)
	if len(autos) != 1 {
		t.Errorf("expected 1 auto, got %d", len(autos))
	}
}

func TestAutoHandler_ListByCategoria_Empty(t *testing.T) {
	r, _ := setupAutoTest(t)
	req := httptest.NewRequest(http.MethodGet, "/api/categorias/cat-99/autos", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAutoHandler_Create_Success(t *testing.T) {
	r, _ := setupAutoTest(t)

	body := bytes.NewBufferString(`{"numero":1,"nombre":"Turbo","creador":"Juan","edad":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categorias/cat-1/autos", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAutoHandler_GetByID_Found(t *testing.T) {
	r, svc := setupAutoTest(t)
	svc.autos["a1"] = &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "A", Creador: "J", Edad: 10}

	req := httptest.NewRequest(http.MethodGet, "/api/autos/a1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAutoHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupAutoTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/autos/no-existe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestAutoHandler_Update_Success(t *testing.T) {
	r, svc := setupAutoTest(t)
	svc.autos["a1"] = &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "Old", Creador: "J", Edad: 10}

	body := bytes.NewBufferString(`{"numero":2,"nombre":"New","creador":"M","edad":12}`)
	req := httptest.NewRequest(http.MethodPut, "/api/autos/a1", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAutoHandler_Delete_Success(t *testing.T) {
	r, svc := setupAutoTest(t)
	svc.autos["a1"] = &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "A", Creador: "J", Edad: 10}

	req := httptest.NewRequest(http.MethodDelete, "/api/autos/a1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
