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

func NewAutoRepository(db *sql.DB) *AutoRepository {
	return &AutoRepository{db: db}
}

const autoCols = "id, categoria_id, numero, nombre, creador, edad, peso, foto_url, created_at, updated_at"

func scanAuto(scanner interface {
	Scan(dest ...interface{}) error
}, a *domain.Auto) error {
	return scanner.Scan(&a.ID, &a.CategoriaID, &a.Numero, &a.Nombre, &a.Creador, &a.Edad, &a.Peso, &a.FotoURL, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AutoRepository) ListAll() ([]domain.Auto, error) {
	rows, err := r.db.Query(
		`SELECT ` + autoCols + ` FROM autos ORDER BY categoria_id, numero ASC`)
	if err != nil {
		return nil, fmt.Errorf("query all autos: %w", err)
	}
	defer rows.Close()

	var autos []domain.Auto
	for rows.Next() {
		var a domain.Auto
		if err := scanAuto(rows, &a); err != nil {
			return nil, fmt.Errorf("scan auto: %w", err)
		}
		autos = append(autos, a)
	}
	return autos, rows.Err()
}

func (r *AutoRepository) ListByCategoria(categoriaID string) ([]domain.Auto, error) {
	rows, err := r.db.Query(
		`SELECT `+autoCols+` FROM autos WHERE categoria_id = ? ORDER BY numero ASC`, categoriaID)
	if err != nil {
		return nil, fmt.Errorf("query autos: %w", err)
	}
	defer rows.Close()

	var autos []domain.Auto
	for rows.Next() {
		var a domain.Auto
		if err := scanAuto(rows, &a); err != nil {
			return nil, fmt.Errorf("scan auto: %w", err)
		}
		autos = append(autos, a)
	}
	return autos, rows.Err()
}

func (r *AutoRepository) GetByID(id string) (*domain.Auto, error) {
	var a domain.Auto
	err := scanAuto(r.db.QueryRow(
		`SELECT `+autoCols+` FROM autos WHERE id = ?`, id), &a)
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
		`INSERT INTO autos (id, categoria_id, numero, nombre, creador, edad, peso, foto_url, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.CategoriaID, a.Numero, a.Nombre, a.Creador, a.Edad, a.Peso, a.FotoURL, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert auto: %w", err)
	}
	return nil
}

func (r *AutoRepository) Update(a *domain.Auto) error {
	a.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	res, err := r.db.Exec(
		`UPDATE autos SET numero = ?, nombre = ?, creador = ?, edad = ?, peso = ?, foto_url = ?, updated_at = ?
		 WHERE id = ?`,
		a.Numero, a.Nombre, a.Creador, a.Edad, a.Peso, a.FotoURL, a.UpdatedAt, a.ID)
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
