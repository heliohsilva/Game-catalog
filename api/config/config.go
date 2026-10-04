package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Port               string
	GinMode            string
	DatabasePath       string
	CORSAllowedOrigins []string
	AutoSeed           bool
	MaxBodySizeBytes   int64
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	port := getEnv("API_PORT", "")
	if port == "" {
		port = getEnv("PORT", "8080")
	}

	ginMode := getEnv("GIN_MODE", "release")

	// Determine SQLite DB Path
	// Check DB_PATH, SQLITE_PATH, and parse DATABASE_URL if provided
	dbPath := getEnv("DB_PATH", "")
	if dbPath == "" {
		dbPath = getEnv("SQLITE_PATH", "")
	}
	if dbPath == "" {
		dbURL := getEnv("DATABASE_URL", "")
		if dbURL != "" {
			if strings.HasPrefix(dbURL, "sqlite://") {
				dbPath = strings.TrimPrefix(dbURL, "sqlite://")
			} else if strings.HasPrefix(dbURL, "file:") {
				dbPath = dbURL
			} else if strings.HasSuffix(dbURL, ".db") || strings.HasSuffix(dbURL, ".sqlite") {
				dbPath = dbURL
			}
		}
	}
	if dbPath == "" {
		dbPath = "game_catalog.db"
	}

	// CORS Origins
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "*")
	var origins []string
	for _, o := range strings.Split(originsStr, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	// AutoSeed: default true if empty database
	autoSeedStr := strings.ToLower(getEnv("AUTO_SEED", "true"))
	autoSeed := autoSeedStr == "true" || autoSeedStr == "1" || autoSeedStr == "yes"

	// Max Body Size in bytes (default: 2MB)
	maxBodyStr := getEnv("MAX_BODY_SIZE_BYTES", "2097152")
	maxBody, err := strconv.ParseInt(maxBodyStr, 10, 64)
	if err != nil || maxBody <= 0 {
		maxBody = 2 * 1024 * 1024
	}

	return &Config{
		Port:               port,
		GinMode:            ginMode,
		DatabasePath:       dbPath,
		CORSAllowedOrigins: origins,
		AutoSeed:           autoSeed,
		MaxBodySizeBytes:   maxBody,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		trimmed := strings.TrimSpace(val)
		if trimmed != "" {
			return trimmed
		}
	}
	return defaultVal
}
