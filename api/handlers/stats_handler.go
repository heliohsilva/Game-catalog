package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"game-catalog/api/repository"
)

// StatsHandler handles HTTP requests for catalog statistics
type StatsHandler struct {
	repo repository.GameRepository
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(repo repository.GameRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// GetStats handles GET /api/v1/stats
func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.repo.GetStats(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to compute catalog statistics: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, stats)
}
