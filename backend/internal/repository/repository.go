// Package repository define interfaces de persistencia.
package repository

import "github.com/ema/fixture/backend/internal/domain"

// CategoriaRepository define operaciones sobre categorías.
type CategoriaRepository interface {
	List() ([]domain.Categoria, error)
	GetByID(id string) (*domain.Categoria, error)
	Save(c *domain.Categoria) error
	Update(c *domain.Categoria) error
	Delete(id string) error
}

// AutoRepository define operaciones sobre autos.
type AutoRepository interface {
	ListByCategoria(categoriaID string) ([]domain.Auto, error)
	GetByID(id string) (*domain.Auto, error)
	Save(a *domain.Auto) error
	Update(a *domain.Auto) error
	Delete(id string) error
	ExistsByNumero(categoriaID string, numero int) (bool, error)
}

// CarreraRepository define operaciones sobre carreras.
type CarreraRepository interface {
	ListByCategoria(categoriaID string) ([]domain.Carrera, error)
	GetByID(id string) (*domain.Carrera, error)
	Save(c *domain.Carrera) error
	UpdateOrdenLlegada(id string, orden []string) error
}
