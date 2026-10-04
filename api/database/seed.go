package database

import (
	"database/sql"
	"fmt"
	"time"

	"game-catalog/api/models"
)

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

// InitialSeedGames contains default games matching the frontend initial dataset
var InitialSeedGames = []models.Game{
	// PC
	{
		ID:             "pc-1",
		Title:          "Hollow Knight: Silksong",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Steam"),
		Genre:          "Metroidvania",
		TimeToBeat:     strPtr("35h"),
		TimeToBeatMain: floatPtr(35),
		AddedAt:        "2026-01-10T12:00:00Z",
	},
	{
		ID:             "pc-2",
		Title:          "Cyberpunk 2077: Phantom Liberty",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("GOG"),
		Genre:          "Action RPG",
		TimeToBeat:     strPtr("19h"),
		TimeToBeatMain: floatPtr(19),
		AddedAt:        "2026-02-14T15:30:00Z",
	},
	{
		ID:             "pc-3",
		Title:          "Baldur's Gate 3",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Steam"),
		Genre:          "CRPG",
		TimeToBeat:     strPtr("75h"),
		TimeToBeatMain: floatPtr(75),
		AddedAt:        "2026-01-05T09:00:00Z",
	},
	{
		ID:             "pc-4",
		Title:          "Hades II",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Epic"),
		Genre:          "Roguelike",
		TimeToBeat:     strPtr("22h"),
		TimeToBeatMain: floatPtr(22),
		AddedAt:        "2026-03-01T18:00:00Z",
	},
	{
		ID:             "pc-5",
		Title:          "Factorio: Space Age",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Steam"),
		Genre:          "Automation",
		TimeToBeat:     strPtr("80h"),
		TimeToBeatMain: floatPtr(80),
		AddedAt:        "2026-03-12T10:00:00Z",
	},
	{
		ID:             "pc-6",
		Title:          "Alan Wake 2",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("Epic"),
		Genre:          "Survival Horror",
		TimeToBeat:     strPtr("18h"),
		TimeToBeatMain: floatPtr(18),
		AddedAt:        "2026-02-20T10:00:00Z",
	},
	{
		ID:             "pc-7",
		Title:          "The Witcher 3: Wild Hunt",
		Platform:       models.PlatformPC,
		Subcategory:    strPtr("GOG"),
		Genre:          "Action RPG",
		TimeToBeat:     strPtr("52h"),
		TimeToBeatMain: floatPtr(52),
		AddedAt:        "2026-01-02T10:00:00Z",
	},

	// PlayStation
	{
		ID:             "ps-1",
		Title:          "Final Fantasy VII Rebirth",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS5"),
		Genre:          "Action RPG",
		TimeToBeat:     strPtr("48h"),
		TimeToBeatMain: floatPtr(48),
		AddedAt:        "2026-02-28T20:00:00Z",
	},
	{
		ID:             "ps-2",
		Title:          "Demon's Souls",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS5"),
		Genre:          "Souls-like",
		TimeToBeat:     strPtr("30h"),
		TimeToBeatMain: floatPtr(30),
		AddedAt:        "2026-01-18T10:10:00Z",
	},
	{
		ID:             "ps-3",
		Title:          "Bloodborne",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS4"),
		Genre:          "Action RPG",
		TimeToBeat:     strPtr("34h"),
		TimeToBeatMain: floatPtr(34),
		AddedAt:        "2026-01-11T12:00:00Z",
	},
	{
		ID:             "ps-4",
		Title:          "Metal Gear Solid 4: Guns of the Patriots",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS3"),
		Genre:          "Stealth Action",
		TimeToBeat:     strPtr("18h"),
		TimeToBeatMain: floatPtr(18),
		AddedAt:        "2026-02-05T14:00:00Z",
	},
	{
		ID:             "ps-5",
		Title:          "Shadow of the Colossus",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS2"),
		Genre:          "Action Adventure",
		TimeToBeat:     strPtr("9h"),
		TimeToBeatMain: floatPtr(9),
		AddedAt:        "2026-01-22T10:00:00Z",
	},
	{
		ID:             "ps-6",
		Title:          "Silent Hill 2",
		Platform:       models.PlatformPlayStation,
		Subcategory:    strPtr("PS2"),
		Genre:          "Psychological Horror",
		TimeToBeat:     strPtr("8h"),
		TimeToBeatMain: floatPtr(8),
		AddedAt:        "2026-01-25T11:00:00Z",
	},

	// Nintendo Switch
	{
		ID:             "switch-1",
		Title:          "The Legend of Zelda: Tears of the Kingdom",
		Platform:       models.PlatformNintendoSwitch,
		Subcategory:    strPtr("Switch"),
		Genre:          "Action Adventure",
		TimeToBeat:     strPtr("59h"),
		TimeToBeatMain: floatPtr(59),
		AddedAt:        "2026-01-20T14:00:00Z",
	},
	{
		ID:             "switch-2",
		Title:          "Metroid Dread",
		Platform:       models.PlatformNintendoSwitch,
		Subcategory:    strPtr("Switch"),
		Genre:          "Metroidvania",
		TimeToBeat:     strPtr("10h"),
		TimeToBeatMain: floatPtr(10),
		AddedAt:        "2026-02-02T16:00:00Z",
	},
	{
		ID:             "switch-3",
		Title:          "Sea of Stars",
		Platform:       models.PlatformNintendoSwitch,
		Subcategory:    strPtr("Switch"),
		Genre:          "RPG",
		TimeToBeat:     strPtr("28h"),
		TimeToBeatMain: floatPtr(28),
		AddedAt:        "2026-03-05T11:20:00Z",
	},
	{
		ID:             "switch-4",
		Title:          "Super Mario Bros. Wonder",
		Platform:       models.PlatformNintendoSwitch,
		Subcategory:    strPtr("Switch"),
		Genre:          "Platformer",
		TimeToBeat:     strPtr("10h"),
		TimeToBeatMain: floatPtr(10),
		AddedAt:        "2026-01-25T13:40:00Z",
	},

	// Xbox
	{
		ID:             "xbox-1",
		Title:          "Avowed",
		Platform:       models.PlatformXbox,
		Subcategory:    strPtr("Xbox Series X/S"),
		Genre:          "Action RPG",
		TimeToBeat:     strPtr("25h"),
		TimeToBeatMain: floatPtr(25),
		AddedAt:        "2026-03-02T14:00:00Z",
	},
	{
		ID:             "xbox-2",
		Title:          "Hi-Fi RUSH",
		Platform:       models.PlatformXbox,
		Subcategory:    strPtr("Xbox Series X/S"),
		Genre:          "Rhythm Action",
		TimeToBeat:     strPtr("11h"),
		TimeToBeatMain: floatPtr(11),
		AddedAt:        "2026-01-30T17:30:00Z",
	},
	{
		ID:             "xbox-3",
		Title:          "Halo 3",
		Platform:       models.PlatformXbox,
		Subcategory:    strPtr("Xbox 360"),
		Genre:          "FPS",
		TimeToBeat:     strPtr("12h"),
		TimeToBeatMain: floatPtr(12),
		AddedAt:        "2026-01-08T15:00:00Z",
	},
	{
		ID:             "xbox-4",
		Title:          "Gears of War",
		Platform:       models.PlatformXbox,
		Subcategory:    strPtr("Xbox 360"),
		Genre:          "Third-Person Shooter",
		TimeToBeat:     strPtr("9h"),
		TimeToBeatMain: floatPtr(9),
		AddedAt:        "2026-01-15T16:00:00Z",
	},

	// Retro / Emulation
	{
		ID:             "retro-ps1-1",
		Title:          "Castlevania: Symphony of the Night",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("PS1"),
		Genre:          "Metroidvania",
		TimeToBeat:     strPtr("9h"),
		TimeToBeatMain: floatPtr(9),
		AddedAt:        "2026-01-15T11:00:00Z",
	},
	{
		ID:             "retro-n64-1",
		Title:          "The Legend of Zelda: Ocarina of Time",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("N64"),
		Genre:          "Action Adventure",
		TimeToBeat:     strPtr("30h"),
		TimeToBeatMain: floatPtr(30),
		AddedAt:        "2026-01-12T10:00:00Z",
	},
	{
		ID:             "retro-n64-2",
		Title:          "Super Mario 64",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("N64"),
		Genre:          "Platformer",
		TimeToBeat:     strPtr("12h"),
		TimeToBeatMain: floatPtr(12),
		AddedAt:        "2026-01-14T15:00:00Z",
	},
	{
		ID:             "retro-snes-1",
		Title:          "Chrono Trigger",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("SNES"),
		Genre:          "RPG",
		TimeToBeat:     strPtr("23h"),
		TimeToBeatMain: floatPtr(23),
		AddedAt:        "2026-01-01T00:00:00Z",
	},
	{
		ID:             "retro-snes-2",
		Title:          "Super Metroid",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("SNES"),
		Genre:          "Metroidvania",
		TimeToBeat:     strPtr("8h"),
		TimeToBeatMain: floatPtr(8),
		AddedAt:        "2026-01-03T14:00:00Z",
	},
	{
		ID:             "retro-snes-3",
		Title:          "EarthBound (Mother 2)",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("SNES"),
		Genre:          "RPG",
		TimeToBeat:     strPtr("28h"),
		TimeToBeatMain: floatPtr(28),
		AddedAt:        "2026-03-14T19:30:00Z",
	},
	{
		ID:             "retro-gen-1",
		Title:          "Sonic the Hedgehog 2",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("Genesis"),
		Genre:          "Platformer",
		TimeToBeat:     strPtr("3h"),
		TimeToBeatMain: floatPtr(3),
		AddedAt:        "2026-02-10T12:00:00Z",
	},
	{
		ID:             "retro-gen-2",
		Title:          "Streets of Rage 2",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("Genesis"),
		Genre:          "Beat 'em up",
		TimeToBeat:     strPtr("2h"),
		TimeToBeatMain: floatPtr(2),
		AddedAt:        "2026-02-12T16:00:00Z",
	},
	{
		ID:             "retro-nes-1",
		Title:          "Mega Man 2",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("NES"),
		Genre:          "Action Platformer",
		TimeToBeat:     strPtr("3h"),
		TimeToBeatMain: floatPtr(3),
		AddedAt:        "2026-02-01T10:00:00Z",
	},
	{
		ID:             "retro-ms-1",
		Title:          "Alex Kidd in Miracle World",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("Master System"),
		Genre:          "Platformer",
		TimeToBeat:     strPtr("2h"),
		TimeToBeatMain: floatPtr(2),
		AddedAt:        "2026-02-04T11:00:00Z",
	},
	{
		ID:             "retro-neo-1",
		Title:          "Metal Slug 3",
		Platform:       models.PlatformRetroEmulation,
		Subcategory:    strPtr("Neo Geo"),
		Genre:          "Run & Gun",
		TimeToBeat:     strPtr("2h"),
		TimeToBeatMain: floatPtr(2),
		AddedAt:        "2026-02-18T18:00:00Z",
	},
}

// SeedDatabase seeds the database with InitialSeedGames
func SeedDatabase(db *sql.DB, reset bool) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	if reset {
		if _, err := tx.Exec("DELETE FROM games"); err != nil {
			return 0, fmt.Errorf("failed to clear games table: %w", err)
		}
	}

	insertQuery := `
	INSERT INTO games (
		id, title, platform, subcategory, genre,
		time_to_beat, time_to_beat_main, added_at, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO NOTHING;
	`
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	insertedCount := 0

	for _, g := range InitialSeedGames {
		// Also check if title on same platform and subcat exists
		var existingID string
		checkQuery := `
			SELECT id FROM games 
			WHERE LOWER(TRIM(title)) = LOWER(TRIM(?)) 
			  AND platform = ? 
			  AND COALESCE(subcategory, '') = COALESCE(?, '')
			LIMIT 1;
		`
		err := tx.QueryRow(checkQuery, g.Title, g.Platform, g.Subcategory).Scan(&existingID)
		if err == nil {
			// Already exists
			continue
		} else if err != sql.ErrNoRows {
			return 0, fmt.Errorf("failed checking duplicate during seed: %w", err)
		}

		res, err := stmt.Exec(
			g.ID,
			g.Title,
			g.Platform,
			g.Subcategory,
			g.Genre,
			g.TimeToBeat,
			g.TimeToBeatMain,
			g.AddedAt,
			now,
			now,
		)
		if err != nil {
			return 0, fmt.Errorf("failed inserting game %s: %w", g.Title, err)
		}
		rows, _ := res.RowsAffected()
		if rows > 0 {
			insertedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit seed transaction: %w", err)
	}

	return insertedCount, nil
}
