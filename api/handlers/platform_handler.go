package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"game-catalog/api/models"
	"game-catalog/api/repository"
)

// PlatformHandler handles HTTP requests for platform metadata and management
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

// CreatePlatform handles POST /api/v1/platforms
func (h *PlatformHandler) CreatePlatform(c *gin.Context) {
	var req models.CreatePlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Invalid request payload: "+err.Error())
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED",
			"Platform name is required and cannot be empty")
		return
	}

	if len(name) > 100 {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED",
			"Platform name cannot exceed 100 characters")
		return
	}

	info, err := h.repo.AddPlatform(c.Request.Context(), name, req.Subcategories)
	if err != nil {
		if errors.Is(err, repository.ErrPlatformExists) {
			RespondError(c, http.StatusConflict, "PLATFORM_ALREADY_EXISTS",
				"Platform '"+name+"' already exists")
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to create platform: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, gin.H{
		"platform": info,
		"message":  "Platform created successfully",
	})
}

// DeletePlatform handles DELETE /api/v1/platforms/:name
func (h *PlatformHandler) DeletePlatform(c *gin.Context) {
	name := c.Param("name")
	if strings.TrimSpace(name) == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_PARAM", "Platform name parameter is required")
		return
	}

	force := strings.EqualFold(c.Query("force"), "true") || c.Query("force") == "1"

	err := h.repo.DeletePlatform(c.Request.Context(), name, force)
	if err != nil {
		if errors.Is(err, repository.ErrPlatformNotFound) {
			RespondError(c, http.StatusNotFound, "PLATFORM_NOT_FOUND",
				"Platform '"+name+"' not found")
			return
		}
		if strings.Contains(err.Error(), "still assigned") {
			RespondError(c, http.StatusConflict, "PLATFORM_HAS_GAMES", err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to delete platform: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Platform '" + name + "' deleted successfully",
	})
}

// CreateSubcategory handles POST /api/v1/platforms/:name/subcategories
func (h *PlatformHandler) CreateSubcategory(c *gin.Context) {
	platformName := c.Param("name")
	if strings.TrimSpace(platformName) == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_PARAM", "Platform name parameter is required")
		return
	}

	var req models.CreateSubcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Invalid request payload: "+err.Error())
		return
	}

	subName := strings.TrimSpace(req.Name)
	if subName == "" {
		subName = strings.TrimSpace(req.Subcategory)
	}
	if subName == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED",
			"Subcategory name is required and cannot be empty")
		return
	}

	if len(subName) > 100 {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED",
			"Subcategory name cannot exceed 100 characters")
		return
	}

	err := h.repo.AddSubcategory(c.Request.Context(), platformName, subName)
	if err != nil {
		if errors.Is(err, repository.ErrPlatformNotFound) {
			RespondError(c, http.StatusNotFound, "PLATFORM_NOT_FOUND",
				"Platform '"+platformName+"' not found")
			return
		}
		if errors.Is(err, repository.ErrSubcategoryExists) {
			RespondError(c, http.StatusConflict, "SUBCATEGORY_ALREADY_EXISTS",
				"Subcategory '"+subName+"' already exists for platform '"+platformName+"'")
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to add subcategory: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, gin.H{
		"platform":    platformName,
		"subcategory": subName,
		"message":     "Subcategory added successfully",
	})
}

// DeleteSubcategory handles DELETE /api/v1/platforms/:name/subcategories/:sub
func (h *PlatformHandler) DeleteSubcategory(c *gin.Context) {
	platformName := c.Param("name")
	subName := c.Param("sub")

	if strings.TrimSpace(platformName) == "" || strings.TrimSpace(subName) == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_PARAM", "Platform and subcategory parameters are required")
		return
	}

	force := strings.EqualFold(c.Query("force"), "true") || c.Query("force") == "1"

	err := h.repo.DeleteSubcategory(c.Request.Context(), platformName, subName, force)
	if err != nil {
		if errors.Is(err, repository.ErrPlatformNotFound) {
			RespondError(c, http.StatusNotFound, "PLATFORM_NOT_FOUND",
				"Platform '"+platformName+"' not found")
			return
		}
		if errors.Is(err, repository.ErrSubcategoryNotFound) {
			RespondError(c, http.StatusNotFound, "SUBCATEGORY_NOT_FOUND",
				"Subcategory '"+subName+"' not found for platform '"+platformName+"'")
			return
		}
		if strings.Contains(err.Error(), "currently use it") {
			RespondError(c, http.StatusConflict, "SUBCATEGORY_IN_USE", err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to delete subcategory: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Subcategory '" + subName + "' removed from '" + platformName + "' successfully",
	})
}
