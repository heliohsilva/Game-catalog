package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var validIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,64}$`)

// Game represents a game catalog entry
type Game struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Platform       string   `json:"platform"`
	Subcategory    *string  `json:"subcategory,omitempty"`
	Genre          string   `json:"genre"`
	TimeToBeat     *string  `json:"timeToBeat,omitempty"`
	TimeToBeatMain *float64 `json:"timeToBeatMain,omitempty"`
	AddedAt        string   `json:"addedAt"`
	CreatedAt      string   `json:"createdAt,omitempty"`
	UpdatedAt      string   `json:"updatedAt,omitempty"`
}

// CreateGameRequest defines the payload for creating a new game
type CreateGameRequest struct {
	ID             *string  `json:"id"`
	Title          string   `json:"title"`
	Platform       string   `json:"platform"`
	Subcategory    *string  `json:"subcategory"`
	Genre          string   `json:"genre"`
	TimeToBeat     *string  `json:"timeToBeat"`
	TimeToBeatMain *float64 `json:"timeToBeatMain"`
	AddedAt        *string  `json:"addedAt"`
}

// UnmarshalJSON supports both camelCase and snake_case keys
func (r *CreateGameRequest) UnmarshalJSON(data []byte) error {
	type Alias CreateGameRequest
	aux := &struct {
		*Alias
		SnakeTimeToBeat     *string  `json:"time_to_beat"`
		SnakeTimeToBeatMain *float64 `json:"time_to_beat_main"`
		SnakeAddedAt        *string  `json:"added_at"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.TimeToBeat == nil && aux.SnakeTimeToBeat != nil {
		r.TimeToBeat = aux.SnakeTimeToBeat
	}
	if r.TimeToBeatMain == nil && aux.SnakeTimeToBeatMain != nil {
		r.TimeToBeatMain = aux.SnakeTimeToBeatMain
	}
	if r.AddedAt == nil && aux.SnakeAddedAt != nil {
		r.AddedAt = aux.SnakeAddedAt
	}

	return nil
}

// Validate validates and normalizes the create request
func (r *CreateGameRequest) Validate() (*Game, error) {
	// 1. Title validation
	title := strings.TrimSpace(r.Title)
	if title == "" {
		return nil, errors.New("title is required and cannot be empty or whitespace only")
	}
	if len(title) > 255 {
		return nil, errors.New("title cannot exceed 255 characters")
	}

	// 2. Platform validation
	normPlatform, ok := NormalizePlatform(r.Platform)
	if !ok {
		return nil, fmt.Errorf("invalid platform '%s'. Valid platforms are: %s", r.Platform, strings.Join(ValidPlatforms, ", "))
	}

	// 3. Subcategory validation
	var normSubcat *string
	if r.Subcategory != nil {
		subTrimmed := strings.TrimSpace(*r.Subcategory)
		if subTrimmed != "" {
			canonicalSub, subOk := NormalizeSubcategory(normPlatform, subTrimmed)
			if !subOk {
				allowed := ValidSubcategories[normPlatform]
				return nil, fmt.Errorf("invalid subcategory '%s' for platform '%s'. Allowed subcategories: %s",
					subTrimmed, normPlatform, strings.Join(allowed, ", "))
			}
			normSubcat = &canonicalSub
		}
	}

	// 4. Genre validation
	genre := strings.TrimSpace(r.Genre)
	if genre == "" {
		genre = "Gaming"
	} else if len(genre) > 100 {
		return nil, errors.New("genre cannot exceed 100 characters")
	}

	// 5. Time to beat validation
	var ttbMain *float64
	if r.TimeToBeatMain != nil {
		val := *r.TimeToBeatMain
		if val < 0 {
			return nil, errors.New("timeToBeatMain cannot be negative")
		}
		if val > 100000 {
			return nil, errors.New("timeToBeatMain cannot exceed 100,000 hours")
		}
		ttbMain = &val
	}

	var ttbStr *string
	if r.TimeToBeat != nil {
		s := strings.TrimSpace(*r.TimeToBeat)
		if s != "" {
			if len(s) > 50 {
				return nil, errors.New("timeToBeat string cannot exceed 50 characters")
			}
			ttbStr = &s
		}
	}
	if ttbStr == nil && ttbMain != nil {
		derived := fmt.Sprintf("%.0fh", *ttbMain)
		ttbStr = &derived
	}

	// 6. AddedAt validation
	var addedAt string
	if r.AddedAt != nil && strings.TrimSpace(*r.AddedAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*r.AddedAt))
		if err != nil {
			return nil, fmt.Errorf("addedAt must be in valid RFC3339 / ISO8601 format (e.g. 2026-01-01T12:00:00Z): %w", err)
		}
		addedAt = parsed.UTC().Format(time.RFC3339)
	} else {
		addedAt = time.Now().UTC().Format(time.RFC3339)
	}

	// 7. ID validation & generation
	var id string
	if r.ID != nil && strings.TrimSpace(*r.ID) != "" {
		trimmedID := strings.TrimSpace(*r.ID)
		if !validIDRegex.MatchString(trimmedID) {
			return nil, errors.New("custom id must be between 1 and 64 characters and contain only alphanumeric, dash, dot, or underscore")
		}
		id = trimmedID
	} else {
		id = "game-" + uuid.New().String()
	}

	now := time.Now().UTC().Format(time.RFC3339)

	return &Game{
		ID:             id,
		Title:          title,
		Platform:       normPlatform,
		Subcategory:    normSubcat,
		Genre:          genre,
		TimeToBeat:     ttbStr,
		TimeToBeatMain: ttbMain,
		AddedAt:        addedAt,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// UpdateGameRequest defines the payload for fully replacing a game
type UpdateGameRequest struct {
	Title          string   `json:"title"`
	Platform       string   `json:"platform"`
	Subcategory    *string  `json:"subcategory"`
	Genre          string   `json:"genre"`
	TimeToBeat     *string  `json:"timeToBeat"`
	TimeToBeatMain *float64 `json:"timeToBeatMain"`
	AddedAt        *string  `json:"addedAt"`
}

// UnmarshalJSON supports both camelCase and snake_case keys
func (r *UpdateGameRequest) UnmarshalJSON(data []byte) error {
	type Alias UpdateGameRequest
	aux := &struct {
		*Alias
		SnakeTimeToBeat     *string  `json:"time_to_beat"`
		SnakeTimeToBeatMain *float64 `json:"time_to_beat_main"`
		SnakeAddedAt        *string  `json:"added_at"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.TimeToBeat == nil && aux.SnakeTimeToBeat != nil {
		r.TimeToBeat = aux.SnakeTimeToBeat
	}
	if r.TimeToBeatMain == nil && aux.SnakeTimeToBeatMain != nil {
		r.TimeToBeatMain = aux.SnakeTimeToBeatMain
	}
	if r.AddedAt == nil && aux.SnakeAddedAt != nil {
		r.AddedAt = aux.SnakeAddedAt
	}

	return nil
}

// Validate validates and applies update onto an existing game
func (r *UpdateGameRequest) Validate(existing *Game) (*Game, error) {
	title := strings.TrimSpace(r.Title)
	if title == "" {
		return nil, errors.New("title is required and cannot be empty or whitespace only")
	}
	if len(title) > 255 {
		return nil, errors.New("title cannot exceed 255 characters")
	}

	normPlatform, ok := NormalizePlatform(r.Platform)
	if !ok {
		return nil, fmt.Errorf("invalid platform '%s'. Valid platforms are: %s", r.Platform, strings.Join(ValidPlatforms, ", "))
	}

	var normSubcat *string
	if r.Subcategory != nil {
		subTrimmed := strings.TrimSpace(*r.Subcategory)
		if subTrimmed != "" {
			canonicalSub, subOk := NormalizeSubcategory(normPlatform, subTrimmed)
			if !subOk {
				allowed := ValidSubcategories[normPlatform]
				return nil, fmt.Errorf("invalid subcategory '%s' for platform '%s'. Allowed subcategories: %s",
					subTrimmed, normPlatform, strings.Join(allowed, ", "))
			}
			normSubcat = &canonicalSub
		}
	}

	genre := strings.TrimSpace(r.Genre)
	if genre == "" {
		genre = "Gaming"
	} else if len(genre) > 100 {
		return nil, errors.New("genre cannot exceed 100 characters")
	}

	var ttbMain *float64
	if r.TimeToBeatMain != nil {
		val := *r.TimeToBeatMain
		if val < 0 {
			return nil, errors.New("timeToBeatMain cannot be negative")
		}
		if val > 100000 {
			return nil, errors.New("timeToBeatMain cannot exceed 100,000 hours")
		}
		ttbMain = &val
	}

	var ttbStr *string
	if r.TimeToBeat != nil {
		s := strings.TrimSpace(*r.TimeToBeat)
		if s != "" {
			if len(s) > 50 {
				return nil, errors.New("timeToBeat string cannot exceed 50 characters")
			}
			ttbStr = &s
		}
	}
	if ttbStr == nil && ttbMain != nil {
		derived := fmt.Sprintf("%.0fh", *ttbMain)
		ttbStr = &derived
	}

	addedAt := existing.AddedAt
	if r.AddedAt != nil && strings.TrimSpace(*r.AddedAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*r.AddedAt))
		if err != nil {
			return nil, fmt.Errorf("addedAt must be in valid RFC3339 / ISO8601 format: %w", err)
		}
		addedAt = parsed.UTC().Format(time.RFC3339)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	return &Game{
		ID:             existing.ID,
		Title:          title,
		Platform:       normPlatform,
		Subcategory:    normSubcat,
		Genre:          genre,
		TimeToBeat:     ttbStr,
		TimeToBeatMain: ttbMain,
		AddedAt:        addedAt,
		CreatedAt:      existing.CreatedAt,
		UpdatedAt:      now,
	}, nil
}

// PatchGameRequest defines the payload for partially updating a game
type PatchGameRequest struct {
	Title          *string  `json:"title"`
	Platform       *string  `json:"platform"`
	Subcategory    *string  `json:"subcategory"`
	Genre          *string  `json:"genre"`
	TimeToBeat     *string  `json:"timeToBeat"`
	TimeToBeatMain *float64 `json:"timeToBeatMain"`
	AddedAt        *string  `json:"addedAt"`
}

// UnmarshalJSON supports both camelCase and snake_case keys for patch
func (r *PatchGameRequest) UnmarshalJSON(data []byte) error {
	type Alias PatchGameRequest
	aux := &struct {
		*Alias
		SnakeTimeToBeat     *string  `json:"time_to_beat"`
		SnakeTimeToBeatMain *float64 `json:"time_to_beat_main"`
		SnakeAddedAt        *string  `json:"added_at"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.TimeToBeat == nil && aux.SnakeTimeToBeat != nil {
		r.TimeToBeat = aux.SnakeTimeToBeat
	}
	if r.TimeToBeatMain == nil && aux.SnakeTimeToBeatMain != nil {
		r.TimeToBeatMain = aux.SnakeTimeToBeatMain
	}
	if r.AddedAt == nil && aux.SnakeAddedAt != nil {
		r.AddedAt = aux.SnakeAddedAt
	}

	return nil
}

// Apply applies patch updates onto an existing game and validates
func (r *PatchGameRequest) Apply(existing *Game) (*Game, error) {
	updated := *existing

	// Platform
	if r.Platform != nil {
		normPlatform, ok := NormalizePlatform(*r.Platform)
		if !ok {
			return nil, fmt.Errorf("invalid platform '%s'. Valid platforms are: %s", *r.Platform, strings.Join(ValidPlatforms, ", "))
		}
		updated.Platform = normPlatform
	}

	// Title
	if r.Title != nil {
		title := strings.TrimSpace(*r.Title)
		if title == "" {
			return nil, errors.New("title cannot be empty or whitespace only")
		}
		if len(title) > 255 {
			return nil, errors.New("title cannot exceed 255 characters")
		}
		updated.Title = title
	}

	// Subcategory
	if r.Subcategory != nil {
		subTrimmed := strings.TrimSpace(*r.Subcategory)
		if subTrimmed == "" {
			updated.Subcategory = nil
		} else {
			canonicalSub, subOk := NormalizeSubcategory(updated.Platform, subTrimmed)
			if !subOk {
				allowed := ValidSubcategories[updated.Platform]
				return nil, fmt.Errorf("invalid subcategory '%s' for platform '%s'. Allowed subcategories: %s",
					subTrimmed, updated.Platform, strings.Join(allowed, ", "))
			}
			updated.Subcategory = &canonicalSub
		}
	} else if r.Platform != nil && updated.Subcategory != nil {
		// If platform changed but subcategory was not provided, verify existing subcategory is still valid
		_, subOk := NormalizeSubcategory(updated.Platform, *updated.Subcategory)
		if !subOk {
			// Subcategory does not match new platform, clear it
			updated.Subcategory = nil
		}
	}

	// Genre
	if r.Genre != nil {
		genre := strings.TrimSpace(*r.Genre)
		if genre == "" {
			genre = "Gaming"
		} else if len(genre) > 100 {
			return nil, errors.New("genre cannot exceed 100 characters")
		}
		updated.Genre = genre
	}

	// TimeToBeatMain
	if r.TimeToBeatMain != nil {
		val := *r.TimeToBeatMain
		if val < 0 {
			return nil, errors.New("timeToBeatMain cannot be negative")
		}
		if val > 100000 {
			return nil, errors.New("timeToBeatMain cannot exceed 100,000 hours")
		}
		updated.TimeToBeatMain = &val
	}

	// TimeToBeat string
	if r.TimeToBeat != nil {
		s := strings.TrimSpace(*r.TimeToBeat)
		if s == "" {
			updated.TimeToBeat = nil
		} else {
			if len(s) > 50 {
				return nil, errors.New("timeToBeat string cannot exceed 50 characters")
			}
			updated.TimeToBeat = &s
		}
	} else if r.TimeToBeatMain != nil && updated.TimeToBeat == nil {
		derived := fmt.Sprintf("%.0fh", *updated.TimeToBeatMain)
		updated.TimeToBeat = &derived
	}

	// AddedAt
	if r.AddedAt != nil && strings.TrimSpace(*r.AddedAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*r.AddedAt))
		if err != nil {
			return nil, fmt.Errorf("addedAt must be in valid RFC3339 / ISO8601 format: %w", err)
		}
		updated.AddedAt = parsed.UTC().Format(time.RFC3339)
	}

	updated.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return &updated, nil
}

// ListGamesQuery holds query parameters for filtering, sorting, and pagination
type ListGamesQuery struct {
	Platform    string
	Subcategory string
	Genre       string
	Search      string
	MinHours    *float64
	MaxHours    *float64
	SortBy      string
	Order       string
	Page        int
	PageSize    int
}

// PaginatedResponse wraps a list of games with pagination metadata
type PaginatedResponse struct {
	Data       []Game `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	TotalPages int    `json:"totalPages"`
}
