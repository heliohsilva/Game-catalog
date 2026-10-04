package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDB_InMemory(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	// Verify games table exists
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='games'").Scan(&name)
	if err != nil || name != "games" {
		t.Fatalf("expected table 'games' to exist, got error: %v", err)
	}

	// Verify indexes exist
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index'")
	if err != nil {
		t.Fatalf("failed to query indexes: %v", err)
	}
	defer rows.Close()

	indexes := make(map[string]bool)
	for rows.Next() {
		var idxName string
		if err := rows.Scan(&idxName); err == nil {
			indexes[idxName] = true
		}
	}

	requiredIndexes := []string{
		"idx_games_platform",
		"idx_games_genre",
		"idx_games_title",
		"idx_games_unique_entry",
	}
	for _, req := range requiredIndexes {
		if !indexes[req] {
			t.Errorf("missing required index: %s", req)
		}
	}
}

func TestInitDB_DiskFileWithNestedDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "game_catalog_db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "nested", "sub", "test.db")
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db with nested dir: %v", err)
	}
	defer db.Close()

	// Verify file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("expected db file to exist at %s", dbPath)
	}

	// Seed database
	inserted, err := SeedDatabase(db, false)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	if inserted != len(InitialSeedGames) {
		t.Errorf("expected %d initial games seeded, got %d", len(InitialSeedGames), inserted)
	}

	// Re-seeding without reset is idempotent
	reinserted, err := SeedDatabase(db, false)
	if err != nil {
		t.Fatalf("second seed failed: %v", err)
	}
	if reinserted != 0 {
		t.Errorf("expected 0 new games on duplicate seed, got %d", reinserted)
	}

	// Re-seeding with reset=true resets and reinserts
	resetCount, err := SeedDatabase(db, true)
	if err != nil {
		t.Fatalf("reset seed failed: %v", err)
	}
	if resetCount != len(InitialSeedGames) {
		t.Errorf("expected %d games after reset seed, got %d", len(InitialSeedGames), resetCount)
	}
}

func TestInitDB_DefaultPath(t *testing.T) {
	// Empty path defaults to game_catalog.sqlite
	tmpDir, err := os.MkdirTemp("", "game_catalog_default_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(cwd)

	db, err := InitDB("")
	if err != nil {
		t.Fatalf("failed to init db with empty path: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat("game_catalog.sqlite"); os.IsNotExist(err) {
		t.Error("expected default game_catalog.sqlite to be created")
	}
}
