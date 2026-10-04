package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"game-catalog/api/models"
)

var (
	ErrNotFound       = errors.New("game not found")
	ErrDuplicate      = errors.New("a game with this title, platform, and subcategory already exists")
	ErrDuplicateID    = errors.New("a game with this ID already exists")
	ErrDatabaseClosed = errors.New("database connection is closed")
)

// GameRepository defines database operations for the game catalog
type GameRepository interface {
	GetByID(ctx context.Context, id string) (*models.Game, error)
	List(ctx context.Context, query models.ListGamesQuery) ([]models.Game, int, error)
	Create(ctx context.Context, game *models.Game) error
	Update(ctx context.Context, game *models.Game) error
	Delete(ctx context.Context, id string) error
	BatchCreate(ctx context.Context, games []*models.Game) (int, []error)
	DeleteAll(ctx context.Context) (int64, error)
	GetStats(ctx context.Context) (*models.CatalogStats, error)
	GetPlatformStats(ctx context.Context) ([]models.PlatformInfo, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsDuplicate(ctx context.Context, title, platform string, subcategory *string, excludeID string) (bool, error)
	Count(ctx context.Context) (int, error)
}

type sqliteGameRepository struct {
	db *sql.DB
}

// NewGameRepository creates a new SQLite-backed game repository
func NewGameRepository(db *sql.DB) GameRepository {
	return &sqliteGameRepository{db: db}
}

func (r *sqliteGameRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM games").Scan(&count)
	return count, err
}

func (r *sqliteGameRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM games WHERE id = ?", id).Scan(&count)
	return count > 0, err
}

func (r *sqliteGameRepository) ExistsDuplicate(ctx context.Context, title, platform string, subcategory *string, excludeID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM games 
		WHERE LOWER(TRIM(title)) = LOWER(TRIM(?))
		  AND platform = ?
		  AND COALESCE(subcategory, '') = COALESCE(?, '')
	`
	args := []interface{}{title, platform, subcategory}
	if excludeID != "" {
		query += " AND id != ?"
		args = append(args, excludeID)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count > 0, err
}

func (r *sqliteGameRepository) GetByID(ctx context.Context, id string) (*models.Game, error) {
	query := `
		SELECT id, title, platform, subcategory, genre,
		       time_to_beat, time_to_beat_main, added_at, created_at, updated_at
		FROM games
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var g models.Game
	var subcat, ttb sql.NullString
	var ttbMain sql.NullFloat64

	err := row.Scan(
		&g.ID,
		&g.Title,
		&g.Platform,
		&subcat,
		&g.Genre,
		&ttb,
		&ttbMain,
		&g.AddedAt,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if subcat.Valid {
		g.Subcategory = &subcat.String
	}
	if ttb.Valid {
		g.TimeToBeat = &ttb.String
	}
	if ttbMain.Valid {
		g.TimeToBeatMain = &ttbMain.Float64
	}

	return &g, nil
}

func (r *sqliteGameRepository) List(ctx context.Context, q models.ListGamesQuery) ([]models.Game, int, error) {
	var whereClauses []string
	var args []interface{}

	// Platform filter
	if q.Platform != "" && !strings.EqualFold(q.Platform, "All") {
		normPlatform, ok := models.NormalizePlatform(q.Platform)
		if ok {
			whereClauses = append(whereClauses, "platform = ?")
			args = append(args, normPlatform)
		} else {
			// Filtering by an unknown platform returns empty list
			whereClauses = append(whereClauses, "1 = 0")
		}
	}

	// Subcategory filter
	if q.Subcategory != "" {
		whereClauses = append(whereClauses, "LOWER(subcategory) = LOWER(?)")
		args = append(args, strings.TrimSpace(q.Subcategory))
	}

	// Genre filter
	if q.Genre != "" {
		whereClauses = append(whereClauses, "LOWER(genre) = LOWER(?)")
		args = append(args, strings.TrimSpace(q.Genre))
	}

	// Search filter across title, genre, subcategory
	if strings.TrimSpace(q.Search) != "" {
		whereClauses = append(whereClauses, "(title LIKE ? ESCAPE '\\' OR genre LIKE ? ESCAPE '\\' OR subcategory LIKE ? ESCAPE '\\')")
		// Escape special LIKE chars in search pattern if needed, or use simple pattern
		escaped := strings.ReplaceAll(q.Search, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "%", "\\%")
		escaped = strings.ReplaceAll(escaped, "_", "\\_")
		escapedPattern := "%" + strings.TrimSpace(escaped) + "%"
		args = append(args, escapedPattern, escapedPattern, escapedPattern)
	}

	// Playtime bounds
	if q.MinHours != nil {
		whereClauses = append(whereClauses, "time_to_beat_main >= ?")
		args = append(args, *q.MinHours)
	}
	if q.MaxHours != nil {
		whereClauses = append(whereClauses, "time_to_beat_main <= ?")
		args = append(args, *q.MaxHours)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total matching
	countQuery := "SELECT COUNT(*) FROM games" + whereSQL
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count games: %w", err)
	}

	// Sorting
	sortByCol := "title"
	switch strings.ToLower(q.SortBy) {
	case "hours", "time_to_beat_main", "timetobeatmain":
		sortByCol = "time_to_beat_main"
	case "added_at", "addedat":
		sortByCol = "added_at"
	case "created_at", "createdat":
		sortByCol = "created_at"
	case "title":
		sortByCol = "title"
	}

	orderDir := "ASC"
	if strings.EqualFold(q.Order, "desc") {
		orderDir = "DESC"
	} else if strings.EqualFold(q.Order, "asc") {
		orderDir = "ASC"
	} else {
		// Default directions: hours and added_at desc by default, title asc
		if sortByCol == "time_to_beat_main" || sortByCol == "added_at" || sortByCol == "created_at" {
			orderDir = "DESC"
		} else {
			orderDir = "ASC"
		}
	}

	// For hours sorting, put NULLs last
	orderBySQL := ""
	if sortByCol == "time_to_beat_main" {
		if orderDir == "DESC" {
			orderBySQL = " ORDER BY time_to_beat_main IS NULL ASC, time_to_beat_main DESC, title ASC"
		} else {
			orderBySQL = " ORDER BY time_to_beat_main IS NULL ASC, time_to_beat_main ASC, title ASC"
		}
	} else {
		orderBySQL = fmt.Sprintf(" ORDER BY %s %s", sortByCol, orderDir)
	}

	// Pagination
	pageSize := q.PageSize
	if pageSize <= 0 {
		pageSize = 50
	} else if pageSize > 100 {
		pageSize = 100
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`
		SELECT id, title, platform, subcategory, genre,
		       time_to_beat, time_to_beat_main, added_at, created_at, updated_at
		FROM games
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereSQL, orderBySQL)

	fetchArgs := append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, fetchArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query games: %w", err)
	}
	defer rows.Close()

	games := make([]models.Game, 0)
	for rows.Next() {
		var g models.Game
		var subcat, ttb sql.NullString
		var ttbMain sql.NullFloat64

		if err := rows.Scan(
			&g.ID,
			&g.Title,
			&g.Platform,
			&subcat,
			&g.Genre,
			&ttb,
			&ttbMain,
			&g.AddedAt,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan game row: %w", err)
		}

		if subcat.Valid {
			g.Subcategory = &subcat.String
		}
		if ttb.Valid {
			g.TimeToBeat = &ttb.String
		}
		if ttbMain.Valid {
			g.TimeToBeatMain = &ttbMain.Float64
		}

		games = append(games, g)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during row iteration: %w", err)
	}

	return games, total, nil
}

func (r *sqliteGameRepository) Create(ctx context.Context, game *models.Game) error {
	// Check duplicate ID
	existsID, err := r.ExistsByID(ctx, game.ID)
	if err != nil {
		return err
	}
	if existsID {
		return ErrDuplicateID
	}

	// Check duplicate title+platform+subcategory
	dup, err := r.ExistsDuplicate(ctx, game.Title, game.Platform, game.Subcategory, "")
	if err != nil {
		return err
	}
	if dup {
		return ErrDuplicate
	}

	query := `
		INSERT INTO games (
			id, title, platform, subcategory, genre,
			time_to_beat, time_to_beat_main, added_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		game.ID,
		game.Title,
		game.Platform,
		game.Subcategory,
		game.Genre,
		game.TimeToBeat,
		game.TimeToBeatMain,
		game.AddedAt,
		game.CreatedAt,
		game.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicate
		}
		return err
	}

	return nil
}

func (r *sqliteGameRepository) Update(ctx context.Context, game *models.Game) error {
	// Ensure exists
	exists, err := r.ExistsByID(ctx, game.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	// Check duplicate title+platform+subcategory (excluding this ID)
	dup, err := r.ExistsDuplicate(ctx, game.Title, game.Platform, game.Subcategory, game.ID)
	if err != nil {
		return err
	}
	if dup {
		return ErrDuplicate
	}

	query := `
		UPDATE games SET
			title = ?,
			platform = ?,
			subcategory = ?,
			genre = ?,
			time_to_beat = ?,
			time_to_beat_main = ?,
			added_at = ?,
			updated_at = ?
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query,
		game.Title,
		game.Platform,
		game.Subcategory,
		game.Genre,
		game.TimeToBeat,
		game.TimeToBeatMain,
		game.AddedAt,
		game.UpdatedAt,
		game.ID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicate
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sqliteGameRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM games WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqliteGameRepository) BatchCreate(ctx context.Context, games []*models.Game) (int, []error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, []error{err}
	}
	defer tx.Rollback()

	insertQuery := `
		INSERT INTO games (
			id, title, platform, subcategory, genre,
			time_to_beat, time_to_beat_main, added_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return 0, []error{err}
	}
	defer stmt.Close()

	var inserted int
	var errs []error

	for _, g := range games {
		// Check for duplicate in batch / db
		var existingID string
		checkQuery := `
			SELECT id FROM games 
			WHERE (id = ?) OR (
				LOWER(TRIM(title)) = LOWER(TRIM(?)) 
				AND platform = ? 
				AND COALESCE(subcategory, '') = COALESCE(?, '')
			) LIMIT 1;
		`
		err := tx.QueryRowContext(ctx, checkQuery, g.ID, g.Title, g.Platform, g.Subcategory).Scan(&existingID)
		if err == nil {
			errs = append(errs, fmt.Errorf("game '%s' (%s) already exists", g.Title, g.Platform))
			continue
		} else if err != sql.ErrNoRows {
			errs = append(errs, err)
			continue
		}

		_, err = stmt.ExecContext(ctx,
			g.ID,
			g.Title,
			g.Platform,
			g.Subcategory,
			g.Genre,
			g.TimeToBeat,
			g.TimeToBeatMain,
			g.AddedAt,
			g.CreatedAt,
			g.UpdatedAt,
		)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to insert '%s': %w", g.Title, err))
			continue
		}
		inserted++
	}

	if err := tx.Commit(); err != nil {
		return 0, []error{err}
	}

	return inserted, errs
}

func (r *sqliteGameRepository) DeleteAll(ctx context.Context) (int64, error) {
	res, err := r.db.ExecContext(ctx, "DELETE FROM games")
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *sqliteGameRepository) GetStats(ctx context.Context) (*models.CatalogStats, error) {
	stats := &models.CatalogStats{
		PlatformBreakdown: make([]models.PlatformStat, 0),
		GenreBreakdown:    make([]models.GenreStat, 0),
	}

	// 1. Overall counts & hours
	overallQuery := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(time_to_beat_main), 0),
			COALESCE(AVG(time_to_beat_main), 0)
		FROM games
	`
	if err := r.db.QueryRowContext(ctx, overallQuery).Scan(
		&stats.TotalGames,
		&stats.TotalPlaytimeHours,
		&stats.AveragePlaytimeHours,
	); err != nil {
		return nil, fmt.Errorf("failed to query overall stats: %w", err)
	}

	if stats.TotalGames == 0 {
		return stats, nil
	}

	// 2. Shortest game (with time_to_beat_main > 0)
	shortestQuery := `
		SELECT id, title, platform, time_to_beat_main
		FROM games
		WHERE time_to_beat_main IS NOT NULL AND time_to_beat_main > 0
		ORDER BY time_to_beat_main ASC
		LIMIT 1
	`
	var shortest models.GameSummary
	err := r.db.QueryRowContext(ctx, shortestQuery).Scan(
		&shortest.ID,
		&shortest.Title,
		&shortest.Platform,
		&shortest.Hours,
	)
	if err == nil {
		stats.ShortestGame = &shortest
	}

	// 3. Longest game
	longestQuery := `
		SELECT id, title, platform, time_to_beat_main
		FROM games
		WHERE time_to_beat_main IS NOT NULL
		ORDER BY time_to_beat_main DESC
		LIMIT 1
	`
	var longest models.GameSummary
	err = r.db.QueryRowContext(ctx, longestQuery).Scan(
		&longest.ID,
		&longest.Title,
		&longest.Platform,
		&longest.Hours,
	)
	if err == nil {
		stats.LongestGame = &longest
	}

	// 4. Platform breakdown
	platformQuery := `
		SELECT 
			platform,
			COUNT(*),
			COALESCE(SUM(time_to_beat_main), 0),
			COALESCE(AVG(time_to_beat_main), 0)
		FROM games
		GROUP BY platform
		ORDER BY COUNT(*) DESC
	`
	pRows, err := r.db.QueryContext(ctx, platformQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform breakdown: %w", err)
	}

	var platformStats []models.PlatformStat
	for pRows.Next() {
		var ps models.PlatformStat
		ps.Subcategories = make([]models.SubcategoryStat, 0)
		if err := pRows.Scan(&ps.Platform, &ps.Count, &ps.TotalHours, &ps.AverageHours); err != nil {
			pRows.Close()
			return nil, err
		}
		platformStats = append(platformStats, ps)
	}
	pRows.Close()

	// Query subcategories in a single query to avoid connection starvation
	subQuery := `
		SELECT platform, COALESCE(subcategory, 'Uncategorized'), COUNT(*)
		FROM games
		GROUP BY platform, subcategory
		ORDER BY COUNT(*) DESC
	`
	sRows, err := r.db.QueryContext(ctx, subQuery)
	if err == nil {
		subcatMap := make(map[string][]models.SubcategoryStat)
		for sRows.Next() {
			var plat, sub string
			var cnt int
			if err := sRows.Scan(&plat, &sub, &cnt); err == nil {
				subcatMap[plat] = append(subcatMap[plat], models.SubcategoryStat{
					Subcategory: sub,
					Count:       cnt,
				})
			}
		}
		sRows.Close()

		for i := range platformStats {
			if subs, ok := subcatMap[platformStats[i].Platform]; ok {
				platformStats[i].Subcategories = subs
			}
		}
	}

	stats.PlatformBreakdown = platformStats

	// 5. Genre breakdown
	genreQuery := `
		SELECT genre, COUNT(*)
		FROM games
		GROUP BY genre
		ORDER BY COUNT(*) DESC
	`
	gRows, err := r.db.QueryContext(ctx, genreQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query genre breakdown: %w", err)
	}
	defer gRows.Close()

	for gRows.Next() {
		var gs models.GenreStat
		if err := gRows.Scan(&gs.Genre, &gs.Count); err != nil {
			return nil, err
		}
		stats.GenreBreakdown = append(stats.GenreBreakdown, gs)
	}

	return stats, nil
}

func (r *sqliteGameRepository) GetPlatformStats(ctx context.Context) ([]models.PlatformInfo, error) {
	result := make([]models.PlatformInfo, 0, len(models.ValidPlatforms))

	for _, platform := range models.ValidPlatforms {
		var count int
		err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM games WHERE platform = ?", platform).Scan(&count)
		if err != nil {
			return nil, err
		}

		subcats := models.ValidSubcategories[platform]
		if subcats == nil {
			subcats = []string{}
		}

		result = append(result, models.PlatformInfo{
			Name:          platform,
			Subcategories: subcats,
			GameCount:     count,
		})
	}

	return result, nil
}
