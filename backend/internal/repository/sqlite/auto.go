package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ema/fixture/backend/internal/domain"
)

// AutoRepository implementa repository.AutoRepository con SQLite.
type AutoRepository struct {
	db *sql.DB
}

// NewAutoRepository crea un nuevo repositorio SQLite de autos.
func NewAutoRepository(db *sql.DB) *AutoRepository {
	return &AutoRepository{db: db}
}

func (r *AutoRepository) ListByCategoria(categoriaID string) ([]domain.Auto, error) {
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

func (r *AutoRepository) GetByID(id string) (*domain.Auto, error) {
	var a domain.Auto
	err := r.db.QueryRow(
		`SELECT id, categoria_id, numero, nombre, creador, edad, foto_url, created_at, updated_at
		 FROM autos WHERE id = ?`, id).
		Scan(&a.ID, &a.CategoriaID, &a.Numero, &a.Nombre, &a.Creador, &a.Edad, &a.FotoURL, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get auto: %w", err)
	}
	return &a, nil
}

func (r *AutoRepository) Save(a *domain.Auto) error {
	now := time.Now().UTC().Format(time.RFC3339)
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO autos (id, categoria_id, numero, nombre, creador, edad, foto_url, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.CategoriaID, a.Numero, a.Nombre, a.Creador, a.Edad, a.FotoURL, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert auto: %w", err)
	}
	return nil
}

func (r *AutoRepository) Update(a *domain.Auto) error {
	a.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	res, err := r.db.Exec(
		`UPDATE autos SET numero = ?, nombre = ?, creador = ?, edad = ?, foto_url = ?, updated_at = ?
		 WHERE id = ?`,
		a.Numero, a.Nombre, a.Creador, a.Edad, a.FotoURL, a.UpdatedAt, a.ID)
	if err != nil {
		return fmt.Errorf("update auto: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("auto %s not found", a.ID)
	}
	return nil
}

func (r *AutoRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM autos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete auto: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("auto %s not found", id)
	}
	return nil
}

func (r *AutoRepository) ExistsByNumero(categoriaID string, numero int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM autos WHERE categoria_id = ? AND numero = ?`,
		categoriaID, numero).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check numero exists: %w", err)
	}
	return count > 0, nil
}
