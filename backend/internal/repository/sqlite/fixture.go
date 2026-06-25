package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ema/fixture/backend/internal/domain"
)

// FixtureRepository implementa repository.FixtureRepository con SQLite.
type FixtureRepository struct {
	db *sql.DB
}

func NewFixtureRepository(db *sql.DB) *FixtureRepository {
	return &FixtureRepository{db: db}
}

func (r *FixtureRepository) Create(f *domain.Fixture) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO fixtures (id, categoria_id, rondas, estado, competencia_id) VALUES (?, ?, ?, ?, ?)`,
		f.ID, f.CategoriaID, f.Rondas, f.Estado, f.CompetenciaID)
	if err != nil {
		return fmt.Errorf("insert fixture: %w", err)
	}

	for _, heat := range f.Heats {
		_, err = tx.Exec(`INSERT INTO heats (id, fixture_id, numero) VALUES (?, ?, ?)`,
			heat.ID, f.ID, heat.Numero)
		if err != nil {
			return fmt.Errorf("insert heat %d: %w", heat.Numero, err)
		}
		for _, autoID := range heat.AutoIDs {
			_, err = tx.Exec(`INSERT INTO heat_autos (heat_id, auto_id) VALUES (?, ?)`, heat.ID, autoID)
			if err != nil {
				return fmt.Errorf("insert heat_auto: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *FixtureRepository) GetByCompetencia(competenciaID string) (*domain.Fixture, error) {
	row := r.db.QueryRow(`SELECT id, categoria_id, rondas, estado FROM fixtures WHERE competencia_id = ?`, competenciaID)
	var f domain.Fixture
	err := row.Scan(&f.ID, &f.CategoriaID, &f.Rondas, &f.Estado)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get fixture by competencia: %w", err)
	}
	heats, err := r.GetHeatsByFixture(f.ID)
	if err != nil {
		return nil, err
	}
	f.Heats = heats
	return &f, nil
}

func (r *FixtureRepository) GetByCategoria(categoriaID string) (*domain.Fixture, error) {
	row := r.db.QueryRow(`SELECT id, categoria_id, rondas, estado FROM fixtures WHERE categoria_id = ?`, categoriaID)
	var f domain.Fixture
	err := row.Scan(&f.ID, &f.CategoriaID, &f.Rondas, &f.Estado)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get fixture: %w", err)
	}

	heats, err := r.GetHeatsByFixture(f.ID)
	if err != nil {
		return nil, err
	}
	f.Heats = heats

	return &f, nil
}

func (r *FixtureRepository) SaveHeat(heat *domain.Heat, fixtureID string) error {
	_, err := r.db.Exec(`INSERT INTO heats (id, fixture_id, numero) VALUES (?, ?, ?)`,
		heat.ID, fixtureID, heat.Numero)
	if err != nil {
		return fmt.Errorf("insert heat: %w", err)
	}
	for _, autoID := range heat.AutoIDs {
		_, err = r.db.Exec(`INSERT INTO heat_autos (heat_id, auto_id) VALUES (?, ?)`, heat.ID, autoID)
		if err != nil {
			return fmt.Errorf("insert heat_auto: %w", err)
		}
	}
	return nil
}

func (r *FixtureRepository) UpdateHeatResult(heatID string, ordenLlegada []string) error {
	b, _ := json.Marshal(ordenLlegada)
	res, err := r.db.Exec(
		`UPDATE heats SET completado = 1, orden_llegada = ?, registrado_at = datetime('now') WHERE id = ?`,
		string(b), heatID)
	if err != nil {
		return fmt.Errorf("update heat: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("heat %s not found", heatID)
	}
	return nil
}

func (r *FixtureRepository) GetHeatsByFixture(fixtureID string) ([]domain.Heat, error) {
	rows, err := r.db.Query(
		`SELECT id, numero, completado, orden_llegada FROM heats WHERE fixture_id = ? ORDER BY numero ASC`, fixtureID)
	if err != nil {
		return nil, fmt.Errorf("query heats: %w", err)
	}
	defer rows.Close()

	var heats []domain.Heat
	for rows.Next() {
		var h domain.Heat
		var completado int
		var ordenJSON string
		if err := rows.Scan(&h.ID, &h.Numero, &completado, &ordenJSON); err != nil {
			return nil, fmt.Errorf("scan heat: %w", err)
		}
		h.Completado = completado == 1
		if ordenJSON != "" {
			json.Unmarshal([]byte(ordenJSON), &h.OrdenLlegada)
		}

		// Cargar autos de la heat
		autoRows, err := r.db.Query(`SELECT auto_id FROM heat_autos WHERE heat_id = ? ORDER BY rowid`, h.ID)
		if err != nil {
			return nil, fmt.Errorf("query heat_autos: %w", err)
		}
		for autoRows.Next() {
			var autoID string
			autoRows.Scan(&autoID)
			h.AutoIDs = append(h.AutoIDs, autoID)
		}
		autoRows.Close()

		heats = append(heats, h)
	}
	return heats, rows.Err()
}

func (r *FixtureRepository) GetAutosByCategoria(categoriaID string) ([]domain.Auto, error) {
	rows, err := r.db.Query(
		`SELECT id, categoria_id, numero, nombre, creador, edad, foto_url, created_at, updated_at
		 FROM autos WHERE categoria_id = ? ORDER BY numero ASC`, categoriaID)
	if err != nil {
		return nil, fmt.Errorf("query autos: %w", err)
	}
	defer rows.Close()

	var autos []domain.Auto
	for rows.Next() {
		var a domain.Auto
		if err := rows.Scan(&a.ID, &a.CategoriaID, &a.Numero, &a.Nombre, &a.Creador, &a.Edad, &a.FotoURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan auto: %w", err)
		}
		autos = append(autos, a)
	}
	return autos, rows.Err()
}

func (r *FixtureRepository) SetFixtureEstado(fixtureID, estado string) error {
	_, err := r.db.Exec(`UPDATE fixtures SET estado = ? WHERE id = ?`, estado, fixtureID)
	return err
}

// SetFinal registra el ID de la carrera final en el fixture (almacenado en estado "finalizado").
func (r *FixtureRepository) SetFinal(fixtureID, finalHeatID string) error {
	_, err := r.db.Exec(`UPDATE heats SET completado = 1 WHERE id = ?`, finalHeatID)
	return err
}

// Helper para construir lista de autos de una categoría como slice de strings.
func (r *FixtureRepository) GetAutoIDsByCategoria(categoriaID string) ([]string, error) {
	rows, err := r.db.Query(`SELECT id FROM autos WHERE categoria_id = ? ORDER BY numero`, categoriaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	// Asegurar mínimo 4 autos para validación
	if len(ids) < 4 {
		return ids, fmt.Errorf("se necesitan al menos 4 autos, hay %d", len(ids))
	}
	return ids, rows.Err()
}

// GetHeatByID retorna un heat individual.
func (r *FixtureRepository) GetHeatByID(heatID string) (*domain.Heat, error) {
	var h domain.Heat
	var completado int
	var ordenJSON string
	err := r.db.QueryRow(
		`SELECT id, numero, completado, orden_llegada FROM heats WHERE id = ?`, heatID).
		Scan(&h.ID, &h.Numero, &completado, &ordenJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get heat: %w", err)
	}
	h.Completado = completado == 1
	if ordenJSON != "" {
		json.Unmarshal([]byte(ordenJSON), &h.OrdenLlegada)
	}

	autoRows, err := r.db.Query(`SELECT auto_id FROM heat_autos WHERE heat_id = ? ORDER BY rowid`, h.ID)
	if err != nil {
		return nil, err
	}
	defer autoRows.Close()
	for autoRows.Next() {
		var autoID string
		autoRows.Scan(&autoID)
		h.AutoIDs = append(h.AutoIDs, autoID)
	}

	return &h, nil
}

// GetHeatIDByNumero busca el heat por fixture_id + numero.
func (r *FixtureRepository) GetHeatIDByNumero(fixtureID string, numero int) (string, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM heats WHERE fixture_id = ? AND numero = ?`, fixtureID, numero).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetAllAutos retorna todos los autos filtrados por categoria.
func (r *FixtureRepository) GetAllAutos(categoriaIDs ...string) (map[string]*domain.Auto, error) {
	query := `SELECT id, categoria_id, numero, nombre, creador, edad, peso, foto_url, created_at, updated_at FROM autos`
	var args []interface{}
	if len(categoriaIDs) > 0 {
		placeholders := make([]string, len(categoriaIDs))
		for i, id := range categoriaIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += ` WHERE categoria_id IN (` + strings.Join(placeholders, ",") + `)`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	autos := make(map[string]*domain.Auto)
	for rows.Next() {
		var a domain.Auto
		if err := rows.Scan(&a.ID, &a.CategoriaID, &a.Numero, &a.Nombre, &a.Creador, &a.Edad, &a.Peso, &a.FotoURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		autos[a.ID] = &a
	}
	return autos, rows.Err()
}

// ─── Archivado ────────────────────────────────────────────────

// ArchivarCompetencia guarda los resultados finales y resetea el fixture.
func (r *FixtureRepository) ArchivarCompetencia(categoriaID, categoriaNombre, winnerID, winnerNombre string, winnerNumero int, resultadosJSON string) error {
	id := uuid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO archivos_carrera (id, categoria_id, categoria_nombre, fecha, winner_id, winner_nombre, winner_numero, resultados)
		 VALUES (?, ?, ?, datetime('now'), ?, ?, ?, ?)`,
		id, categoriaID, categoriaNombre, winnerID, winnerNombre, winnerNumero, resultadosJSON)
	return err
}

// GetArchivosByCategoria retorna los archivos de una categoría.
func (r *FixtureRepository) GetArchivosByCategoria(categoriaID string) ([]map[string]any, error) {
	rows, err := r.db.Query(
		`SELECT id, categoria_nombre, fecha, winner_nombre, winner_numero FROM archivos_carrera
		 WHERE categoria_id = ? ORDER BY created_at DESC`, categoriaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var archivos []map[string]any
	for rows.Next() {
		var id, nombre, fecha, wNombre string
		var wNumero int
		if err := rows.Scan(&id, &nombre, &fecha, &wNombre, &wNumero); err != nil {
			return nil, err
		}
		archivos = append(archivos, map[string]any{
			"id": id, "categoria_nombre": nombre, "fecha": fecha,
			"winner_nombre": wNombre, "winner_numero": wNumero,
		})
	}
	return archivos, nil
}

// LimpiarFixture elimina heats, heat_autos y fixture de una categoría.
func (r *FixtureRepository) LimpiarFixture(categoriaID string) error {
	f, err := r.GetByCategoria(categoriaID)
	if err != nil || f == nil {
		return err
	}
	// Borrar heat_autos, heats, y fixture
	r.db.Exec(`DELETE FROM heat_autos WHERE heat_id IN (SELECT id FROM heats WHERE fixture_id = ?)`, f.ID)
	r.db.Exec(`DELETE FROM heats WHERE fixture_id = ?`, f.ID)
	r.db.Exec(`DELETE FROM fixtures WHERE id = ?`, f.ID)
	return nil
}
