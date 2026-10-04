package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"game-catalog/api/config"
	"game-catalog/api/database"
	"game-catalog/api/models"
	"game-catalog/api/server"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestServer(t *testing.T) (*gin.Engine, *sql.DB) {
	t.Helper()
	cfg := &config.Config{
		Port:               "8080",
		GinMode:            "test",
		DatabasePath:       fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano()),
		CORSAllowedOrigins: []string{"*"},
		AutoSeed:           false,
		MaxBodySizeBytes:   1024 * 1024,
	}

	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		t.Fatalf("failed to init test database: %v", err)
	}

	router := server.SetupRouter(cfg, db)
	return router, db
}

func performRequest(r http.Handler, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = bytes.NewBufferString(v)
		case []byte:
			reqBody = bytes.NewBuffer(v)
		default:
			jsonBytes, _ := json.Marshal(v)
			reqBody = bytes.NewBuffer(jsonBytes)
		}
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	if body != nil && headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealthAndPing(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Health check endpoint", func(t *testing.T) {
		w := performRequest(router, "GET", "/health", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		if res["status"] != "healthy" {
			t.Errorf("expected healthy status, got %v", res["status"])
		}
		if res["database"] != "connected" {
			t.Errorf("expected database connected, got %v", res["database"])
		}
	})

	t.Run("Ping endpoint", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/ping", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		if res["message"] != "pong" {
			t.Errorf("expected pong message, got %v", res["message"])
		}
	})

	t.Run("Unknown route returns standardized 404 JSON", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/nonexistent", nil, nil)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
		var errRes map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &errRes)
		if errRes["error"] == nil {
			t.Error("expected standardized error payload in 404 response")
		}
	})
}

func TestPlatformEndpoints(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	w := performRequest(router, "GET", "/api/v1/platforms", nil, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res struct {
		Platforms []models.PlatformInfo `json:"platforms"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(res.Platforms) != len(models.ValidPlatforms) {
		t.Errorf("expected %d platforms, got %d", len(models.ValidPlatforms), len(res.Platforms))
	}
}

func TestSeedEndpoint(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	// 1. Initial seed
	w1 := performRequest(router, "POST", "/api/v1/seed", nil, nil)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}
	var res1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &res1)
	inserted1 := int(res1["inserted"].(float64))
	if inserted1 == 0 {
		t.Errorf("expected initial seed to insert games, got 0")
	}

	// 2. Second seed without reset should be idempotent
	w2 := performRequest(router, "POST", "/api/v1/seed", nil, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var res2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &res2)
	inserted2 := int(res2["inserted"].(float64))
	if inserted2 != 0 {
		t.Errorf("expected idempotent seed to insert 0 new games, got %d", inserted2)
	}

	// 3. Seed with reset=true should reseed
	w3 := performRequest(router, "POST", "/api/v1/seed?reset=true", nil, nil)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w3.Code)
	}
	var res3 map[string]interface{}
	json.Unmarshal(w3.Body.Bytes(), &res3)
	inserted3 := int(res3["inserted"].(float64))
	if inserted3 != inserted1 {
		t.Errorf("expected reset seed to insert %d games, got %d", inserted1, inserted3)
	}
}

func TestGameCreationConstraints(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Valid game creation with all fields", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":          "Elden Ring",
			"platform":       "PC",
			"subcategory":    "Steam",
			"genre":          "Action RPG",
			"timeToBeat":     "58h",
			"timeToBeatMain": 58,
			"addedAt":        "2026-02-01T12:00:00Z",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		if location == "" {
			t.Error("expected Location header in 201 response")
		}

		var created models.Game
		json.Unmarshal(w.Body.Bytes(), &created)
		if created.Title != "Elden Ring" {
			t.Errorf("expected title 'Elden Ring', got '%s'", created.Title)
		}
		if created.TimeToBeatMain == nil || *created.TimeToBeatMain != 58 {
			t.Errorf("expected hours 58, got %v", created.TimeToBeatMain)
		}
	})

	t.Run("Valid game with Unicode and emojis in title", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":    "NieR:Automata™ ⚔️ (Game of the YoRHa Edition)",
			"platform": "PlayStation",
			"genre":    "Action JRPG",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created for unicode title, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Missing title rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"platform": "PC",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Whitespace title rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":    "   \t\n  ",
			"platform": "PC",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Overly long title rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":    string(make([]byte, 256)),
			"platform": "PC",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Invalid platform rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":    "Valid Game",
			"platform": "GameCube",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Incompatible subcategory rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "Valid Game",
			"platform":    "PC",
			"subcategory": "PS5", // PS5 is PlayStation, not PC
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Negative hours rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":          "Valid Game",
			"platform":       "PC",
			"timeToBeatMain": -10,
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Duplicate game rejected with 409 Conflict", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "Elden Ring",
			"platform":    "PC",
			"subcategory": "Steam",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusConflict {
			t.Errorf("expected 409 Conflict for duplicate game, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Same title on different platform allowed", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "Elden Ring",
			"platform":    "Xbox",
			"subcategory": "Xbox Series X/S",
		}
		w := performRequest(router, "POST", "/api/v1/games", payload, nil)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201 Created for cross-platform release, got %d", w.Code)
		}
	})

	t.Run("Custom duplicate ID rejected with 409 Conflict", func(t *testing.T) {
		payload1 := map[string]interface{}{
			"id":       "my-custom-id",
			"title":    "Unique Game 1",
			"platform": "PC",
		}
		w1 := performRequest(router, "POST", "/api/v1/games", payload1, nil)
		if w1.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d", w1.Code)
		}

		payload2 := map[string]interface{}{
			"id":       "my-custom-id",
			"title":    "Unique Game 2",
			"platform": "Xbox",
		}
		w2 := performRequest(router, "POST", "/api/v1/games", payload2, nil)
		if w2.Code != http.StatusConflict {
			t.Errorf("expected 409 Conflict for duplicate ID, got %d", w2.Code)
		}
	})

	t.Run("Malformed JSON rejected with 400", func(t *testing.T) {
		w := performRequest(router, "POST", "/api/v1/games", "not a json string", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Non-JSON content type rejected with 415", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "text/plain"}
		w := performRequest(router, "POST", "/api/v1/games", "text body", headers)
		if w.Code != http.StatusUnsupportedMediaType {
			t.Errorf("expected 415 Unsupported Media Type, got %d", w.Code)
		}
	})
}

func TestGameCRUDAndEdgeCases(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	// Create test game
	createPayload := map[string]interface{}{
		"id":             "crud-game-1",
		"title":          "Super Metroid",
		"platform":       "Retro / Emulation",
		"subcategory":    "SNES",
		"genre":          "Metroidvania",
		"timeToBeatMain": 8,
		"timeToBeat":     "8h",
	}
	wCreate := performRequest(router, "POST", "/api/v1/games", createPayload, nil)
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("failed to create: %d", wCreate.Code)
	}

	// 1. Get Game by ID
	t.Run("Get existing game", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games/crud-game-1", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var g models.Game
		json.Unmarshal(w.Body.Bytes(), &g)
		if g.Title != "Super Metroid" {
			t.Errorf("expected Super Metroid, got %s", g.Title)
		}
	})

	t.Run("Get non-existent game returns 404", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games/does-not-exist", nil, nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// 2. Full Update (PUT)
	t.Run("PUT full update existing game", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"title":          "Super Metroid: Redesign",
			"platform":       "Retro / Emulation",
			"subcategory":    "SNES",
			"genre":          "Metroidvania",
			"timeToBeatMain": 15,
		}
		w := performRequest(router, "PUT", "/api/v1/games/crud-game-1", updatePayload, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		var updated models.Game
		json.Unmarshal(w.Body.Bytes(), &updated)
		if updated.Title != "Super Metroid: Redesign" || *updated.TimeToBeatMain != 15 {
			t.Errorf("unexpected updated values: %+v", updated)
		}
	})

	t.Run("PUT non-existent game returns 404", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"title":    "Non-existent",
			"platform": "PC",
		}
		w := performRequest(router, "PUT", "/api/v1/games/fake-id", updatePayload, nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// 3. Partial Update (PATCH)
	t.Run("PATCH partial update (hours only)", func(t *testing.T) {
		patchPayload := map[string]interface{}{
			"timeToBeatMain": 20,
		}
		w := performRequest(router, "PATCH", "/api/v1/games/crud-game-1", patchPayload, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		var patched models.Game
		json.Unmarshal(w.Body.Bytes(), &patched)
		if *patched.TimeToBeatMain != 20 {
			t.Errorf("expected 20 hours, got %v", patched.TimeToBeatMain)
		}
		// Preserved fields
		if patched.Title != "Super Metroid: Redesign" {
			t.Errorf("expected preserved title, got %s", patched.Title)
		}
	})

	t.Run("PATCH non-existent game returns 404", func(t *testing.T) {
		patchPayload := map[string]interface{}{"title": "New"}
		w := performRequest(router, "PATCH", "/api/v1/games/fake-id", patchPayload, nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// 4. Delete
	t.Run("DELETE existing game", func(t *testing.T) {
		w := performRequest(router, "DELETE", "/api/v1/games/crud-game-1", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Verify deletion
		wCheck := performRequest(router, "GET", "/api/v1/games/crud-game-1", nil, nil)
		if wCheck.Code != http.StatusNotFound {
			t.Errorf("expected 404 after deletion, got %d", wCheck.Code)
		}

		// Second delete returns 404
		wDouble := performRequest(router, "DELETE", "/api/v1/games/crud-game-1", nil, nil)
		if wDouble.Code != http.StatusNotFound {
			t.Errorf("expected 404 on re-deleting, got %d", wDouble.Code)
		}
	})
}

func TestListQueryAndFiltering(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	// Seed database
	performRequest(router, "POST", "/api/v1/seed", nil, nil)

	t.Run("Default list returns games with pagination metadata", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Total <= 0 {
			t.Errorf("expected total > 0, got %d", res.Total)
		}
		if res.Page != 1 {
			t.Errorf("expected page 1, got %d", res.Page)
		}
		if len(res.Data) != res.Total {
			t.Errorf("expected %d items on page 1, got %d", res.Total, len(res.Data))
		}
	})

	t.Run("Filter by platform", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?platform=PC", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		for _, g := range res.Data {
			if g.Platform != "PC" {
				t.Errorf("expected PC game, got %s", g.Platform)
			}
		}
	})

	t.Run("Search query", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?search=Zelda", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Total == 0 {
			t.Error("expected to find Zelda games")
		}
		for _, g := range res.Data {
			if !bytes.Contains([]byte(g.Title), []byte("Zelda")) {
				t.Errorf("expected title to contain Zelda, got %s", g.Title)
			}
		}
	})

	t.Run("Invalid pagination parameters return 400", func(t *testing.T) {
		wPage := performRequest(router, "GET", "/api/v1/games?page=-1", nil, nil)
		if wPage.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for negative page, got %d", wPage.Code)
		}

		wLimit := performRequest(router, "GET", "/api/v1/games?limit=abc", nil, nil)
		if wLimit.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid limit, got %d", wLimit.Code)
		}

		wMin := performRequest(router, "GET", "/api/v1/games?min_hours=-5", nil, nil)
		if wMin.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for negative min_hours, got %d", wMin.Code)
		}
	})
}

func TestBatchAndClearEndpoints(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Batch create games", func(t *testing.T) {
		batch := []map[string]interface{}{
			{
				"title":          "Game Alpha",
				"platform":       "PC",
				"genre":          "Action",
				"timeToBeatMain": 10,
			},
			{
				"title":          "Game Beta",
				"platform":       "Nintendo Switch",
				"genre":          "RPG",
				"timeToBeatMain": 20,
			},
		}
		w := performRequest(router, "POST", "/api/v1/games/batch", batch, nil)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
		}
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		if int(res["inserted"].(float64)) != 2 {
			t.Errorf("expected 2 inserted, got %v", res["inserted"])
		}
	})

	t.Run("Empty batch rejected with 400", func(t *testing.T) {
		w := performRequest(router, "POST", "/api/v1/games/batch", []interface{}{}, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for empty batch, got %d", w.Code)
		}
	})

	t.Run("Clear without confirmation fails with 400", func(t *testing.T) {
		w := performRequest(router, "DELETE", "/api/v1/games", nil, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 without ?confirm=true, got %d", w.Code)
		}
	})

	t.Run("Clear with confirmation succeeds with 200", func(t *testing.T) {
		w := performRequest(router, "DELETE", "/api/v1/games?confirm=true", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Verify empty
		wList := performRequest(router, "GET", "/api/v1/games", nil, nil)
		var res models.PaginatedResponse
		json.Unmarshal(wList.Body.Bytes(), &res)
		if res.Total != 0 || len(res.Data) != 0 {
			t.Errorf("expected 0 games after clear, got %d", res.Total)
		}
	})
}

func TestStatsEndpoint(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Stats on empty catalog", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/stats", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var s models.CatalogStats
		json.Unmarshal(w.Body.Bytes(), &s)
		if s.TotalGames != 0 || s.TotalPlaytimeHours != 0 {
			t.Errorf("expected zeroes for empty stats, got %+v", s)
		}
	})

	t.Run("Stats after seeding", func(t *testing.T) {
		performRequest(router, "POST", "/api/v1/seed", nil, nil)

		w := performRequest(router, "GET", "/api/v1/stats", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var s models.CatalogStats
		json.Unmarshal(w.Body.Bytes(), &s)
		if s.TotalGames == 0 {
			t.Error("expected non-zero total games in stats")
		}
		if s.TotalPlaytimeHours <= 0 {
			t.Error("expected non-zero playtime in stats")
		}
		if s.ShortestGame == nil || s.LongestGame == nil {
			t.Error("expected shortest and longest game summaries")
		}
		if len(s.PlatformBreakdown) == 0 {
			t.Error("expected platform breakdown")
		}
	})
}

func TestCORSOptionsPreflight(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	headers := map[string]string{
		"Origin":                         "http://localhost:3000",
		"Access-Control-Request-Method":  "POST",
		"Access-Control-Request-Headers": "Content-Type",
	}
	w := performRequest(router, "OPTIONS", "/api/v1/games", nil, headers)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for OPTIONS preflight, got %d", w.Code)
	}
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" && allowOrigin != "http://localhost:3000" {
		t.Errorf("expected valid CORS Allow-Origin header, got %s", allowOrigin)
	}
}

func TestConcurrentReadsAndWrites(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	performRequest(router, "POST", "/api/v1/seed", nil, nil)

	var wg sync.WaitGroup
	errChan := make(chan error, 50)
	workers := 20

	// Run concurrent read and write operations
	for i := 0; i < workers; i++ {
		wg.Add(2)

		// Concurrent reader
		go func(workerID int) {
			defer wg.Done()
			w := performRequest(router, "GET", "/api/v1/games?limit=5", nil, nil)
			if w.Code != http.StatusOK {
				errChan <- fmt.Errorf("worker %d read failed with %d: %s", workerID, w.Code, w.Body.String())
			}
		}(i)

		// Concurrent writer
		go func(workerID int) {
			defer wg.Done()
			payload := map[string]interface{}{
				"title":          fmt.Sprintf("Concurrent Game %d", workerID),
				"platform":       "PC",
				"genre":          "Concurrency Test",
				"timeToBeatMain": workerID + 1,
			}
			w := performRequest(router, "POST", "/api/v1/games", payload, nil)
			if w.Code != http.StatusCreated {
				errChan <- fmt.Errorf("worker %d write failed with %d: %s", workerID, w.Code, w.Body.String())
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrency error: %v", err)
	}
}

func TestGameUpdateEdgeCasesAndConstraints(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	// Create 2 test games
	g1Payload := map[string]interface{}{
		"id":          "up-game-1",
		"title":       "Hollow Knight",
		"platform":    "PC",
		"subcategory": "Steam",
		"genre":       "Metroidvania",
	}
	g2Payload := map[string]interface{}{
		"id":          "up-game-2",
		"title":       "Shovel Knight",
		"platform":    "Nintendo Switch",
		"subcategory": "Switch",
		"genre":       "Platformer",
	}
	performRequest(router, "POST", "/api/v1/games", g1Payload, nil)
	performRequest(router, "POST", "/api/v1/games", g2Payload, nil)

	t.Run("PUT with invalid JSON rejected with 400", func(t *testing.T) {
		w := performRequest(router, "PUT", "/api/v1/games/up-game-1", "{invalid-json", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PUT with invalid platform rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":    "Hollow Knight",
			"platform": "Atari 2600",
		}
		w := performRequest(router, "PUT", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PUT with incompatible subcategory rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "Hollow Knight",
			"platform":    "PC",
			"subcategory": "PS5",
		}
		w := performRequest(router, "PUT", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PUT causing duplicate collision with another game rejected with 409", func(t *testing.T) {
		// Try updating game 1 to match game 2's title, platform, and subcat
		payload := map[string]interface{}{
			"title":       "Shovel Knight",
			"platform":    "Nintendo Switch",
			"subcategory": "Switch",
			"genre":       "Platformer",
		}
		w := performRequest(router, "PUT", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusConflict {
			t.Errorf("expected 409 Conflict, got %d", w.Code)
		}
	})

	t.Run("PATCH with invalid JSON rejected with 400", func(t *testing.T) {
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", "bad-payload", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PATCH with empty title rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{"title": "   "}
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PATCH with invalid platform rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{"platform": "UnknownPlatform"}
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PATCH with invalid subcategory rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{"subcategory": "InvalidSubcategory"}
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PATCH with negative hours rejected with 400", func(t *testing.T) {
		payload := map[string]interface{}{"timeToBeatMain": -25}
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("PATCH causing duplicate collision rejected with 409", func(t *testing.T) {
		title := "Shovel Knight"
		plat := "Nintendo Switch"
		sub := "Switch"
		payload := map[string]interface{}{
			"title":       title,
			"platform":    plat,
			"subcategory": sub,
		}
		w := performRequest(router, "PATCH", "/api/v1/games/up-game-1", payload, nil)
		if w.Code != http.StatusConflict {
			t.Errorf("expected 409 Conflict, got %d", w.Code)
		}
	})
}

func TestBatchCreateEdgeCases(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Batch create with invalid JSON body rejected with 400", func(t *testing.T) {
		w := performRequest(router, "POST", "/api/v1/games/batch", "invalid-batch", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Batch create with item validation failure rejected with 400", func(t *testing.T) {
		batch := []map[string]interface{}{
			{"title": "Valid Game", "platform": "PC"},
			{"title": "", "platform": "PC"}, // Invalid item
		}
		w := performRequest(router, "POST", "/api/v1/games/batch", batch, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Batch create exceeding 500 items rejected with 400", func(t *testing.T) {
		largeBatch := make([]map[string]interface{}, 501)
		for i := 0; i < 501; i++ {
			largeBatch[i] = map[string]interface{}{
				"title":    fmt.Sprintf("Batch Game %d", i),
				"platform": "PC",
			}
		}
		w := performRequest(router, "POST", "/api/v1/games/batch", largeBatch, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 BATCH_TOO_LARGE, got %d", w.Code)
		}
	})

	t.Run("Batch create with duplicate item produces 207 MultiStatus", func(t *testing.T) {

		// First insert one game
		performRequest(router, "POST", "/api/v1/games", map[string]interface{}{
			"title":    "Existing Game",
			"platform": "PC",
		}, nil)

		// Now send a batch with 1 new and 1 duplicate
		batch := []map[string]interface{}{
			{"title": "Brand New Game", "platform": "PC"},
			{"title": "Existing Game", "platform": "PC"},
		}
		w := performRequest(router, "POST", "/api/v1/games/batch", batch, nil)
		if w.Code != http.StatusMultiStatus {
			t.Errorf("expected 207 MultiStatus, got %d", w.Code)
		}
	})
}

func TestGameListEdgeCases(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	// Seed data
	performRequest(router, "POST", "/api/v1/seed", nil, nil)

	t.Run("Filter by subcategory and genre combined", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?platform=PC&subcategory=Steam&genre=Metroidvania", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Total < 1 {
			t.Errorf("expected at least 1 match for Hollow Knight Silksong, got %d", res.Total)
		}
	})

	t.Run("Sort by added_at desc", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?sort_by=added_at&order=desc", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Sort by invalid field defaults gracefully", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?sort_by=random_column", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Page beyond max pages returns empty data array", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?page=999&limit=10", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) != 0 {
			t.Errorf("expected 0 items on page 999, got %d", len(res.Data))
		}
		if res.Total == 0 {
			t.Error("expected non-zero total items")
		}
	})

	t.Run("Query using q= parameter instead of search=", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/games?q=Cyberpunk", nil, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var res models.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Total != 1 {
			t.Errorf("expected 1 result for Cyberpunk, got %d", res.Total)
		}
	})
}

func TestHealthDegradedWhenDBClosed(t *testing.T) {
	router, db := setupTestServer(t)
	// Close database immediately to test degraded health check
	db.Close()

	w := performRequest(router, "GET", "/health", nil, nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 Service Unavailable when DB is closed, got %d", w.Code)
	}

	wV1 := performRequest(router, "GET", "/api/v1/health", nil, nil)
	if wV1.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 on /api/v1/health when DB is closed, got %d", wV1.Code)
	}
}

func TestHLTBEndpoint(t *testing.T) {
	router, db := setupTestServer(t)
	defer db.Close()

	t.Run("Missing query returns 400", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/hltb", nil, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Empty query returns 400", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/hltb?q=", nil, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Punctuation only returns found false", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/hltb?q=---", nil, nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", w.Code)
		}
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		if res["found"] != false {
			t.Errorf("expected found: false, got %v", res["found"])
		}
	})
}


