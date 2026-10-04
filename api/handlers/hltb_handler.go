package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	hltbTokenLock  sync.RWMutex
	hltbCachedToken string
	hltbTokenExpiry time.Time
	hltbClient     = &http.Client{Timeout: 10 * time.Second}
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
)

const hltbUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"

type hltbInitResponse struct {
	Token string `json:"token"`
}

type hltbGameItem struct {
	GameID       int    `json:"game_id"`
	GameName     string `json:"game_name"`
	CompMain     int    `json:"comp_main"` // in seconds
	CompPlus     int    `json:"comp_plus"`
	Comp100      int    `json:"comp_100"`
	CompAll      int    `json:"comp_all"`
	ReleaseWorld int    `json:"release_world"`
}

type hltbSearchResponse struct {
	Data []hltbGameItem `json:"data"`
}

// HLTBHandler handles HowLongToBeat proxy requests
type HLTBHandler struct{}

func NewHLTBHandler() *HLTBHandler {
	return &HLTBHandler{}
}

func getHLTBToken(forceRefresh bool) (string, error) {
	hltbTokenLock.RLock()
	if !forceRefresh && hltbCachedToken != "" && time.Now().Before(hltbTokenExpiry) {
		token := hltbCachedToken
		hltbTokenLock.RUnlock()
		return token, nil
	}
	hltbTokenLock.RUnlock()

	hltbTokenLock.Lock()
	defer hltbTokenLock.Unlock()

	// Double-check after acquiring write lock
	if !forceRefresh && hltbCachedToken != "" && time.Now().Before(hltbTokenExpiry) {
		return hltbCachedToken, nil
	}

	url := fmt.Sprintf("https://howlongtobeat.com/api/search/site/init?t=%d", time.Now().UnixMilli())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", hltbUserAgent)
	req.Header.Set("Referer", "https://howlongtobeat.com/")
	req.Header.Set("Origin", "https://howlongtobeat.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := hltbClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HLTB init returned status %d", resp.StatusCode)
	}

	var initData hltbInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&initData); err != nil {
		return "", err
	}

	if initData.Token == "" {
		return "", fmt.Errorf("empty token received from HLTB init")
	}

	hltbCachedToken = initData.Token
	hltbTokenExpiry = time.Now().Add(5 * time.Minute)
	return hltbCachedToken, nil
}

func searchHLTB(terms []string, token string) ([]hltbGameItem, error) {
	payload := map[string]interface{}{
		"searchType":  "games",
		"searchTerms": terms,
		"searchPage":  1,
		"size":        5,
		"searchOptions": map[string]interface{}{
			"games": map[string]interface{}{
				"userId":        0,
				"platform":      "",
				"sortCategory":  "popular",
				"rangeCategory": "main",
				"rangeTime":     map[string]interface{}{"min": 0, "max": 0},
				"gameplay":      map[string]interface{}{"perspective": "", "flow": "", "genre": ""},
				"modifier":      "",
			},
		},
		"useCache": true,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://howlongtobeat.com/api/search/site", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", hltbUserAgent)
	req.Header.Set("Referer", "https://howlongtobeat.com/")
	req.Header.Set("Origin", "https://howlongtobeat.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", token)

	resp, err := hltbClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HLTB search returned status %d", resp.StatusCode)
	}

	var searchResp hltbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return searchResp.Data, nil
}

// Search handles GET /api/v1/hltb?q=...
func (h *HLTBHandler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		RespondError(c, http.StatusBadRequest, "MISSING_QUERY", "Search query parameter 'q' is required")
		return
	}

	// Clean punctuation
	cleaned := nonAlphanumericRegex.ReplaceAllString(query, " ")
	rawTerms := strings.Fields(cleaned)
	if len(rawTerms) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"found":   false,
			"message": "Invalid search query",
		})
		return
	}

	token, err := getHLTBToken(false)
	if err != nil {
		token, err = getHLTBToken(true)
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to connect to HowLongToBeat",
		})
		return
	}

	games, err := searchHLTB(rawTerms, token)
	if err != nil {
		// Retry once with refreshed token
		token, err = getHLTBToken(true)
		if err == nil {
			games, _ = searchHLTB(rawTerms, token)
		}
	}

	// If no match and query had multiple words, try leading words
	if len(games) == 0 && len(rawTerms) > 2 && token != "" {
		fallbackTerms := rawTerms[:2]
		games, _ = searchHLTB(fallbackTerms, token)
	}

	if len(games) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"found":   false,
			"message": "No results found on HowLongToBeat",
		})
		return
	}

	top := games[0]
	mainHours := int(math.Round(float64(top.CompMain) / 3600.0))
	extraHours := int(math.Round(float64(top.CompPlus) / 3600.0))
	compHours := int(math.Round(float64(top.Comp100) / 3600.0))

	formatted := "N/A"
	if mainHours > 0 {
		formatted = fmt.Sprintf("%dh", mainHours)
	} else if extraHours > 0 {
		formatted = fmt.Sprintf("%dh", extraHours)
	} else if compHours > 0 {
		formatted = fmt.Sprintf("%dh", compHours)
	}

	response := gin.H{
		"found":       true,
		"gameId":      top.GameID,
		"title":       top.GameName,
		"timeToBeat":  formatted,
		"hltbUrl":     fmt.Sprintf("https://howlongtobeat.com/game/%d", top.GameID),
	}
	if mainHours > 0 {
		response["mainHours"] = mainHours
	}
	if extraHours > 0 {
		response["extraHours"] = extraHours
	}
	if compHours > 0 {
		response["completionistHours"] = compHours
	}

	matches := make([]gin.H, 0, len(games))
	for _, g := range games {
		hrs := int(math.Round(float64(g.CompMain) / 3600.0))
		matches = append(matches, gin.H{
			"id":        g.GameID,
			"name":      g.GameName,
			"mainHours": hrs,
		})
	}
	response["matches"] = matches

	c.JSON(http.StatusOK, response)
}
