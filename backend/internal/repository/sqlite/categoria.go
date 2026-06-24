// Package sqlite contiene implementaciones concretas de repositorios con SQLite.
package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ema/fixture/backend/internal/domain"
)

// CategoriaRepository implementa repository.CategoriaRepository con SQLite.
type CategoriaRepository struct {
	db *sql.DB
}

// NewCategoriaRepository crea un nuevo repositorio SQLite de categorías.
func NewCategoriaRepository(db *sql.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

// List retorna todas las categorías ordenadas por edad mínima.
func (r *CategoriaRepository) List() ([]domain.Categoria, error) {
	rows, err := r.db.Query(
		`SELECT id, nombre, edad_min, edad_max, created_at, updated_at
		 FROM categorias ORDER BY edad_min ASC`)
	if err != nil {
		return nil, fmt.Errorf("query categorias: %w", err)
	}
	defer rows.Close()

	var categorias []domain.Categoria
	for rows.Next() {
		var c domain.Categoria
		if err := rows.Scan(&c.ID, &c.Nombre, &c.EdadMin, &c.EdadMax, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan categoria: %w", err)
		}
		categorias = append(categorias, c)
	}
	return categorias, rows.Err()
}

// GetByID retorna una categoría por su ID.
func (r *CategoriaRepository) GetByID(id string) (*domain.Categoria, error) {
	var c domain.Categoria
	err := r.db.QueryRow(
		`SELECT id, nombre, edad_min, edad_max, created_at, updated_at
		 FROM categorias WHERE id = ?`, id).
		Scan(&c.ID, &c.Nombre, &c.EdadMin, &c.EdadMax, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get categoria %s: %w", id, err)
	}
	return &c, nil
}

// Save inserta una nueva categoría.
func (r *CategoriaRepository) Save(c *domain.Categoria) error {
	now := time.Now().UTC().Format(time.RFC3339)
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO categorias (id, nombre, edad_min, edad_max, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		c.ID, c.Nombre, c.EdadMin, c.EdadMax, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert categoria: %w", err)
	}
	return nil
}

// Update actualiza una categoría existente.
func (r *CategoriaRepository) Update(c *domain.Categoria) error {
	c.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	res, err := r.db.Exec(
		`UPDATE categorias SET nombre = ?, edad_min = ?, edad_max = ?, updated_at = ?
		 WHERE id = ?`,
		c.Nombre, c.EdadMin, c.EdadMax, c.UpdatedAt, c.ID)
	if err != nil {
		return fmt.Errorf("update categoria: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("categoria %s not found", c.ID)
	}
	return nil
}

// Delete elimina una categoría por ID.
func (r *CategoriaRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM categorias WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete categoria: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("categoria %s not found", id)
	}
	return nil
}
