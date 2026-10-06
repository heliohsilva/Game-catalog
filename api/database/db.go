package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"game-catalog/api/models"

	_ "modernc.org/sqlite"
)

// InitDB initializes SQLite database connection, applies PRAGMAs and creates tables
func InitDB(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		dbPath = "game_catalog.sqlite"
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

	dotSQLPath := "game_catalog.sql"

	if err := loadFrom(db, dotSQLPath); err != nil{
		fmt.Printf("Not possible to load db from .sql file: %v", err)

		if err := runMigrations(db); err != nil{
			return nil, fmt.Errorf("Cannot perform migrations: %w", err)
		}
	}

	return db, nil
}

func loadFrom(db *sql.DB, path string) error {
	dbBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("File not found")
	}

	_, err = db.Exec(string(dbBytes))

	return err
}

func runMigrations(db *sql.DB) error {
	// Run migrations
	if err := migrate(db); err != nil {
		db.Close()
		return fmt.Errorf("database migration failed: %w", err)
	}

	if err := seedDefaultPlatforms(db); err != nil {
		db.Close()
		return fmt.Errorf("default platform seeding failed: %w", err)
	}

	if err := SyncDynamicPlatforms(db); err != nil {
		db.Close()
		return fmt.Errorf("platform sync failed: %w", err)
	}

	return nil
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

	-- Platforms table
	CREATE TABLE IF NOT EXISTS platforms (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT UNIQUE NOT NULL COLLATE NOCASE,
		created_at TEXT NOT NULL
	);

	-- Platform subcategories table
	CREATE TABLE IF NOT EXISTS platform_subcategories (
		id TEXT PRIMARY KEY NOT NULL,
		platform_name TEXT NOT NULL COLLATE NOCASE,
		name TEXT NOT NULL COLLATE NOCASE,
		created_at TEXT NOT NULL,
		UNIQUE(platform_name, name),
		FOREIGN KEY (platform_name) REFERENCES platforms(name) ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_platform_subcategories_platform ON platform_subcategories (platform_name);
	`

	_, err := db.Exec(schema)
	return err
}

func seedDefaultPlatforms(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM platforms").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	defaults := []struct {
		Name string
		Subs []string
	}{
		{models.PlatformPC, []string{"Steam", "GOG", "Epic"}},
		{models.PlatformPlayStation, []string{"PS2", "PS3", "PS4", "PS5"}},
		{models.PlatformNintendoSwitch, []string{"Switch"}},
		{models.PlatformXbox, []string{"Xbox Series X/S", "Xbox One", "Xbox 360"}},
		{models.PlatformRetroEmulation, []string{"PS1", "N64", "SNES", "Genesis", "NES", "Master System", "Neo Geo"}},
	}

	now := time.Now().UTC().Format(time.RFC3339)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, d := range defaults {
		platID := fmt.Sprintf("plat-default-%d", i+1)
		_, err := tx.Exec("INSERT INTO platforms (id, name, created_at) VALUES (?, ?, ?)", platID, d.Name, now)
		if err != nil {
			return err
		}
		for j, s := range d.Subs {
			subID := fmt.Sprintf("sub-default-%d-%d", i+1, j+1)
			_, err := tx.Exec("INSERT INTO platform_subcategories (id, platform_name, name, created_at) VALUES (?, ?, ?, ?)",
				subID, d.Name, s, now)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func SyncDynamicPlatforms(db *sql.DB) error {
	rows, err := db.Query("SELECT p.name, COALESCE(s.name, '') FROM platforms p LEFT JOIN platform_subcategories s ON s.platform_name = p.name ORDER BY p.rowid ASC, s.rowid ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	platformMap := make(map[string][]string)
	var platformOrder []string

	for rows.Next() {
		var pName, sName string
		if err := rows.Scan(&pName, &sName); err != nil {
			return err
		}
		if _, exists := platformMap[pName]; !exists {
			platformMap[pName] = []string{}
			platformOrder = append(platformOrder, pName)
		}
		if sName != "" {
			platformMap[pName] = append(platformMap[pName], sName)
		}
	}

	platformList := make([]models.PlatformInfo, 0, len(platformOrder))
	for _, name := range platformOrder {
		platformList = append(platformList, models.PlatformInfo{
			Name:          name,
			Subcategories: platformMap[name],
		})
	}

	models.SetDynamicPlatforms(platformList)
	return nil
}
