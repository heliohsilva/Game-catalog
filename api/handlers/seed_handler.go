package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"game-catalog/api/database"
)

// SeedHandler handles requests to seed or reset initial catalog data
type SeedHandler struct {
	db *sql.DB
}

// NewSeedHandler creates a new SeedHandler
func NewSeedHandler(db *sql.DB) *SeedHandler {
	return &SeedHandler{db: db}
}

// Seed handles POST /api/v1/seed
func (h *SeedHandler) Seed(c *gin.Context) {
	resetStr := strings.ToLower(c.Query("reset"))
	reset := resetStr == "true" || resetStr == "1" || resetStr == "yes"

	count, err := database.SeedDatabase(h.db, reset)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "SEED_FAILED",
			"Failed to seed database: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message":  "Database seeded successfully",
		"inserted": count,
		"reset":    reset,
	})
}
