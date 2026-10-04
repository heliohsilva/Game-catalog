package config

import (
	"os"
	"testing"
)

func TestConfigLoadDefaults(t *testing.T) {
	// Clear relevant env vars
	envKeys := []string{"API_PORT", "PORT", "GIN_MODE", "DB_PATH", "SQLITE_PATH", "DATABASE_URL", "CORS_ALLOWED_ORIGINS", "AUTO_SEED", "MAX_BODY_SIZE_BYTES"}
	for _, k := range envKeys {
		os.Unsetenv(k)
	}

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.GinMode != "release" {
		t.Errorf("expected default ginMode release, got %s", cfg.GinMode)
	}
	if cfg.DatabasePath != "game_catalog.sqlite" {
		t.Errorf("expected default dbPath game_catalog.sqlite, got %s", cfg.DatabasePath)
	}
	if cfg.AutoSeed {
		t.Error("expected default autoSeed false")
	}
	if len(cfg.CORSAllowedOrigins) != 1 || cfg.CORSAllowedOrigins[0] != "*" {
		t.Errorf("expected CORS allowed origins [*], got %v", cfg.CORSAllowedOrigins)
	}
	if cfg.MaxBodySizeBytes != 2*1024*1024 {
		t.Errorf("expected default maxBodySizeBytes 2MB, got %d", cfg.MaxBodySizeBytes)
	}
}

func TestConfigLoadCustomEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("API_PORT", "9090")
	os.Setenv("GIN_MODE", "debug")
	os.Setenv("DB_PATH", "custom_dir/test.db")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000, https://example.com")
	os.Setenv("AUTO_SEED", "true")
	os.Setenv("MAX_BODY_SIZE_BYTES", "1048576")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("API_PORT")
		os.Unsetenv("GIN_MODE")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		os.Unsetenv("AUTO_SEED")
		os.Unsetenv("MAX_BODY_SIZE_BYTES")
	}()

	cfg := Load()
	if cfg.Port != "9090" {
		t.Errorf("expected API_PORT 9090, got %s", cfg.Port)
	}
	if cfg.GinMode != "debug" {
		t.Errorf("expected GIN_MODE debug, got %s", cfg.GinMode)
	}
	if cfg.DatabasePath != "custom_dir/test.db" {
		t.Errorf("expected DB_PATH custom_dir/test.db, got %s", cfg.DatabasePath)
	}
	if !cfg.AutoSeed {
		t.Error("expected autoSeed true")
	}
	if len(cfg.CORSAllowedOrigins) != 2 || cfg.CORSAllowedOrigins[0] != "http://localhost:3000" || cfg.CORSAllowedOrigins[1] != "https://example.com" {
		t.Errorf("unexpected CORS allowed origins: %v", cfg.CORSAllowedOrigins)
	}
	if cfg.MaxBodySizeBytes != 1048576 {
		t.Errorf("expected max body size 1048576, got %d", cfg.MaxBodySizeBytes)
	}
}

func TestConfigDatabaseURLFallback(t *testing.T) {
	os.Unsetenv("DB_PATH")
	os.Unsetenv("SQLITE_PATH")

	t.Run("sqlite:// prefix", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "sqlite:///tmp/db.sqlite")
		defer os.Unsetenv("DATABASE_URL")
		cfg := Load()
		if cfg.DatabasePath != "/tmp/db.sqlite" {
			t.Errorf("expected /tmp/db.sqlite, got %s", cfg.DatabasePath)
		}
	})

	t.Run("file: prefix", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "file:memdb?mode=memory")
		defer os.Unsetenv("DATABASE_URL")
		cfg := Load()
		if cfg.DatabasePath != "file:memdb?mode=memory" {
			t.Errorf("expected file:memdb?mode=memory, got %s", cfg.DatabasePath)
		}
	})

	t.Run(".db suffix", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "my_games.db")
		defer os.Unsetenv("DATABASE_URL")
		cfg := Load()
		if cfg.DatabasePath != "my_games.db" {
			t.Errorf("expected my_games.db, got %s", cfg.DatabasePath)
		}
	})
}
