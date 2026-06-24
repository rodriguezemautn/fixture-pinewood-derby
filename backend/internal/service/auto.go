package service

import (
	"fmt"
	"strings"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/repository"
	"github.com/google/uuid"
)

type autoService struct {
	autoRepo      repository.AutoRepository
	categoriaRepo repository.CategoriaRepository
}

// NewAutoService crea un service concreto de autos.
func NewAutoService(autoRepo repository.AutoRepository, categoriaRepo repository.CategoriaRepository) AutoService {
	return &autoService{autoRepo: autoRepo, categoriaRepo: categoriaRepo}
}

func (s *autoService) ListByCategoria(categoriaID string) ([]domain.Auto, error) {
	return s.autoRepo.ListByCategoria(categoriaID)
}

func (s *autoService) ListAll() ([]domain.Auto, error) {
	return s.autoRepo.ListAll()
}

func (s *autoService) GetByID(id string) (*domain.Auto, error) {
	return s.autoRepo.GetByID(id)
}

func (s *autoService) Create(categoriaID string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error) {
	cat, err := s.categoriaRepo.GetByID(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("verificar categoria: %w", err)
	}
	if cat == nil {
		return nil, fmt.Errorf("categoria %s no encontrada", categoriaID)
	}

	if err := validateAuto(numero, nombre, creador, edad, peso); err != nil {
		return nil, err
	}

	exists, err := s.autoRepo.ExistsByNumero(categoriaID, numero)
	if err != nil {
		return nil, fmt.Errorf("verificar numero: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("ya existe un auto con número %d en esta categoría", numero)
	}

	a := &domain.Auto{
		ID:          uuid.New().String(),
		CategoriaID: categoriaID,
		Numero:      numero,
		Nombre:      strings.TrimSpace(nombre),
		Creador:     strings.TrimSpace(creador),
		Edad:        edad,
		Peso:        peso,
		FotoURL:     fotoURL,
	}

	if err := s.autoRepo.Save(a); err != nil {
		return nil, fmt.Errorf("save auto: %w", err)
	}
	return a, nil
}

func (s *autoService) Update(id string, numero int, nombre, creador string, edad int, peso int, fotoURL string) (*domain.Auto, error) {
	if err := validateAuto(numero, nombre, creador, edad, peso); err != nil {
		return nil, err
	}

	a, err := s.autoRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get auto for update: %w", err)
	}
	if a == nil {
		return nil, fmt.Errorf("auto %s not found", id)
	}

	if a.Numero != numero {
		exists, err := s.autoRepo.ExistsByNumero(a.CategoriaID, numero)
		if err != nil {
			return nil, fmt.Errorf("verificar numero: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("ya existe un auto con número %d en esta categoría", numero)
		}
	}

	a.Numero = numero
	a.Nombre = strings.TrimSpace(nombre)
	a.Creador = strings.TrimSpace(creador)
	a.Edad = edad
	a.Peso = peso
	a.FotoURL = fotoURL

	if err := s.autoRepo.Update(a); err != nil {
		return nil, fmt.Errorf("update auto: %w", err)
	}
	return a, nil
}

func validateAuto(numero int, nombre, creador string, edad int, peso int) error {
	if strings.TrimSpace(nombre) == "" {
		return fmt.Errorf("nombre es requerido")
	}
	if strings.TrimSpace(creador) == "" {
		return fmt.Errorf("creador es requerido")
	}
	if numero < 1 {
		return fmt.Errorf("número debe ser mayor a 0")
	}
	if edad < 1 {
		return fmt.Errorf("edad debe ser mayor a 0")
	}
	if peso < 0 {
		return fmt.Errorf("peso no puede ser negativo")
	}
	if peso > 5000 {
		return fmt.Errorf("peso no puede superar 5000g (5kg)")
	}
	return nil
}

func (s *autoService) Delete(id string) error {
	return s.autoRepo.Delete(id)
}
