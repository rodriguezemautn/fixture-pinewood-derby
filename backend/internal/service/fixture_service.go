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
	CrearCompetencia(categoriaID string, rondas int) (*repository.Competencia, *domain.Fixture, error)
	ListarCompetencias(categoriaID string) ([]repository.Competencia, error)
	ObtenerFixture(categoriaID string) (*domain.Fixture, error)
	ObtenerFixturePorCompetencia(competenciaID string) (*domain.Fixture, error)
	RegistrarResultado(heatID string, ordenLlegada []string) error
	ObtenerPosiciones(categoriaID string) ([]domain.Standing, error)
	ObtenerPosicionesPorCompetencia(competenciaID string) ([]domain.Standing, error)
	GenerarFinal(competenciaID string) (*domain.Heat, error)
	FinalizarCompetencia(competenciaID string) error
	AgregarDesempate(competenciaID string, autoIDs []string) (*domain.Heat, error)
	Archivar(categoriaID string) error
	GetArchivos(categoriaID string) ([]map[string]any, error)
}

type fixtureService struct {
	fixtureRepo    repository.FixtureRepository
	categoriaRepo  repository.CategoriaRepository
	competenciaRepo repository.CompetenciaRepository
}

func NewFixtureService(fixtureRepo repository.FixtureRepository, categoriaRepo repository.CategoriaRepository, competenciaRepo repository.CompetenciaRepository) FixtureService {
	return &fixtureService{fixtureRepo: fixtureRepo, categoriaRepo: categoriaRepo, competenciaRepo: competenciaRepo}
}

// CrearCompetencia crea una nueva competencia con su fixture.
func (s *fixtureService) CrearCompetencia(categoriaID string, rondas int) (*repository.Competencia, *domain.Fixture, error) {
	// Obtener autos
	autoIDs, err := s.fixtureRepo.GetAutoIDsByCategoria(categoriaID)
	if err != nil {
		return nil, nil, fmt.Errorf("obtener autos: %w", err)
	}

	// Número de competencia
	competencias, _ := s.competenciaRepo.ListByCategoria(categoriaID)
	num := len(competencias) + 1

	cat, _ := s.categoriaRepo.GetByID(categoriaID)
	catNombre := categoriaID
	if cat != nil {
		catNombre = cat.Nombre
	}

	comp := &repository.Competencia{
		CategoriaID: categoriaID,
		Numero:      num,
		Nombre:      fmt.Sprintf("%s - Competencia #%d", catNombre, num),
		Estado:      "abierta",
		Rondas:      rondas,
	}

	if err := s.competenciaRepo.Create(comp); err != nil {
		return nil, nil, fmt.Errorf("crear competencia: %w", err)
	}

	// Generar heats
	heats := GenerarHeats(autoIDs, rondas)
	fixtureID := uuid.New().String()
	for i := range heats {
		heats[i].ID = uuid.New().String()
	}

	f := &domain.Fixture{
		ID:            fixtureID,
		CategoriaID:   categoriaID,
		CompetenciaID: comp.ID,
		Estado:        "pendiente",
		Rondas:        rondas,
		Heats:         heats,
	}

	if err := s.fixtureRepo.Create(f); err != nil {
		return nil, nil, fmt.Errorf("guardar fixture: %w", err)
	}

	return comp, f, nil
}

func (s *fixtureService) ListarCompetencias(categoriaID string) ([]repository.Competencia, error) {
	return s.competenciaRepo.ListByCategoria(categoriaID)
}

func (s *fixtureService) ObtenerFixture(categoriaID string) (*domain.Fixture, error) {
	return s.fixtureRepo.GetByCategoria(categoriaID)
}

func (s *fixtureService) ObtenerFixturePorCompetencia(competenciaID string) (*domain.Fixture, error) {
	return s.fixtureRepo.GetByCompetencia(competenciaID)
}

func (s *fixtureService) RegistrarResultado(heatID string, ordenLlegada []string) error {
	if len(ordenLlegada) < 2 {
		return fmt.Errorf("se requiere al menos 2 autos para registrar resultado")
	}
	return s.fixtureRepo.UpdateHeatResult(heatID, ordenLlegada)
}

func (s *fixtureService) ObtenerPosiciones(categoriaID string) ([]domain.Standing, error) {
	f, err := s.fixtureRepo.GetByCategoria(categoriaID)
	if err != nil || f == nil {
		return nil, err
	}
	autos, err := s.fixtureRepo.GetAllAutos(categoriaID)
	if err != nil {
		return nil, err
	}
	return CalcularStandings(autos, f.Heats), nil
}

func (s *fixtureService) ObtenerPosicionesPorCompetencia(competenciaID string) ([]domain.Standing, error) {
	f, err := s.fixtureRepo.GetByCompetencia(competenciaID)
	if err != nil || f == nil {
		return nil, err
	}
	autos, err := s.fixtureRepo.GetAllAutos(f.CategoriaID)
	if err != nil {
		return nil, err
	}
	return CalcularStandings(autos, f.Heats), nil
}

func (s *fixtureService) GenerarFinal(competenciaID string) (*domain.Heat, error) {
	comp, err := s.competenciaRepo.GetByID(competenciaID)
	if err != nil || comp == nil {
		return nil, fmt.Errorf("competencia no encontrada")
	}
	if comp.Estado == "finalizada" {
		return nil, fmt.Errorf("competencia ya finalizada")
	}

	f, err := s.fixtureRepo.GetByCompetencia(competenciaID)
	if err != nil || f == nil {
		return nil, fmt.Errorf("no hay fixture para esta competencia")
	}

	for _, h := range f.Heats {
		if !h.Completado {
			return nil, fmt.Errorf("heat %d no está completado", h.Numero)
		}
	}

	autos, err := s.fixtureRepo.GetAllAutos(comp.CategoriaID)
	if err != nil {
		return nil, err
	}

	standings := CalcularStandings(autos, f.Heats)
	finalistas := SeleccionarFinal(standings)

	final := &domain.Heat{
		ID:      uuid.New().String(),
		Numero:  f.Rondas*((len(finalistas)+3)/4) + 1,
		AutoIDs: finalistas,
	}

	if err := s.fixtureRepo.SaveHeat(final, f.ID); err != nil {
		return nil, fmt.Errorf("guardar final: %w", err)
	}

	s.fixtureRepo.SetFixtureEstado(f.ID, "finalizado")
	return final, nil
}

func (s *fixtureService) FinalizarCompetencia(competenciaID string) error {
	comp, err := s.competenciaRepo.GetByID(competenciaID)
	if err != nil || comp == nil {
		return fmt.Errorf("competencia no encontrada")
	}
	return s.competenciaRepo.SetEstado(competenciaID, "finalizada")
}

func (s *fixtureService) AgregarDesempate(competenciaID string, autoIDs []string) (*domain.Heat, error) {
	comp, err := s.competenciaRepo.GetByID(competenciaID)
	if err != nil || comp == nil {
		return nil, fmt.Errorf("competencia no encontrada")
	}
	if comp.Estado != "finalizada" {
		return nil, fmt.Errorf("solo se pueden agregar desempates a competencias finalizadas")
	}
	if len(autoIDs) < 2 {
		return nil, fmt.Errorf("se necesitan al menos 2 autos para el desempate")
	}

	f, err := s.fixtureRepo.GetByCompetencia(competenciaID)
	if err != nil || f == nil {
		return nil, fmt.Errorf("no hay fixture")
	}

	heat := &domain.Heat{
		ID:      uuid.New().String(),
		Numero:  len(f.Heats) + 1,
		AutoIDs: autoIDs,
	}

	if err := s.fixtureRepo.SaveHeat(heat, f.ID); err != nil {
		return nil, fmt.Errorf("guardar desempate: %w", err)
	}
	return heat, nil
}

func (s *fixtureService) Archivar(categoriaID string) error {
	posiciones, err := s.ObtenerPosiciones(categoriaID)
	if err != nil {
		return err
	}
	if len(posiciones) == 0 {
		return fmt.Errorf("no hay posiciones para archivar")
	}

	catNombre := categoriaID
	if cat, err := s.categoriaRepo.GetByID(categoriaID); err == nil && cat != nil {
		catNombre = cat.Nombre
	}

	resultadosB, _ := json.Marshal(posiciones)
	winner := posiciones[0]
	if err := s.fixtureRepo.ArchivarCompetencia(
		categoriaID, catNombre,
		winner.AutoID, winner.Nombre, winner.Numero,
		string(resultadosB),
	); err != nil {
		return fmt.Errorf("archivar: %w", err)
	}
	return s.fixtureRepo.LimpiarFixture(categoriaID)
}

func (s *fixtureService) GetArchivos(categoriaID string) ([]map[string]any, error) {
	return s.fixtureRepo.GetArchivosByCategoria(categoriaID)
}
