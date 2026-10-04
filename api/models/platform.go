package models

import "strings"

// Platform constants
const (
	PlatformPC             = "PC"
	PlatformPlayStation    = "PlayStation"
	PlatformNintendoSwitch = "Nintendo Switch"
	PlatformXbox           = "Xbox"
	PlatformRetroEmulation = "Retro / Emulation"
)

// ValidPlatforms is the canonical list of supported platforms
var ValidPlatforms = []string{
	PlatformPC,
	PlatformPlayStation,
	PlatformNintendoSwitch,
	PlatformXbox,
	PlatformRetroEmulation,
}

// ValidSubcategories maps each platform to its canonical allowed subcategories
var ValidSubcategories = map[string][]string{
	PlatformPC:             {"Steam", "GOG", "Epic"},
	PlatformPlayStation:    {"PS2", "PS3", "PS4", "PS5"},
	PlatformNintendoSwitch: {"Switch"},
	PlatformXbox:           {"Xbox Series X/S", "Xbox One", "Xbox 360"},
	PlatformRetroEmulation: {"PS1", "N64", "SNES", "Genesis", "NES", "Master System", "Neo Geo"},
}

// PlatformInfo contains metadata about a platform and its game count
type PlatformInfo struct {
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories"`
	GameCount     int      `json:"gameCount"`
}

// NormalizePlatform normalizes and matches a platform string case-insensitively
func NormalizePlatform(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	for _, p := range ValidPlatforms {
		if strings.EqualFold(trimmed, p) {
			return p, true
		}
	}
	return "", false
}

// NormalizeSubcategory normalizes a subcategory for a given platform
func NormalizeSubcategory(platform, subcategory string) (string, bool) {
	trimmed := strings.TrimSpace(subcategory)
	if trimmed == "" {
		return "", false
	}

	allowed, ok := ValidSubcategories[platform]
	if !ok {
		return "", false
	}

	for _, s := range allowed {
		if strings.EqualFold(trimmed, s) {
			return s, true
		}
	}

	return "", false
}
