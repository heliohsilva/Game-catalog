package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"game-catalog/api/config"
	"game-catalog/api/handlers"
	"game-catalog/api/middleware"
	"game-catalog/api/repository"
)

// SetupRouter sets up Gin engine with all routes, middleware, and handlers
func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(gin.Logger())
	r.Use(middleware.JSONRecovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.MaxBodySize(cfg.MaxBodySizeBytes))

	// Repositories & Handlers
	repo := repository.NewGameRepository(db)
	gameHandler := handlers.NewGameHandler(repo)
	statsHandler := handlers.NewStatsHandler(repo)
	platformHandler := handlers.NewPlatformHandler(repo)
	seedHandler := handlers.NewSeedHandler(db)
	healthHandler := handlers.NewHealthHandler(db)
	hltbHandler := handlers.NewHLTBHandler()

	// Health Check
	r.GET("/health", healthHandler.Health)

	// API Group v1
	v1 := r.Group("/api/v1")
	{
		// Ping & Health
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		v1.GET("/health", healthHandler.Health)

		// HowLongToBeat proxy
		v1.GET("/hltb", hltbHandler.Search)

		// Games CRUD & operations
		v1.GET("/games", gameHandler.List)
		v1.GET("/games/:id", gameHandler.Get)
		v1.POST("/games", gameHandler.Create)
		v1.PUT("/games/:id", gameHandler.Update)
		v1.PATCH("/games/:id", gameHandler.Patch)
		v1.DELETE("/games/:id", gameHandler.Delete)
		v1.POST("/games/batch", gameHandler.BatchCreate)
		v1.POST("/games/bulk", gameHandler.BatchCreate)
		v1.DELETE("/games", gameHandler.Clear)

		// Analytics & Stats
		v1.GET("/stats", statsHandler.GetStats)
		v1.GET("/games/stats", statsHandler.GetStats)

		// Platform metadata & management
		v1.GET("/platforms", platformHandler.ListPlatforms)
		v1.POST("/platforms", platformHandler.CreatePlatform)
		v1.DELETE("/platforms/:name", platformHandler.DeletePlatform)
		v1.POST("/platforms/:name/subcategories", platformHandler.CreateSubcategory)
		v1.DELETE("/platforms/:name/subcategories/:sub", platformHandler.DeleteSubcategory)

		// Seed initial games
		v1.POST("/seed", seedHandler.Seed)
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		handlers.RespondError(c, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
	})

	return r
}
