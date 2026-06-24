// Package service define interfaces de lógica de negocio.
package service

import "github.com/ema/fixture/backend/internal/domain"

// CategoriaService define operaciones de negocio sobre categorías.
type CategoriaService interface {
	List() ([]domain.Categoria, error)
	Create(nombre string, edadMin, edadMax int) (*domain.Categoria, error)
}

// CarreraService define operaciones de negocio sobre carreras.
type CarreraService interface {
	RegistrarOrdenLlegada(carreraID string, orden []string) error
}
