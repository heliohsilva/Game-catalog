package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"game-catalog/api/models"
	"game-catalog/api/repository"
)

// GameHandler handles HTTP requests for game operations
type GameHandler struct {
	repo repository.GameRepository
}

// NewGameHandler creates a new GameHandler
func NewGameHandler(repo repository.GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

// List handles GET /api/v1/games
func (h *GameHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var q models.ListGamesQuery
	q.Platform = c.Query("platform")
	q.Subcategory = c.Query("subcategory")
	q.Genre = c.Query("genre")

	// Support both "search" and "q"
	search := c.Query("search")
	if search == "" {
		search = c.Query("q")
	}
	q.Search = search

	// Min hours
	if minStr := c.Query("min_hours"); minStr != "" {
		val, err := strconv.ParseFloat(minStr, 64)
		if err != nil || val < 0 {
			RespondError(c, http.StatusBadRequest, "INVALID_QUERY_PARAM",
				"min_hours must be a non-negative number")
			return
		}
		q.MinHours = &val
	}

	// Max hours
	if maxStr := c.Query("max_hours"); maxStr != "" {
		val, err := strconv.ParseFloat(maxStr, 64)
		if err != nil || val < 0 {
			RespondError(c, http.StatusBadRequest, "INVALID_QUERY_PARAM",
				"max_hours must be a non-negative number")
			return
		}
		q.MaxHours = &val
	}

	// Sort & Order
	q.SortBy = c.DefaultQuery("sort_by", "title")
	q.Order = c.DefaultQuery("order", "")

	// Pagination
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		RespondError(c, http.StatusBadRequest, "INVALID_QUERY_PARAM",
			"page must be a positive integer greater than 0")
		return
	}
	q.Page = page

	pageSizeStr := c.DefaultQuery("page_size", c.DefaultQuery("limit", "50"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		RespondError(c, http.StatusBadRequest, "INVALID_QUERY_PARAM",
			"page_size (or limit) must be a positive integer greater than 0")
		return
	}
	if pageSize > 100 {
		pageSize = 100
	}
	q.PageSize = pageSize

	games, total, err := h.repo.List(ctx, q)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to retrieve games from database")
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	RespondSuccess(c, http.StatusOK, models.PaginatedResponse{
		Data:       games,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
	})
}

// Get handles GET /api/v1/games/:id
func (h *GameHandler) Get(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_ID", "Game ID cannot be empty")
		return
	}

	game, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to retrieve game")
		return
	}

	RespondSuccess(c, http.StatusOK, game)
}

// Create handles POST /api/v1/games
func (h *GameHandler) Create(c *gin.Context) {
	if !ValidateJSONContentType(c) {
		return
	}

	var req models.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			RespondError(c, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
				"Request body exceeds maximum allowed size")
			return
		}
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Failed to parse JSON body: "+err.Error())
		return
	}

	game, err := req.Validate()
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := h.repo.Create(c.Request.Context(), game); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			subStr := ""
			if game.Subcategory != nil {
				subStr = fmt.Sprintf(" (%s)", *game.Subcategory)
			}
			RespondError(c, http.StatusConflict, "DUPLICATE_GAME",
				fmt.Sprintf("Game '%s' on platform '%s'%s already exists in the catalog",
					game.Title, game.Platform, subStr))
			return
		}
		if errors.Is(err, repository.ErrDuplicateID) {
			RespondError(c, http.StatusConflict, "DUPLICATE_ID",
				fmt.Sprintf("Game with ID '%s' already exists", game.ID))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to create game: "+err.Error())
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/games/%s", game.ID))
	RespondSuccess(c, http.StatusCreated, game)
}

// Update handles PUT /api/v1/games/:id (full replace)
func (h *GameHandler) Update(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_ID", "Game ID cannot be empty")
		return
	}

	if !ValidateJSONContentType(c) {
		return
	}

	ctx := c.Request.Context()
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to retrieve game")
		return
	}

	var req models.UpdateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			RespondError(c, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
				"Request body exceeds maximum allowed size")
			return
		}
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Failed to parse JSON body: "+err.Error())
		return
	}

	updated, err := req.Validate(existing)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := h.repo.Update(ctx, updated); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			subStr := ""
			if updated.Subcategory != nil {
				subStr = fmt.Sprintf(" (%s)", *updated.Subcategory)
			}
			RespondError(c, http.StatusConflict, "DUPLICATE_GAME",
				fmt.Sprintf("Game '%s' on platform '%s'%s already exists in the catalog",
					updated.Title, updated.Platform, subStr))
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to update game: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, updated)
}

// Patch handles PATCH /api/v1/games/:id (partial update)
func (h *GameHandler) Patch(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_ID", "Game ID cannot be empty")
		return
	}

	if !ValidateJSONContentType(c) {
		return
	}

	ctx := c.Request.Context()
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to retrieve game")
		return
	}

	var req models.PatchGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			RespondError(c, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
				"Request body exceeds maximum allowed size")
			return
		}
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Failed to parse JSON body: "+err.Error())
		return
	}

	updated, err := req.Apply(existing)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := h.repo.Update(ctx, updated); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			subStr := ""
			if updated.Subcategory != nil {
				subStr = fmt.Sprintf(" (%s)", *updated.Subcategory)
			}
			RespondError(c, http.StatusConflict, "DUPLICATE_GAME",
				fmt.Sprintf("Game '%s' on platform '%s'%s already exists in the catalog",
					updated.Title, updated.Platform, subStr))
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to update game: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, updated)
}

// Delete handles DELETE /api/v1/games/:id
func (h *GameHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		RespondError(c, http.StatusBadRequest, "INVALID_ID", "Game ID cannot be empty")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			RespondError(c, http.StatusNotFound, "GAME_NOT_FOUND",
				fmt.Sprintf("Game with ID '%s' was not found", id))
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to delete game: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Game deleted successfully",
		"id":      id,
	})
}

// BatchCreate handles POST /api/v1/games/batch
func (h *GameHandler) BatchCreate(c *gin.Context) {
	if !ValidateJSONContentType(c) {
		return
	}

	var reqs []models.CreateGameRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			RespondError(c, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
				"Request body exceeds maximum allowed size")
			return
		}
		RespondError(c, http.StatusBadRequest, "INVALID_JSON",
			"Failed to parse JSON body array: "+err.Error())
		return
	}

	if len(reqs) == 0 {
		RespondError(c, http.StatusBadRequest, "EMPTY_BATCH", "Batch array cannot be empty")
		return
	}

	if len(reqs) > 500 {
		RespondError(c, http.StatusBadRequest, "BATCH_TOO_LARGE",
			"Batch size exceeds maximum allowed of 500 items")
		return
	}

	games := make([]*models.Game, 0, len(reqs))
	for idx, r := range reqs {
		g, err := r.Validate()
		if err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED",
				fmt.Sprintf("Item at index %d is invalid: %s", idx, err.Error()))
			return
		}
		games = append(games, g)
	}

	inserted, errs := h.repo.BatchCreate(c.Request.Context(), games)
	errMsgs := make([]string, 0, len(errs))
	for _, e := range errs {
		errMsgs = append(errMsgs, e.Error())
	}

	status := http.StatusCreated
	if inserted == 0 && len(errs) > 0 {
		status = http.StatusConflict
	} else if len(errs) > 0 {
		status = http.StatusMultiStatus
	}

	RespondSuccess(c, status, gin.H{
		"inserted": inserted,
		"total":    len(games),
		"errors":   errMsgs,
	})
}

// Clear handles DELETE /api/v1/games (bulk clear with ?confirm=true)
func (h *GameHandler) Clear(c *gin.Context) {
	confirm := strings.ToLower(c.Query("confirm"))
	if confirm != "true" && confirm != "1" && confirm != "yes" {
		RespondError(c, http.StatusBadRequest, "CONFIRMATION_REQUIRED",
			"Clearing the catalog requires confirmation query param ?confirm=true")
		return
	}

	deleted, err := h.repo.DeleteAll(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Failed to clear catalog: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "All games cleared successfully",
		"deleted": deleted,
	})
}
