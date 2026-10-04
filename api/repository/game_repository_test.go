package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"game-catalog/api/database"
	"game-catalog/api/models"
)

func setupTestDB(t *testing.T) (*sql.DB, GameRepository) {
	t.Helper()
	// Use an in-memory SQLite database for fast and isolated test execution
	db, err := database.InitDB(fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("failed to init test sqlite db: %v", err)
	}

	repo := NewGameRepository(db)
	return db, repo
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestGameRepository_CRUD(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// 1. Create a game
	g1 := &models.Game{
		ID:             "test-1",
		Title:          "Celeste",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Steam"),
		Genre:          "Platformer",
		TimeToBeat:     strPtr("8h"),
		TimeToBeatMain: floatPtr(8),
		AddedAt:        "2026-01-01T00:00:00Z",
		CreatedAt:      "2026-01-01T00:00:00Z",
		UpdatedAt:      "2026-01-01T00:00:00Z",
	}

	if err := repo.Create(ctx, g1); err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// 2. Retrieve game by ID
	retrieved, err := repo.GetByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("failed to get game by id: %v", err)
	}
	if retrieved.Title != "Celeste" {
		t.Errorf("expected title 'Celeste', got '%s'", retrieved.Title)
	}
	if retrieved.Subcategory == nil || *retrieved.Subcategory != "Steam" {
		t.Errorf("expected subcategory 'Steam', got '%v'", retrieved.Subcategory)
	}

	// 3. Duplicate ID prevention
	gDuplicateID := &models.Game{
		ID:          "test-1",
		Title:       "Other Game",
		Platform:    models.PlatformPlayStation,
		Genre:       "RPG",
		AddedAt:     "2026-01-01T00:00:00Z",
		CreatedAt:   "2026-01-01T00:00:00Z",
		UpdatedAt:   "2026-01-01T00:00:00Z",
	}
	if err := repo.Create(ctx, gDuplicateID); !errors.Is(err, ErrDuplicateID) {
		t.Errorf("expected ErrDuplicateID, got %v", err)
	}

	// 4. Duplicate (Title + Platform + Subcategory) prevention (case-insensitive)
	gDuplicateEntry := &models.Game{
		ID:          "test-2",
		Title:       "  celeste  ", // case and space variations
		Platform:    models.PlatformPC,
		Subcategory: strPtr("Steam"),
		Genre:       "Platformer",
		AddedAt:     "2026-01-01T00:00:00Z",
		CreatedAt:   "2026-01-01T00:00:00Z",
		UpdatedAt:   "2026-01-01T00:00:00Z",
	}
	if err := repo.Create(ctx, gDuplicateEntry); !errors.Is(err, ErrDuplicate) {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}

	// 5. Same title on DIFFERENT platform is allowed
	gSwitch := &models.Game{
		ID:          "test-switch-1",
		Title:       "Celeste",
		Platform:    models.PlatformNintendoSwitch,
		Subcategory: strPtr("Switch"),
		Genre:       "Platformer",
		AddedAt:     "2026-01-01T00:00:00Z",
		CreatedAt:   "2026-01-01T00:00:00Z",
		UpdatedAt:   "2026-01-01T00:00:00Z",
	}
	if err := repo.Create(ctx, gSwitch); err != nil {
		t.Errorf("expected Celeste on Switch to be allowed, got error: %v", err)
	}

	// 6. Update game
	retrieved.Title = "Celeste: Farewell Edition"
	retrieved.TimeToBeatMain = floatPtr(12)
	retrieved.TimeToBeat = strPtr("12h")
	if err := repo.Update(ctx, retrieved); err != nil {
		t.Fatalf("failed to update game: %v", err)
	}

	updated, err := repo.GetByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("failed to get updated game: %v", err)
	}
	if updated.Title != "Celeste: Farewell Edition" || *updated.TimeToBeatMain != 12 {
		t.Errorf("updated values mismatch: title=%s, hours=%v", updated.Title, *updated.TimeToBeatMain)
	}

	// 7. Update to collision with another game
	updated.Title = "Celeste"
	updated.Platform = models.PlatformNintendoSwitch
	updated.Subcategory = strPtr("Switch")
	if err := repo.Update(ctx, updated); !errors.Is(err, ErrDuplicate) {
		t.Errorf("expected ErrDuplicate when updating to existing title/platform/subcat, got %v", err)
	}

	// 8. Delete game
	if err := repo.Delete(ctx, "test-1"); err != nil {
		t.Fatalf("failed to delete game: %v", err)
	}

	// Verify deleted
	_, err = repo.GetByID(ctx, "test-1")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}

	// Delete non-existent game
	if err := repo.Delete(ctx, "non-existent-id"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound when deleting non-existent game, got %v", err)
	}
}

func TestGameRepository_ListFilterSearchSort(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// Seed sample games
	games := []*models.Game{
		{
			ID:             "g1",
			Title:          "The Witcher 3: Wild Hunt",
			Platform:       models.PlatformPC,
			Subcategory:    strPtr("GOG"),
			Genre:          "Action RPG",
			TimeToBeatMain: floatPtr(52),
			AddedAt:        "2026-01-01T10:00:00Z",
			CreatedAt:      "2026-01-01T10:00:00Z",
			UpdatedAt:      "2026-01-01T10:00:00Z",
		},
		{
			ID:             "g2",
			Title:          "Hollow Knight",
			Platform:       models.PlatformPC,
			Subcategory:    strPtr("Steam"),
			Genre:          "Metroidvania",
			TimeToBeatMain: floatPtr(30),
			AddedAt:        "2026-01-02T10:00:00Z",
			CreatedAt:      "2026-01-02T10:00:00Z",
			UpdatedAt:      "2026-01-02T10:00:00Z",
		},
		{
			ID:             "g3",
			Title:          "Bloodborne",
			Platform:       models.PlatformPlayStation,
			Subcategory:    strPtr("PS4"),
			Genre:          "Action RPG",
			TimeToBeatMain: floatPtr(34),
			AddedAt:        "2026-01-03T10:00:00Z",
			CreatedAt:      "2026-01-03T10:00:00Z",
			UpdatedAt:      "2026-01-03T10:00:00Z",
		},
		{
			ID:             "g4",
			Title:          "Super Mario 64",
			Platform:       models.PlatformRetroEmulation,
			Subcategory:    strPtr("N64"),
			Genre:          "Platformer",
			TimeToBeatMain: floatPtr(12),
			AddedAt:        "2026-01-04T10:00:00Z",
			CreatedAt:      "2026-01-04T10:00:00Z",
			UpdatedAt:      "2026-01-04T10:00:00Z",
		},
		{
			ID:             "g5",
			Title:          "Game Without Playtime",
			Platform:       models.PlatformPC,
			Subcategory:    strPtr("Epic"),
			Genre:          "Indie",
			TimeToBeatMain: nil,
			AddedAt:        "2026-01-05T10:00:00Z",
			CreatedAt:      "2026-01-05T10:00:00Z",
			UpdatedAt:      "2026-01-05T10:00:00Z",
		},
	}

	for _, g := range games {
		if err := repo.Create(ctx, g); err != nil {
			t.Fatalf("failed inserting test game %s: %v", g.Title, err)
		}
	}

	t.Run("List all", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 5 || len(res) != 5 {
			t.Errorf("expected 5 games, got total=%d, len=%d", total, len(res))
		}
	})

	t.Run("Filter by platform", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Platform: "PC"})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 3 || len(res) != 3 {
			t.Errorf("expected 3 PC games, got total=%d, len=%d", total, len(res))
		}
	})

	t.Run("Filter by unknown platform returns empty list", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Platform: "VirtualBoy"})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 0 || len(res) != 0 {
			t.Errorf("expected 0 games, got total=%d, len=%d", total, len(res))
		}
	})

	t.Run("Filter by subcategory", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Subcategory: "PS4"})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 1 || res[0].Title != "Bloodborne" {
			t.Errorf("expected 1 game (Bloodborne), got total=%d", total)
		}
	})

	t.Run("Filter by genre", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Genre: "Action RPG"})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 2 || len(res) != 2 {
			t.Errorf("expected 2 Action RPG games, got %d (len %d)", total, len(res))
		}
	})

	t.Run("Search query by title keyword", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Search: "Witcher"})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 1 || res[0].Title != "The Witcher 3: Wild Hunt" {
			t.Errorf("expected The Witcher 3, got total=%d", total)
		}
	})

	t.Run("Search with SQL injection chars does not crash or inject", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Search: "' OR '1'='1"})
		if err != nil {
			t.Fatalf("query error with sql injection string: %v", err)
		}
		if total != 0 || len(res) != 0 {
			t.Errorf("expected 0 results for injection payload, got %d", total)
		}
	})

	t.Run("Search with wildcards", func(t *testing.T) {
		res, total, err := repo.List(ctx, models.ListGamesQuery{Search: "%"})
		if err != nil {
			t.Fatalf("query error with wildcard: %v", err)
		}
		// Literal % is escaped, so 0 matches unless title contains literal %
		if total != 0 || len(res) != 0 {
			t.Errorf("expected 0 matches for literal %%, got %d", total)
		}
	})

	t.Run("Filter by playtime bounds", func(t *testing.T) {
		min := 20.0
		max := 40.0
		res, total, err := repo.List(ctx, models.ListGamesQuery{
			MinHours: &min,
			MaxHours: &max,
		})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		// Hollow Knight (30) and Bloodborne (34)
		if total != 2 {
			t.Errorf("expected 2 games between 20h and 40h, got %d", total)
		}
		for _, g := range res {
			if *g.TimeToBeatMain < 20 || *g.TimeToBeatMain > 40 {
				t.Errorf("game %s out of hours range: %v", g.Title, *g.TimeToBeatMain)
			}
		}
	})

	t.Run("Sort by hours DESC (NULLs last)", func(t *testing.T) {
		res, _, err := repo.List(ctx, models.ListGamesQuery{
			SortBy: "hours",
			Order:  "desc",
		})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		// First should be The Witcher 3 (52), last should be Game Without Playtime (nil)
		if res[0].Title != "The Witcher 3: Wild Hunt" {
			t.Errorf("expected first to be The Witcher 3, got %s", res[0].Title)
		}
		if res[len(res)-1].Title != "Game Without Playtime" {
			t.Errorf("expected last to be Game Without Playtime, got %s", res[len(res)-1].Title)
		}
	})

	t.Run("Pagination", func(t *testing.T) {
		page1, total, err := repo.List(ctx, models.ListGamesQuery{
			Page:     1,
			PageSize: 2,
			SortBy:   "title",
			Order:    "asc",
		})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if total != 5 || len(page1) != 2 {
			t.Errorf("expected total=5, len=2, got total=%d, len=%d", total, len(page1))
		}

		page2, _, err := repo.List(ctx, models.ListGamesQuery{
			Page:     2,
			PageSize: 2,
			SortBy:   "title",
			Order:    "asc",
		})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if len(page2) != 2 {
			t.Errorf("expected len=2 on page 2, got %d", len(page2))
		}
		if page1[0].ID == page2[0].ID {
			t.Errorf("page 1 and page 2 returned same game: %s", page1[0].ID)
		}
	})
}

func TestGameRepository_BatchAndStats(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// 1. Stats on empty DB
	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("failed to get stats on empty DB: %v", err)
	}
	if stats.TotalGames != 0 || stats.TotalPlaytimeHours != 0 || stats.AveragePlaytimeHours != 0 {
		t.Errorf("expected zeroes for empty DB stats, got %+v", stats)
	}

	// 2. BatchCreate
	batch := []*models.Game{
		{
			ID:             "b1",
			Title:          "Hades II",
			Platform:       models.PlatformPC,
			Subcategory:    strPtr("Epic"),
			Genre:          "Roguelike",
			TimeToBeatMain: floatPtr(22),
			AddedAt:        "2026-01-01T00:00:00Z",
			CreatedAt:      "2026-01-01T00:00:00Z",
			UpdatedAt:      "2026-01-01T00:00:00Z",
		},
		{
			ID:             "b2",
			Title:          "Sea of Stars",
			Platform:       models.PlatformNintendoSwitch,
			Subcategory:    strPtr("Switch"),
			Genre:          "RPG",
			TimeToBeatMain: floatPtr(28),
			AddedAt:        "2026-01-02T00:00:00Z",
			CreatedAt:      "2026-01-02T00:00:00Z",
			UpdatedAt:      "2026-01-02T00:00:00Z",
		},
		{
			ID:             "b3",
			Title:          "Halo 3",
			Platform:       models.PlatformXbox,
			Subcategory:    strPtr("Xbox 360"),
			Genre:          "FPS",
			TimeToBeatMain: floatPtr(12),
			AddedAt:        "2026-01-03T00:00:00Z",
			CreatedAt:      "2026-01-03T00:00:00Z",
			UpdatedAt:      "2026-01-03T00:00:00Z",
		},
	}

	inserted, errs := repo.BatchCreate(ctx, batch)
	if len(errs) > 0 {
		t.Fatalf("unexpected batch errors: %v", errs)
	}
	if inserted != 3 {
		t.Errorf("expected 3 inserted, got %d", inserted)
	}

	// 3. Batch with duplicate item
	dupBatch := []*models.Game{
		{
			ID:          "b4",
			Title:       "New Unique Game",
			Platform:    models.PlatformPC,
			Genre:       "Puzzle",
			AddedAt:     "2026-01-01T00:00:00Z",
			CreatedAt:   "2026-01-01T00:00:00Z",
			UpdatedAt:   "2026-01-01T00:00:00Z",
		},
		{
			ID:          "b1", // duplicate ID
			Title:       "Hades II",
			Platform:    models.PlatformPC,
			Subcategory: strPtr("Epic"),
			Genre:       "Roguelike",
			AddedAt:     "2026-01-01T00:00:00Z",
			CreatedAt:   "2026-01-01T00:00:00Z",
			UpdatedAt:   "2026-01-01T00:00:00Z",
		},
	}
	inserted2, errs2 := repo.BatchCreate(ctx, dupBatch)
	if inserted2 != 1 {
		t.Errorf("expected 1 inserted from dup batch, got %d", inserted2)
	}
	if len(errs2) != 1 {
		t.Errorf("expected 1 error from dup batch, got %d", len(errs2))
	}

	// 4. Stats with data
	statsWithData, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("stats error: %v", err)
	}
	if statsWithData.TotalGames != 4 {
		t.Errorf("expected 4 total games, got %d", statsWithData.TotalGames)
	}
	// Playtimes: 22 + 28 + 12 = 62. Shortest: Halo 3 (12). Longest: Sea of Stars (28).
	if statsWithData.TotalPlaytimeHours != 62 {
		t.Errorf("expected 62 total hours, got %f", statsWithData.TotalPlaytimeHours)
	}
	if statsWithData.ShortestGame == nil || statsWithData.ShortestGame.Title != "Halo 3" {
		t.Errorf("expected shortest game Halo 3, got %+v", statsWithData.ShortestGame)
	}
	if statsWithData.LongestGame == nil || statsWithData.LongestGame.Title != "Sea of Stars" {
		t.Errorf("expected longest game Sea of Stars, got %+v", statsWithData.LongestGame)
	}

	// 5. Platform metadata stats
	platforms, err := repo.GetPlatformStats(ctx)
	if err != nil {
		t.Fatalf("platform stats error: %v", err)
	}
	if len(platforms) != len(models.ValidPlatforms) {
		t.Errorf("expected %d platforms, got %d", len(models.ValidPlatforms), len(platforms))
	}

	// 6. DeleteAll
	deleted, err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatalf("deleteAll error: %v", err)
	}
	if deleted != 4 {
		t.Errorf("expected 4 deleted, got %d", deleted)
	}
	count, _ := repo.Count(ctx)
	if count != 0 {
		t.Errorf("expected 0 count after DeleteAll, got %d", count)
	}
}
