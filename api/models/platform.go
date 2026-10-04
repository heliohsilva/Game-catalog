package models

import (
	"strings"
	"sync"
)

// Platform constants
const (
	PlatformPC             = "PC"
	PlatformPlayStation    = "PlayStation"
	PlatformNintendoSwitch = "Nintendo Switch"
	PlatformXbox           = "Xbox"
	PlatformRetroEmulation = "Retro / Emulation"
)

var (
	platformMu sync.RWMutex

	// ValidPlatforms is the canonical list of supported platforms
	ValidPlatforms = []string{
		PlatformPC,
		PlatformPlayStation,
		PlatformNintendoSwitch,
		PlatformXbox,
		PlatformRetroEmulation,
	}

	// ValidSubcategories maps each platform to its canonical allowed subcategories
	ValidSubcategories = map[string][]string{
		PlatformPC:             {"Steam", "GOG", "Epic"},
		PlatformPlayStation:    {"PS2", "PS3", "PS4", "PS5"},
		PlatformNintendoSwitch: {"Switch"},
		PlatformXbox:           {"Xbox Series X/S", "Xbox One", "Xbox 360"},
		PlatformRetroEmulation: {"PS1", "N64", "SNES", "Genesis", "NES", "Master System", "Neo Geo"},
	}
)

// PlatformInfo contains metadata about a platform and its game count
type PlatformInfo struct {
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories"`
	GameCount     int      `json:"gameCount"`
}

// CreatePlatformRequest defines the payload for creating a new platform
type CreatePlatformRequest struct {
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories,omitempty"`
}

// CreateSubcategoryRequest defines the payload for adding a subplatform
type CreateSubcategoryRequest struct {
	Name        string `json:"name"`
	Subcategory string `json:"subcategory,omitempty"`
}

// SetDynamicPlatforms updates the global ValidPlatforms and ValidSubcategories
func SetDynamicPlatforms(platforms []PlatformInfo) {
	platformMu.Lock()
	defer platformMu.Unlock()

	newPlatforms := make([]string, 0, len(platforms))
	newSubcats := make(map[string][]string, len(platforms))

	for _, p := range platforms {
		name := strings.TrimSpace(p.Name)
		if name == "" {
			continue
		}
		newPlatforms = append(newPlatforms, name)
		if p.Subcategories == nil {
			newSubcats[name] = []string{}
		} else {
			subs := make([]string, 0, len(p.Subcategories))
			for _, s := range p.Subcategories {
				st := strings.TrimSpace(s)
				if st != "" {
					subs = append(subs, st)
				}
			}
			newSubcats[name] = subs
		}
	}

	ValidPlatforms = newPlatforms
	ValidSubcategories = newSubcats
}

// AddPlatformRegistry updates in-memory registry with a new platform
func AddPlatformRegistry(name string, subcategories []string) {
	platformMu.Lock()
	defer platformMu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return
	}

	for _, p := range ValidPlatforms {
		if strings.EqualFold(p, name) {
			// Update subcategories if already exists
			for k := range ValidSubcategories {
				if strings.EqualFold(k, name) {
					ValidSubcategories[k] = subcategories
					return
				}
			}
			return
		}
	}

	ValidPlatforms = append(ValidPlatforms, name)
	if subcategories == nil {
		ValidSubcategories[name] = []string{}
	} else {
		ValidSubcategories[name] = subcategories
	}
}

// RemovePlatformRegistry removes a platform from the in-memory registry
func RemovePlatformRegistry(name string) {
	platformMu.Lock()
	defer platformMu.Unlock()

	name = strings.TrimSpace(name)
	filtered := make([]string, 0, len(ValidPlatforms))
	for _, p := range ValidPlatforms {
		if !strings.EqualFold(p, name) {
			filtered = append(filtered, p)
		}
	}
	ValidPlatforms = filtered

	for k := range ValidSubcategories {
		if strings.EqualFold(k, name) {
			delete(ValidSubcategories, k)
		}
	}
}

// AddSubcategoryRegistry adds a subcategory to a platform in-memory
func AddSubcategoryRegistry(platform, subcategory string) {
	platformMu.Lock()
	defer platformMu.Unlock()

	platform = strings.TrimSpace(platform)
	subcategory = strings.TrimSpace(subcategory)
	if platform == "" || subcategory == "" {
		return
	}

	for pName, subs := range ValidSubcategories {
		if strings.EqualFold(pName, platform) {
			for _, s := range subs {
				if strings.EqualFold(s, subcategory) {
					return
				}
			}
			ValidSubcategories[pName] = append(subs, subcategory)
			return
		}
	}
}

// RemoveSubcategoryRegistry removes a subcategory from a platform in-memory
func RemoveSubcategoryRegistry(platform, subcategory string) {
	platformMu.Lock()
	defer platformMu.Unlock()

	platform = strings.TrimSpace(platform)
	subcategory = strings.TrimSpace(subcategory)

	for pName, subs := range ValidSubcategories {
		if strings.EqualFold(pName, platform) {
			filtered := make([]string, 0, len(subs))
			for _, s := range subs {
				if !strings.EqualFold(s, subcategory) {
					filtered = append(filtered, s)
				}
			}
			ValidSubcategories[pName] = filtered
			return
		}
	}
}

// NormalizePlatform normalizes and matches a platform string case-insensitively
func NormalizePlatform(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	platformMu.RLock()
	defer platformMu.RUnlock()

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

	platformMu.RLock()
	defer platformMu.RUnlock()

	for pName, allowed := range ValidSubcategories {
		if strings.EqualFold(pName, platform) {
			// If platform has no subcategories defined, any trimmed subcategory is valid
			if len(allowed) == 0 {
				return trimmed, true
			}
			for _, s := range allowed {
				if strings.EqualFold(trimmed, s) {
					return s, true
				}
			}
			return "", false
		}
	}

	return "", false
}

// GetValidPlatformsCopy returns a thread-safe copy of valid platforms
func GetValidPlatformsCopy() []string {
	platformMu.RLock()
	defer platformMu.RUnlock()

	res := make([]string, len(ValidPlatforms))
	copy(res, ValidPlatforms)
	return res
}

// GetValidSubcategoriesCopy returns a thread-safe copy of subcategories for a platform
func GetValidSubcategoriesCopy(platform string) ([]string, bool) {
	platformMu.RLock()
	defer platformMu.RUnlock()

	for pName, subs := range ValidSubcategories {
		if strings.EqualFold(pName, platform) {
			res := make([]string, len(subs))
			copy(res, subs)
			return res, true
		}
	}
	return nil, false
}
