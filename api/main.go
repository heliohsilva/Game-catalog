package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game-catalog/api/config"
	"game-catalog/api/database"
	"game-catalog/api/repository"
	"game-catalog/api/server"
)

func main() {
	cfg := config.Load()

	log.Printf("[SERVER] Starting Game Catalog API on port %s (Gin mode: %s)...", cfg.Port, cfg.GinMode)
	log.Printf("[DATABASE] Connecting to SQLite database at: %s", cfg.DatabasePath)

	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize SQLite database: %v", err)
	}
	defer db.Close()

	repo := repository.NewGameRepository(db)

	// Check if database needs seeding
	if cfg.AutoSeed {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		count, err := repo.Count(ctx)
		cancel()
		if err != nil {
			log.Printf("[WARN] Failed to count games in database: %v", err)
		} else if count == 0 {
			log.Printf("[DATABASE] Empty database detected. Seeding initial catalog games...")
			inserted, err := database.SeedDatabase(db, false)
			if err != nil {
				log.Printf("[ERROR] Seeding initial games failed: %v", err)
			} else {
				log.Printf("[DATABASE] Successfully seeded %d initial games!", inserted)
			}
		} else {
			log.Printf("[DATABASE] Database already contains %d games.", count)
		}
	}

	router := server.SetupRouter(cfg, db)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[SERVER] HTTP server listening on http://0.0.0.0:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server listen error: %v", err)
		}
	}()

	<-quit
	log.Println("[SERVER] Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[SERVER] Server forced to shutdown: %v", err)
	}

	log.Println("[SERVER] Server exited cleanly.")
}
