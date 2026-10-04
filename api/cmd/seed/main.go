package main

import (
	"flag"
	"log"
	"os"

	"game-catalog/api/database"
)

func main() {
	dbPath := flag.String("path", "data/game_catalog.sqlite", "Path to SQLite database file")
	reset := flag.Bool("reset", true, "Reset and reseed database")
	flag.Parse()

	// Also check environment variable if set
	if envPath := os.Getenv("DB_PATH"); envPath != "" && *dbPath == "data/game_catalog.sqlite" {
		*dbPath = envPath
	}

	log.Printf("[SEED] Initializing SQLite database at: %s", *dbPath)
	db, err := database.InitDB(*dbPath)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("[SEED] Seeding example games (reset=%v)...", *reset)
	count, err := database.SeedDatabase(db, *reset)
	if err != nil {
		log.Fatalf("[FATAL] Failed to seed database: %v", err)
	}

	log.Printf("[SUCCESS] Successfully seeded %d example games into %s", count, *dbPath)
}
