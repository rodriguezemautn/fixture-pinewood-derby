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
		`CREATE TABLE IF NOT EXISTS autos (
			id TEXT PRIMARY KEY,
			categoria_id TEXT NOT NULL,
			numero INTEGER NOT NULL,
			nombre TEXT NOT NULL,
			creador TEXT NOT NULL,
			edad INTEGER NOT NULL,
			peso INTEGER NOT NULL DEFAULT 0,
			foto_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE(categoria_id, numero),
			FOREIGN KEY (categoria_id) REFERENCES categorias(id)
		)`,
		`ALTER TABLE autos ADD COLUMN peso INTEGER NOT NULL DEFAULT 0`,
		`CREATE TABLE IF NOT EXISTS fixtures (
			id TEXT PRIMARY KEY,
			categoria_id TEXT NOT NULL,
			rondas INTEGER NOT NULL DEFAULT 3,
			estado TEXT NOT NULL DEFAULT 'pendiente',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (categoria_id) REFERENCES categorias(id)
		)`,
		`CREATE TABLE IF NOT EXISTS heats (
			id TEXT PRIMARY KEY,
			fixture_id TEXT NOT NULL,
			numero INTEGER NOT NULL,
			completado INTEGER NOT NULL DEFAULT 0,
			orden_llegada TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (fixture_id) REFERENCES fixtures(id)
		)`,
		`CREATE TABLE IF NOT EXISTS heat_autos (
			heat_id TEXT NOT NULL,
			auto_id TEXT NOT NULL,
			PRIMARY KEY (heat_id, auto_id),
			FOREIGN KEY (heat_id) REFERENCES heats(id),
			FOREIGN KEY (auto_id) REFERENCES autos(id)
		)`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			// ALTER TABLE puede fallar si la columna ya existe
			if len(m) >= 11 && m[:11] == "ALTER TABLE" {
				continue
			}
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("✅ Migraciones ejecutadas")
	return nil
}
