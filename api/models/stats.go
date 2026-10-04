package models

// GameSummary provides brief info about a game in stats
type GameSummary struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Platform string  `json:"platform"`
	Hours    float64 `json:"hours"`
}

// SubcategoryStat aggregates count per subcategory
type SubcategoryStat struct {
	Subcategory string `json:"subcategory"`
	Count       int    `json:"count"`
}

// PlatformStat aggregates metrics per platform
type PlatformStat struct {
	Platform      string            `json:"platform"`
	Count         int               `json:"count"`
	TotalHours    float64           `json:"totalHours"`
	AverageHours  float64           `json:"averageHours"`
	Subcategories []SubcategoryStat `json:"subcategories"`
}

// GenreStat aggregates count per genre
type GenreStat struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

// CatalogStats aggregates overall catalog metrics
type CatalogStats struct {
	TotalGames          int            `json:"totalGames"`
	TotalPlaytimeHours  float64        `json:"totalPlaytimeHours"`
	AveragePlaytimeHours float64       `json:"averagePlaytimeHours"`
	ShortestGame        *GameSummary   `json:"shortestGame,omitempty"`
	LongestGame         *GameSummary   `json:"longestGame,omitempty"`
	PlatformBreakdown   []PlatformStat `json:"platformBreakdown"`
	GenreBreakdown      []GenreStat    `json:"genreBreakdown"`
}
