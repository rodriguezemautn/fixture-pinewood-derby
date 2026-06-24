// Package router configura las rutas HTTP usando chi.
package router

import (
	"github.com/ema/fixture/backend/internal/auth"
	"github.com/ema/fixture/backend/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New crea un router chi con middleware y handlers registrados.
// Las rutas de escritura (POST/PUT/DELETE) requieren autenticación admin.
func New(handlers ...handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware global
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Auth middleware: protege POST/PUT/DELETE, GET es público
	r.Use(auth.Middleware)

	// Registrar handlers
	for _, h := range handlers {
		h.Register(r)
	}

	return r
}
