package service

import (
	"fmt"
	"strings"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/repository"
	"github.com/google/uuid"
)

type categoriaService struct {
	repo repository.CategoriaRepository
}

// NewCategoriaService crea un service concreto de categorías.
func NewCategoriaService(repo repository.CategoriaRepository) CategoriaService {
	return &categoriaService{repo: repo}
}

func (s *categoriaService) List() ([]domain.Categoria, error) {
	return s.repo.List()
}

func (s *categoriaService) GetByID(id string) (*domain.Categoria, error) {
	return s.repo.GetByID(id)
}

func (s *categoriaService) Create(nombre string, edadMin, edadMax int) (*domain.Categoria, error) {
	if err := validate(nombre, edadMin, edadMax); err != nil {
		return nil, err
	}

	c := &domain.Categoria{
		ID:      uuid.New().String(),
		Nombre:  strings.TrimSpace(nombre),
		EdadMin: edadMin,
		EdadMax: edadMax,
	}

	if err := s.repo.Save(c); err != nil {
		return nil, fmt.Errorf("save categoria: %w", err)
	}
	return c, nil
}

func (s *categoriaService) Update(id, nombre string, edadMin, edadMax int) (*domain.Categoria, error) {
	if err := validate(nombre, edadMin, edadMax); err != nil {
		return nil, err
	}

	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get categoria for update: %w", err)
	}
	if c == nil {
		return nil, fmt.Errorf("categoria %s not found", id)
	}

	c.Nombre = strings.TrimSpace(nombre)
	c.EdadMin = edadMin
	c.EdadMax = edadMax

	if err := s.repo.Update(c); err != nil {
		return nil, fmt.Errorf("update categoria: %w", err)
	}
	return c, nil
}

func (s *categoriaService) Delete(id string) error {
	return s.repo.Delete(id)
}

func validate(nombre string, edadMin, edadMax int) error {
	if strings.TrimSpace(nombre) == "" {
		return fmt.Errorf("nombre es requerido")
	}
	if edadMin < 1 || edadMax < 1 {
		return fmt.Errorf("las edades deben ser mayores a 0")
	}
	if edadMin >= edadMax {
		return fmt.Errorf("edad mínima debe ser menor a edad máxima")
	}
	return nil
}
