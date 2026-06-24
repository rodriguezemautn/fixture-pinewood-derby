// Package database maneja la conexión a SQLite y migraciones.
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// New abre conexión SQLite y ejecuta migraciones.
func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// WAL mode para mejor concurrencia
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("error setting WAL mode: %w", err)
	}

	if err := Migrate(db); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("🗄️  Base de datos inicializada")
	return db, nil
}

// Migrate ejecuta migraciones idempotentes.
func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS categorias (
			id TEXT PRIMARY KEY,
			nombre TEXT NOT NULL,
			edad_min INTEGER NOT NULL,
			edad_max INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now'))
		)`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("✅ Migraciones ejecutadas")
	return nil
}
