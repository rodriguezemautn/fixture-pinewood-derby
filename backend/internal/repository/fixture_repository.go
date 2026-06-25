package repository

import "github.com/ema/fixture/backend/internal/domain"

// FixtureRepository define operaciones sobre fixtures.
type FixtureRepository interface {
	Create(f *domain.Fixture) error
	GetByCategoria(categoriaID string) (*domain.Fixture, error)
	GetByCompetencia(competenciaID string) (*domain.Fixture, error)
	SaveHeat(heat *domain.Heat, fixtureID string) error
	UpdateHeatResult(heatID string, ordenLlegada []string) error
	GetHeatsByFixture(fixtureID string) ([]domain.Heat, error)
	GetAutosByCategoria(categoriaID string) ([]domain.Auto, error)
	SetFixtureEstado(fixtureID, estado string) error
	GetAutoIDsByCategoria(categoriaID string) ([]string, error)
	GetAllAutos(categoriaIDs ...string) (map[string]*domain.Auto, error)
	ArchivarCompetencia(categoriaID, categoriaNombre, winnerID, winnerNombre string, winnerNumero int, resultadosJSON string) error
	GetArchivosByCategoria(categoriaID string) ([]map[string]any, error)
	LimpiarFixture(fixtureID string) error
}
