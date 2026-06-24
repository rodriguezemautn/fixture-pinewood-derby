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

// mockCategoriaService implementa service.CategoriaService en memoria.
type mockCategoriaService struct {
	data map[string]*domain.Categoria
}

func (m *mockCategoriaService) List() ([]domain.Categoria, error) {
	var res []domain.Categoria
	for _, c := range m.data {
		res = append(res, *c)
	}
	return res, nil
}

func (m *mockCategoriaService) GetByID(id string) (*domain.Categoria, error) {
	c, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockCategoriaService) Create(nombre string, edadMin, edadMax int) (*domain.Categoria, error) {
	c := &domain.Categoria{ID: "new-id", Nombre: nombre, EdadMin: edadMin, EdadMax: edadMax}
	m.data[c.ID] = c
	return c, nil
}

func (m *mockCategoriaService) Update(id, nombre string, edadMin, edadMax int) (*domain.Categoria, error) {
	c, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	c.Nombre = nombre
	c.EdadMin = edadMin
	c.EdadMax = edadMax
	return c, nil
}

func (m *mockCategoriaService) Delete(id string) error {
	delete(m.data, id)
	return nil
}

func newMockSvc() *mockCategoriaService {
	return &mockCategoriaService{data: make(map[string]*domain.Categoria)}
}

func setupCategoriaTest(t *testing.T) (*chi.Mux, *mockCategoriaService) {
	t.Helper()
	r := chi.NewRouter()
	svc := newMockSvc()
	h := NewCategoriaHandler(svc)
	h.Register(r)
	return r, svc
}

func TestCategoriaHandler_List(t *testing.T) {
	r, svc := setupCategoriaTest(t)
	svc.data["1"] = &domain.Categoria{ID: "1", Nombre: "Test", EdadMin: 1, EdadMax: 5}

	req := httptest.NewRequest(http.MethodGet, "/api/categorias", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var categorias []domain.Categoria
	json.NewDecoder(rec.Body).Decode(&categorias)
	if len(categorias) != 1 {
		t.Errorf("expected 1 categoria, got %d", len(categorias))
	}
}

func TestCategoriaHandler_Create_Success(t *testing.T) {
	r, _ := setupCategoriaTest(t)

	body := bytes.NewBufferString(`{"nombre":"Juveniles","edad_min":13,"edad_max":15}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categorias", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCategoriaHandler_Create_BadJSON(t *testing.T) {
	r, _ := setupCategoriaTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/categorias", bytes.NewBufferString("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCategoriaHandler_GetByID_Found(t *testing.T) {
	r, svc := setupCategoriaTest(t)
	svc.data["abc"] = &domain.Categoria{ID: "abc", Nombre: "Test", EdadMin: 1, EdadMax: 5}

	req := httptest.NewRequest(http.MethodGet, "/api/categorias/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCategoriaHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupCategoriaTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/categorias/no-existe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCategoriaHandler_Update_Success(t *testing.T) {
	r, svc := setupCategoriaTest(t)
	svc.data["1"] = &domain.Categoria{ID: "1", Nombre: "Old", EdadMin: 1, EdadMax: 5}

	body := bytes.NewBufferString(`{"nombre":"New","edad_min":6,"edad_max":10}`)
	req := httptest.NewRequest(http.MethodPut, "/api/categorias/1", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCategoriaHandler_Delete_Success(t *testing.T) {
	r, svc := setupCategoriaTest(t)
	svc.data["1"] = &domain.Categoria{ID: "1", Nombre: "Test", EdadMin: 1, EdadMax: 5}

	req := httptest.NewRequest(http.MethodDelete, "/api/categorias/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
