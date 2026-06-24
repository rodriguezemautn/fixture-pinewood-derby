package sqlite

import (
	"database/sql"
	"testing"

	"github.com/ema/fixture/backend/internal/domain"
	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categorias (
		id TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		edad_min INTEGER NOT NULL,
		edad_max INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestCategoriaRepository_SaveAndList(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewCategoriaRepository(db)

	c := &domain.Categoria{
		ID:      "test-1",
		Nombre:  "Pre-Juveniles",
		EdadMin: 10,
		EdadMax: 12,
	}

	if err := repo.Save(c); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	list, err := repo.List()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 categoria, got %d", len(list))
	}
	if list[0].Nombre != "Pre-Juveniles" {
		t.Errorf("expected nombre=Pre-Juveniles, got %s", list[0].Nombre)
	}
}

func TestCategoriaRepository_GetByID(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewCategoriaRepository(db)

	c := &domain.Categoria{ID: "test-2", Nombre: "Juveniles", EdadMin: 13, EdadMax: 15}
	repo.Save(c)

	got, err := repo.GetByID("test-2")
	if err != nil {
		t.Fatalf("getbyid failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected categoria, got nil")
	}
	if got.Nombre != "Juveniles" {
		t.Errorf("expected Juveniles, got %s", got.Nombre)
	}

	notFound, err := repo.GetByID("no-existe")
	if err != nil {
		t.Fatalf("getbyid nonexistent failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for nonexistent categoria")
	}
}

func TestCategoriaRepository_Update(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewCategoriaRepository(db)

	c := &domain.Categoria{ID: "test-3", Nombre: "Viejos", EdadMin: 40, EdadMax: 50}
	repo.Save(c)

	c.Nombre = "Master"
	c.EdadMin = 30
	if err := repo.Update(c); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, _ := repo.GetByID("test-3")
	if got.Nombre != "Master" || got.EdadMin != 30 {
		t.Errorf("update not applied: %+v", got)
	}
}

func TestCategoriaRepository_Delete(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewCategoriaRepository(db)

	c := &domain.Categoria{ID: "test-4", Nombre: "Eliminar", EdadMin: 1, EdadMax: 5}
	repo.Save(c)

	if err := repo.Delete("test-4"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	got, _ := repo.GetByID("test-4")
	if got != nil {
		t.Error("expected nil after delete")
	}
}
