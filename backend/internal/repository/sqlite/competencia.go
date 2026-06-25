package sqlite

import (
	"database/sql"

	"github.com/ema/fixture/backend/internal/repository"
	"github.com/google/uuid"
)

type CompetenciaRepository struct {
	db *sql.DB
}

func NewCompetenciaRepository(db *sql.DB) *CompetenciaRepository {
	return &CompetenciaRepository{db: db}
}

func (r *CompetenciaRepository) Create(c *repository.Competencia) error {
	c.ID = uuid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO competencias (id, categoria_id, numero, nombre, estado, rondas)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		c.ID, c.CategoriaID, c.Numero, c.Nombre, c.Estado, c.Rondas)
	return err
}

func (r *CompetenciaRepository) ListByCategoria(categoriaID string) ([]repository.Competencia, error) {
	rows, err := r.db.Query(
		`SELECT id, categoria_id, numero, nombre, estado, rondas, created_at
		 FROM competencias WHERE categoria_id = ? ORDER BY created_at DESC`, categoriaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comp []repository.Competencia
	for rows.Next() {
		var c repository.Competencia
		if err := rows.Scan(&c.ID, &c.CategoriaID, &c.Numero, &c.Nombre, &c.Estado, &c.Rondas, &c.CreatedAt); err != nil {
			return nil, err
		}
		comp = append(comp, c)
	}
	return comp, rows.Err()
}

func (r *CompetenciaRepository) GetByID(id string) (*repository.Competencia, error) {
	var c repository.Competencia
	err := r.db.QueryRow(
		`SELECT id, categoria_id, numero, nombre, estado, rondas, created_at
		 FROM competencias WHERE id = ?`, id).
		Scan(&c.ID, &c.CategoriaID, &c.Numero, &c.Nombre, &c.Estado, &c.Rondas, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompetenciaRepository) SetEstado(id, estado string) error {
	_, err := r.db.Exec(`UPDATE competencias SET estado = ? WHERE id = ?`, estado, id)
	return err
}

func (r *CompetenciaRepository) SetNombre(id, nombre string) error {
	_, err := r.db.Exec(`UPDATE competencias SET nombre = ? WHERE id = ?`, nombre, id)
	return err
}
