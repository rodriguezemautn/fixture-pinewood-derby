// Package service define interfaces de lógica de negocio.
package service

import "github.com/ema/fixture/backend/internal/domain"

// CategoriaService define operaciones de negocio sobre categorías.
type CategoriaService interface {
	List() ([]domain.Categoria, error)
	GetByID(id string) (*domain.Categoria, error)
	Create(nombre string, edadMin, edadMax int) (*domain.Categoria, error)
	Update(id, nombre string, edadMin, edadMax int) (*domain.Categoria, error)
	Delete(id string) error
}

// AutoService define operaciones de negocio sobre autos.
type AutoService interface {
	ListByCategoria(categoriaID string) ([]domain.Auto, error)
	ListAll() ([]domain.Auto, error)
	GetByID(id string) (*domain.Auto, error)
	Create(categoriaID string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error)
	Update(id string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error)
	Delete(id string) error
}

// CarreraService define operaciones de negocio sobre carreras.
type CarreraService interface {
	RegistrarOrdenLlegada(carreraID string, orden []string) error
}
