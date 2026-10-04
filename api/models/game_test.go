package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNormalizePlatform(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedOk  bool
	}{
		{"PC", PlatformPC, true},
		{"pc", PlatformPC, true},
		{"  PC  ", PlatformPC, true},
		{"PlayStation", PlatformPlayStation, true},
		{"playstation", PlatformPlayStation, true},
		{"Nintendo Switch", PlatformNintendoSwitch, true},
		{"nintendo switch", PlatformNintendoSwitch, true},
		{"Xbox", PlatformXbox, true},
		{"xbox", PlatformXbox, true},
		{"Retro / Emulation", PlatformRetroEmulation, true},
		{"retro / emulation", PlatformRetroEmulation, true},
		{"Atari", "", false},
		{"", "", false},
		{"   ", "", false},
	}

	for _, tt := range tests {
		actual, ok := NormalizePlatform(tt.input)
		if ok != tt.expectedOk {
			t.Errorf("NormalizePlatform(%q) ok = %v, expected %v", tt.input, ok, tt.expectedOk)
		}
		if actual != tt.expected {
			t.Errorf("NormalizePlatform(%q) = %q, expected %q", tt.input, actual, tt.expected)
		}
	}
}

func TestNormalizeSubcategory(t *testing.T) {
	tests := []struct {
		platform   string
		subcat     string
		expected   string
		expectedOk bool
	}{
		{PlatformPC, "Steam", "Steam", true},
		{PlatformPC, "steam", "Steam", true},
		{PlatformPC, "GOG", "GOG", true},
		{PlatformPC, "Epic", "Epic", true},
		{PlatformPC, "PS5", "", false}, // PS5 does not belong to PC
		{PlatformPlayStation, "PS5", "PS5", true},
		{PlatformPlayStation, "ps4", "PS4", true},
		{PlatformPlayStation, "Steam", "", false},
		{PlatformNintendoSwitch, "Switch", "Switch", true},
		{PlatformXbox, "Xbox Series X/S", "Xbox Series X/S", true},
		{PlatformXbox, "xbox 360", "Xbox 360", true},
		{PlatformRetroEmulation, "SNES", "SNES", true},
		{PlatformRetroEmulation, "genesis", "Genesis", true},
		{"UnknownPlatform", "Steam", "", false},
		{PlatformPC, "", "", false},
		{PlatformPC, "   ", "", false},
	}

	for _, tt := range tests {
		actual, ok := NormalizeSubcategory(tt.platform, tt.subcat)
		if ok != tt.expectedOk {
			t.Errorf("NormalizeSubcategory(%q, %q) ok = %v, expected %v", tt.platform, tt.subcat, ok, tt.expectedOk)
		}
		if actual != tt.expected {
			t.Errorf("NormalizeSubcategory(%q, %q) = %q, expected %q", tt.platform, tt.subcat, actual, tt.expected)
		}
	}
}

func TestCreateGameRequestValidation(t *testing.T) {
	hours25 := 25.0
	negHours := -5.0
	hugeHours := 200000.0
	validSub := "Steam"
	invalidSub := "PS5"
	validAddedAt := "2026-03-01T12:00:00Z"
	invalidAddedAt := "not-a-date"

	t.Run("Valid request with all fields", func(t *testing.T) {
		req := CreateGameRequest{
			Title:          "Portal 2",
			Platform:       "PC",
			Subcategory:    &validSub,
			Genre:          "Puzzle",
			TimeToBeat:     &validSub,
			TimeToBeatMain: &hours25,
			AddedAt:        &validAddedAt,
		}
		g, err := req.Validate()
		if err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
		if g.Title != "Portal 2" {
			t.Errorf("expected title Portal 2, got %s", g.Title)
		}
		if g.Platform != PlatformPC {
			t.Errorf("expected platform PC, got %s", g.Platform)
		}
		if g.Subcategory == nil || *g.Subcategory != "Steam" {
			t.Errorf("expected subcategory Steam, got %v", g.Subcategory)
		}
		if g.AddedAt != validAddedAt {
			t.Errorf("expected addedAt %s, got %s", validAddedAt, g.AddedAt)
		}
		if !strings.HasPrefix(g.ID, "game-") {
			t.Errorf("expected auto-generated ID starting with game-, got %s", g.ID)
		}
	})

	t.Run("Empty and whitespace title fails", func(t *testing.T) {
		reqEmpty := CreateGameRequest{Title: "", Platform: "PC"}
		if _, err := reqEmpty.Validate(); err == nil {
			t.Error("expected error for empty title")
		}

		reqWhitespace := CreateGameRequest{Title: "   \t\n  ", Platform: "PC"}
		if _, err := reqWhitespace.Validate(); err == nil {
			t.Error("expected error for whitespace title")
		}
	})

	t.Run("Title exceeding 255 chars fails", func(t *testing.T) {
		longTitle := strings.Repeat("A", 256)
		req := CreateGameRequest{Title: longTitle, Platform: "PC"}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for title > 255 chars")
		}
	})

	t.Run("Invalid platform fails", func(t *testing.T) {
		req := CreateGameRequest{Title: "Game", Platform: "Dreamcast"}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for invalid platform")
		}
	})

	t.Run("Incompatible subcategory fails", func(t *testing.T) {
		req := CreateGameRequest{
			Title:       "Game",
			Platform:    "PC",
			Subcategory: &invalidSub,
		}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for subcategory not matching platform")
		}
	})

	t.Run("Negative playtime fails", func(t *testing.T) {
		req := CreateGameRequest{
			Title:          "Game",
			Platform:       "PC",
			TimeToBeatMain: &negHours,
		}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for negative timeToBeatMain")
		}
	})

	t.Run("Excessive playtime (>100k hours) fails", func(t *testing.T) {
		req := CreateGameRequest{
			Title:          "Game",
			Platform:       "PC",
			TimeToBeatMain: &hugeHours,
		}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for timeToBeatMain > 100,000")
		}
	})

	t.Run("Invalid addedAt timestamp fails", func(t *testing.T) {
		req := CreateGameRequest{
			Title:    "Game",
			Platform: "PC",
			AddedAt:  &invalidAddedAt,
		}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for invalid addedAt RFC3339 string")
		}
	})

	t.Run("Custom valid ID is accepted", func(t *testing.T) {
		customID := "custom-id_123"
		req := CreateGameRequest{
			ID:       &customID,
			Title:    "Game",
			Platform: "PC",
		}
		g, err := req.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.ID != customID {
			t.Errorf("expected ID %s, got %s", customID, g.ID)
		}
	})

	t.Run("Custom invalid ID with illegal chars fails", func(t *testing.T) {
		badID := "bad/id;drop"
		req := CreateGameRequest{
			ID:       &badID,
			Title:    "Game",
			Platform: "PC",
		}
		if _, err := req.Validate(); err == nil {
			t.Error("expected error for ID with illegal characters")
		}
	})

	t.Run("Dual JSON casing support (camelCase & snake_case)", func(t *testing.T) {
		snakeJSON := `{"title":"Hades","platform":"PC","time_to_beat_main":25,"added_at":"2026-01-01T00:00:00Z"}`
		var req CreateGameRequest
		if err := json.Unmarshal([]byte(snakeJSON), &req); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		g, err := req.Validate()
		if err != nil {
			t.Fatalf("validate failed: %v", err)
		}
		if g.TimeToBeatMain == nil || *g.TimeToBeatMain != 25 {
			t.Errorf("expected timeToBeatMain 25, got %v", g.TimeToBeatMain)
		}
		if g.AddedAt != "2026-01-01T00:00:00Z" {
			t.Errorf("expected addedAt 2026-01-01T00:00:00Z, got %s", g.AddedAt)
		}
		if g.TimeToBeat == nil || *g.TimeToBeat != "25h" {
			t.Errorf("expected auto-derived timeToBeat 25h, got %v", g.TimeToBeat)
		}
	})
}

func TestPatchGameRequest(t *testing.T) {
	initialSub := "Steam"
	initialHours := 15.0
	existing := &Game{
		ID:             "game-1",
		Title:          "Original Title",
		Platform:       PlatformPC,
		Subcategory:    &initialSub,
		Genre:          "Action",
		TimeToBeatMain: &initialHours,
		AddedAt:        "2026-01-01T00:00:00Z",
		CreatedAt:      "2026-01-01T00:00:00Z",
		UpdatedAt:      "2026-01-01T00:00:00Z",
	}

	t.Run("Partial title update", func(t *testing.T) {
		newTitle := "Updated Title"
		patch := PatchGameRequest{Title: &newTitle}
		updated, err := patch.Apply(existing)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Title != "Updated Title" {
			t.Errorf("expected updated title, got %s", updated.Title)
		}
		if updated.Platform != PlatformPC {
			t.Errorf("expected preserved platform PC, got %s", updated.Platform)
		}
		if *updated.Subcategory != "Steam" {
			t.Errorf("expected preserved subcategory Steam, got %s", *updated.Subcategory)
		}
	})

	t.Run("Changing platform clears incompatible subcategory", func(t *testing.T) {
		newPlatform := PlatformNintendoSwitch // "Steam" is not valid for Nintendo Switch
		patch := PatchGameRequest{Platform: &newPlatform}
		updated, err := patch.Apply(existing)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Platform != PlatformNintendoSwitch {
			t.Errorf("expected Nintendo Switch, got %s", updated.Platform)
		}
		if updated.Subcategory != nil {
			t.Errorf("expected subcategory to be cleared, got %v", *updated.Subcategory)
		}
	})

	t.Run("Empty title in patch fails", func(t *testing.T) {
		emptyTitle := "   "
		patch := PatchGameRequest{Title: &emptyTitle}
		if _, err := patch.Apply(existing); err == nil {
			t.Error("expected error when patching empty title")
		}
	})

	t.Run("Negative hours in patch fails", func(t *testing.T) {
		neg := -10.0
		patch := PatchGameRequest{TimeToBeatMain: &neg}
		if _, err := patch.Apply(existing); err == nil {
			t.Error("expected error when patching negative hours")
		}
	})
}
