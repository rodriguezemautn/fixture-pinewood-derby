package sqlite

import (
	"database/sql"
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
	_ "modernc.org/sqlite"
)

func setupAutoDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categorias (
			id TEXT PRIMARY KEY, nombre TEXT NOT NULL,
			edad_min INTEGER NOT NULL, edad_max INTEGER NOT NULL,
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS autos (
			id TEXT PRIMARY KEY, categoria_id TEXT NOT NULL,
			numero INTEGER NOT NULL, nombre TEXT NOT NULL,
			creador TEXT NOT NULL, edad INTEGER NOT NULL,
			foto_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
			UNIQUE(categoria_id, numero),
			FOREIGN KEY (categoria_id) REFERENCES categorias(id)
		);
		INSERT INTO categorias (id, nombre, edad_min, edad_max, created_at, updated_at)
		VALUES ('cat-1', 'Test', 1, 10, 'now', 'now');
	`)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func TestAutoRepository_SaveAndListByCategoria(t *testing.T) {
	db := setupAutoDB(t)
	defer db.Close()
	repo := NewAutoRepository(db)

	a := &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "Turbo", Creador: "Juan", Edad: 10}
	if err := repo.Save(a); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	list, err := repo.ListByCategoria("cat-1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 auto, got %d", len(list))
	}
}

func TestAutoRepository_ExistsByNumero(t *testing.T) {
	db := setupAutoDB(t)
	defer db.Close()
	repo := NewAutoRepository(db)

	repo.Save(&domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 5, Nombre: "Test", Creador: "J", Edad: 8})

	exists, err := repo.ExistsByNumero("cat-1", 5)
	if err != nil {
		t.Fatalf("exists check failed: %v", err)
	}
	if !exists {
		t.Error("expected exists=true for numero 5")
	}

	exists, _ = repo.ExistsByNumero("cat-1", 99)
	if exists {
		t.Error("expected exists=false for numero 99")
	}
}

func TestAutoRepository_Update(t *testing.T) {
	db := setupAutoDB(t)
	defer db.Close()
	repo := NewAutoRepository(db)

	repo.Save(&domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "Old", Creador: "J", Edad: 8})
	repo.Save(&domain.Auto{ID: "a2", CategoriaID: "cat-1", Numero: 2, Nombre: "Test", Creador: "J", Edad: 8})

	// Update numero to avoid unique conflict
	updated := &domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 10, Nombre: "New", Creador: "J", Edad: 10, FotoURL: "foto.jpg"}
	if err := repo.Update(updated); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, _ := repo.GetByID("a1")
	if got.Nombre != "New" || got.Numero != 10 {
		t.Errorf("update not applied: %+v", got)
	}
}

func TestAutoRepository_Delete(t *testing.T) {
	db := setupAutoDB(t)
	defer db.Close()
	repo := NewAutoRepository(db)

	repo.Save(&domain.Auto{ID: "a1", CategoriaID: "cat-1", Numero: 1, Nombre: "Test", Creador: "J", Edad: 8})
	if err := repo.Delete("a1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	got, _ := repo.GetByID("a1")
	if got != nil {
		t.Error("expected nil after delete")
	}
}

func TestAutoRepository_GetByID_NotFound(t *testing.T) {
	db := setupAutoDB(t)
	defer db.Close()
	repo := NewAutoRepository(db)

	got, err := repo.GetByID("no-existe")
	if err != nil {
		t.Fatalf("getbyid failed: %v", err)
	}
	if got != nil {
		t.Error("expected nil for nonexistent")
	}
}
