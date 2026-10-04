package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// InitDB initializes SQLite database connection, applies PRAGMAs and creates tables
func InitDB(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		dbPath = "game_catalog.db"
	}

	// If not an in-memory database, ensure the directory exists
	if dbPath != ":memory:" && !strings.Contains(dbPath, "mode=memory") {
		dir := filepath.Dir(dbPath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory '%s': %w", dir, err)
			}
		}
	}

	// Construct DSN with pragmas for robust SQLite concurrency
	dsn := dbPath
	if dbPath == ":memory:" {
		dsn = "file::memory:?cache=shared"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Optimize connection pool for SQLite
	if dbPath == ":memory:" || strings.Contains(dbPath, "mode=memory") {
		// In-memory databases need to keep at least 1 open connection to avoid losing data
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)
	}

	// Apply SQLite PRAGMAs
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 5000;",
	}
	if dbPath != ":memory:" && !strings.Contains(dbPath, "mode=memory") {
		pragmas = append(pragmas,
			"PRAGMA journal_mode = WAL;",
			"PRAGMA synchronous = NORMAL;",
		)
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			// Some PRAGMAs in in-memory mode might warn or no-op, log or proceed
			_ = err
		}
	}

	// Run migrations
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS games (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		platform TEXT NOT NULL,
		subcategory TEXT,
		genre TEXT NOT NULL DEFAULT 'Gaming',
		time_to_beat TEXT,
		time_to_beat_main REAL,
		added_at TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_games_platform ON games (platform);
	CREATE INDEX IF NOT EXISTS idx_games_subcategory ON games (subcategory);
	CREATE INDEX IF NOT EXISTS idx_games_genre ON games (genre);
	CREATE INDEX IF NOT EXISTS idx_games_added_at ON games (added_at);
	CREATE INDEX IF NOT EXISTS idx_games_title ON games (title);

	-- Unique constraint preventing exact same title on the same platform and subcategory
	CREATE UNIQUE INDEX IF NOT EXISTS idx_games_unique_entry 
	ON games (LOWER(TRIM(title)), platform, COALESCE(subcategory, ''));
	`

	_, err := db.Exec(schema)
	return err
}
