package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db        *sql.DB
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now().UTC(),
	}
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	ctx := c.Request.Context()

	dbStatus := "connected"
	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "error: " + err.Error()
	}

	uptime := time.Since(h.startTime).String()

	status := http.StatusOK
	if dbStatus != "connected" {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":   "healthy",
		"service":  "game-catalog-api",
		"database": dbStatus,
		"uptime":   uptime,
	})
}
