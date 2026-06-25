package service

import (
	"encoding/json"
	"fmt"

	"github.com/ema/fixture/backend/internal/domain"
	"github.com/ema/fixture/backend/internal/repository"
	"github.com/google/uuid"
)

// FixtureService define operaciones de fixture.
type FixtureService interface {
	Generar(categoriaID string, rondas int) (*domain.Fixture, error)
	ObtenerFixture(categoriaID string) (*domain.Fixture, error)
	RegistrarResultado(heatID string, ordenLlegada []string) error
	ObtenerPosiciones(categoriaID string) ([]domain.Standing, error)
	GenerarFinal(categoriaID string) (*domain.Heat, error)
	Archivar(categoriaID string) error
	GetArchivos(categoriaID string) ([]map[string]any, error)
}

type fixtureService struct {
	fixtureRepo repository.FixtureRepository
}

// NewFixtureService crea un service de fixture.
func NewFixtureService(fixtureRepo repository.FixtureRepository) FixtureService {
	return &fixtureService{fixtureRepo: fixtureRepo}
}

func (s *fixtureService) Generar(categoriaID string, rondas int) (*domain.Fixture, error) {
	// Verificar que no exista un fixture activo
	existente, _ := s.fixtureRepo.GetByCategoria(categoriaID)
	if existente != nil && existente.Estado != "finalizado" {
		return nil, fmt.Errorf("ya existe un fixture activo para esta categoría")
	}

	// Obtener autos
	autoIDs, err := s.fixtureRepo.GetAutoIDsByCategoria(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("obtener autos: %w", err)
	}

	// Generar heats
	heats := GenerarHeats(autoIDs, rondas)

	// Asignar IDs
	fixtureID := uuid.New().String()
	for i := range heats {
		heats[i].ID = uuid.New().String()
	}

	f := &domain.Fixture{
		ID:          fixtureID,
		CategoriaID: categoriaID,
		Estado:      "pendiente",
		Rondas:      rondas,
		Heats:       heats,
	}

	if err := s.fixtureRepo.Create(f); err != nil {
		return nil, fmt.Errorf("guardar fixture: %w", err)
	}

	return f, nil
}

func (s *fixtureService) ObtenerFixture(categoriaID string) (*domain.Fixture, error) {
	return s.fixtureRepo.GetByCategoria(categoriaID)
}

func (s *fixtureService) RegistrarResultado(heatID string, ordenLlegada []string) error {
	if len(ordenLlegada) < 2 {
		return fmt.Errorf("se requiere al menos 2 autos para registrar resultado")
	}
	return s.fixtureRepo.UpdateHeatResult(heatID, ordenLlegada)
}

func (s *fixtureService) ObtenerPosiciones(categoriaID string) ([]domain.Standing, error) {
	f, err := s.fixtureRepo.GetByCategoria(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("obtener fixture: %w", err)
	}
	if f == nil {
		return nil, nil
	}

	autos, err := s.fixtureRepo.GetAllAutos(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("obtener autos: %w", err)
	}

	standings := CalcularStandings(autos, f.Heats)
	return standings, nil
}

func (s *fixtureService) GenerarFinal(categoriaID string) (*domain.Heat, error) {
	f, err := s.fixtureRepo.GetByCategoria(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("obtener fixture: %w", err)
	}
	if f == nil {
		return nil, fmt.Errorf("no hay fixture para esta categoría")
	}

	// Verificar que todos los heats estén completados
	for _, h := range f.Heats {
		if !h.Completado {
			return nil, fmt.Errorf("heat %d no está completado", h.Numero)
		}
	}

	autos, err := s.fixtureRepo.GetAllAutos(categoriaID)
	if err != nil {
		return nil, fmt.Errorf("obtener autos: %w", err)
	}

	standings := CalcularStandings(autos, f.Heats)
	finalistas := SeleccionarFinal(standings)

	final := &domain.Heat{
		ID:      uuid.New().String(),
		Numero:  f.Rondas*((len(finalistas)+3)/4) + 1,
		AutoIDs: finalistas,
		Completado: false,
	}

	if err := s.fixtureRepo.SaveHeat(final, f.ID); err != nil {
		return nil, fmt.Errorf("guardar final: %w", err)
	}

	s.fixtureRepo.SetFixtureEstado(f.ID, "finalizado")
	return final, nil
}

func (s *fixtureService) Archivar(categoriaID string) error {
	// Obtener posiciones finales
	posiciones, err := s.ObtenerPosiciones(categoriaID)
	if err != nil {
		return fmt.Errorf("obtener posiciones: %w", err)
	}
	if len(posiciones) == 0 {
		return fmt.Errorf("no hay posiciones para archivar")
	}

	// Obtener categoría
	f, err := s.fixtureRepo.GetByCategoria(categoriaID)
	if err != nil || f == nil {
		return fmt.Errorf("no hay fixture activo")
	}

	// Obtener nombre de categoría
	autos, err := s.fixtureRepo.GetAllAutos(categoriaID)
	if err != nil {
		return fmt.Errorf("obtener autos: %w", err)
	}

	// Buscar nombre de categoría desde cualquier auto (todos tienen misma categoria_id)
	catNombre := categoriaID
	for _, a := range autos {
		if a.CategoriaID == categoriaID {
			// Intentar obtener el nombre de la categoría
			break
		}
	}

	// Convertir posiciones a JSON
	resultadosB, _ := json.Marshal(posiciones)

	// Guardar archivo
	winner := posiciones[0]
	if err := s.fixtureRepo.ArchivarCompetencia(
		categoriaID, catNombre,
		winner.AutoID, winner.Nombre, winner.Numero,
		string(resultadosB),
	); err != nil {
		return fmt.Errorf("archivar: %w", err)
	}

	// Limpiar fixture para nueva competencia
	return s.fixtureRepo.LimpiarFixture(categoriaID)
}

func (s *fixtureService) GetArchivos(categoriaID string) ([]map[string]any, error) {
	return s.fixtureRepo.GetArchivosByCategoria(categoriaID)
}
