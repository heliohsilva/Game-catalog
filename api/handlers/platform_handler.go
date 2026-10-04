package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"game-catalog/api/repository"
)

// PlatformHandler handles HTTP requests for platform metadata
type PlatformHandler struct {
	repo repository.GameRepository
}

// NewPlatformHandler creates a new PlatformHandler
func NewPlatformHandler(repo repository.GameRepository) *PlatformHandler {
	return &PlatformHandler{repo: repo}
}

// ListPlatforms handles GET /api/v1/platforms
func (h *PlatformHandler) ListPlatforms(c *gin.Context) {
	platforms, err := h.repo.GetPlatformStats(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to retrieve platform data: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"platforms": platforms,
	})
}
