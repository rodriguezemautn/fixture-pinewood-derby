package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerarYVerificarToken(t *testing.T) {
	token := GenerarToken()
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Middleware test: GET should pass without auth
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/categorias", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("GET should pass without auth, got %d", rec.Code)
	}

	// POST without token should fail
	req = httptest.NewRequest(http.MethodPost, "/api/categorias", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("POST without token should return 401, got %d", rec.Code)
	}

	// POST with valid token should pass
	req = httptest.NewRequest(http.MethodPost, "/api/categorias", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("POST with valid token should pass, got %d", rec.Code)
	}

	// POST with invalid token should fail
	req = httptest.NewRequest(http.MethodPost, "/api/categorias", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("POST with invalid token should return 401, got %d", rec.Code)
	}

	// Login endpoint should pass without auth
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("POST /api/auth/login should pass without auth, got %d", rec.Code)
	}
}
