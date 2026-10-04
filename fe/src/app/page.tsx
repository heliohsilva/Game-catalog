'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { Game, Platform, PLATFORMS, OmarchyTheme } from '../types/game';
import { 
  DEFAULT_DARK_THEME_ID, 
  DEFAULT_LIGHT_THEME_ID, 
  getThemeById, 
  applyThemeToCss 
} from '../data/themes';
import { Waybar } from '../components/Waybar';
import { PlatformCluster } from '../components/PlatformCluster';
import { PlatformFilterBar } from '../components/PlatformFilterBar';
import { ThemeSelectorModal } from '../components/ThemeSelectorModal';
import { AddGameModal } from '../components/AddGameModal';
import { GameDetailModal } from '../components/GameDetailModal';
import { SearchOmnibarModal } from '../components/SearchOmnibarModal';

const STORAGE_KEY_THEME = 'game_catalog_theme_id_v6';
const STORAGE_KEY_WALLPAPER = 'game_catalog_wallpaper_v6';

const WALLPAPERS = [
  { id: 'village', name: 'Pixel Village (Night)', url: '/wallpapers/pixel-art-bg.jpg' },
  { id: 'city', name: 'Pixel City Skyline', url: '/wallpapers/pixel-city.gif' },
  { id: 'sunset', name: 'Pixel Sunset City', url: '/wallpapers/sunset-city.gif' },
  { id: 'backyard', name: 'Pixel Backyard', url: '/wallpapers/pixel-backyard.webp' },
];

export default function GameCatalogPage() {
  const [games, setGames] = useState<Game[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentTheme, setCurrentTheme] = useState<OmarchyTheme>(() => 
    getThemeById(DEFAULT_DARK_THEME_ID)
  );
  const [wallpaperIndex, setWallpaperIndex] = useState(0);
  const [selectedPlatform, setSelectedPlatform] = useState<Platform | 'All'>('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'title' | 'hours'>('title');

  // Modals state
  const [isThemeModalOpen, setIsThemeModalOpen] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isSearchModalOpen, setIsSearchModalOpen] = useState(false);
  const [inspectedGame, setInspectedGame] = useState<Game | null>(null);
  const [addPlatformPreset, setAddPlatformPreset] = useState<Platform>('PC');

  // Load saved theme/wallpaper and fetch games from SQLite Backend API on mount
  useEffect(() => {
    try {
      // Clean up any legacy localStorage cached games so browser never shows stale mock data
      localStorage.removeItem('game_catalog_games');
      localStorage.removeItem('game_catalog_games_v5');
      localStorage.removeItem('game_catalog_games_v6');

      const savedThemeId = localStorage.getItem(STORAGE_KEY_THEME);
      if (savedThemeId) {
        const t = getThemeById(savedThemeId);
        setCurrentTheme(t);
        applyThemeToCss(t);
      } else {
        applyThemeToCss(getThemeById(DEFAULT_DARK_THEME_ID));
      }

      const savedWallpaper = localStorage.getItem(STORAGE_KEY_WALLPAPER);
      if (savedWallpaper !== null) {
        const idx = parseInt(savedWallpaper, 10);
        if (!isNaN(idx) && idx >= 0 && idx < WALLPAPERS.length) {
          setWallpaperIndex(idx);
        }
      }

      // Fetch dynamic catalog from SQLite Backend API
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
      fetch(`${apiUrl}/games?limit=200`)
        .then((res) => (res.ok ? res.json() : null))
        .then((res) => {
          if (res && Array.isArray(res.data)) {
            setGames(res.data);
          } else {
            setGames([]);
          }
        })
        .catch((err) => {
          console.error('Failed to fetch catalog from SQLite API:', err);
          setGames([]);
        })
        .finally(() => {
          setIsLoading(false);
        });
    } catch {
      applyThemeToCss(getThemeById(DEFAULT_DARK_THEME_ID));
      setIsLoading(false);
    }
  }, []);

  // Theme change handler
  const handleThemeSelect = (themeId: string) => {
    const theme = getThemeById(themeId);
    setCurrentTheme(theme);
    applyThemeToCss(theme);
    try {
      localStorage.setItem(STORAGE_KEY_THEME, themeId);
    } catch {
      // ignore
    }
  };

  // Cycle Wallpaper handler
  const handleCycleWallpaper = () => {
    setWallpaperIndex((prev) => {
      const next = (prev + 1) % WALLPAPERS.length;
      try {
        localStorage.setItem(STORAGE_KEY_WALLPAPER, String(next));
      } catch {
        // ignore
      }
      return next;
    });
  };

  // Add game
  const handleAddGame = async (newGameData: Omit<Game, 'id' | 'addedAt'>) => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      const res = await fetch(`${apiUrl}/games`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newGameData),
      });
      if (res.ok) {
        const created: Game = await res.json();
        setGames((prev) => [created, ...prev]);
        return;
      }
    } catch (err) {
      console.error('Error adding game to API:', err);
    }

    const fallbackGame: Game = {
      ...newGameData,
      id: `game-${Date.now()}`,
      addedAt: new Date().toISOString(),
    };
    setGames((prev) => [fallbackGame, ...prev]);
  };

  // Update game
  const handleUpdateGame = async (updated: Game) => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      await fetch(`${apiUrl}/games/${updated.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updated),
      });
    } catch (err) {
      console.error('Error updating game in API:', err);
    }
    setGames((prev) => prev.map((g) => (g.id === updated.id ? updated : g)));
  };

  // Delete game
  const handleDeleteGame = async (id: string) => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      await fetch(`${apiUrl}/games/${id}`, {
        method: 'DELETE',
      });
    } catch (err) {
      console.error('Error deleting game from API:', err);
    }
    setGames((prev) => prev.filter((g) => g.id !== id));
    if (inspectedGame?.id === id) {
      setInspectedGame(null);
    }
  };

  // Open add modal targeted for a specific platform
  const handleOpenAddForPlatform = (platform: Platform) => {
    setAddPlatformPreset(platform);
    setIsAddModalOpen(true);
  };

  // Filter & Sort Logic
  const filteredGames = useMemo(() => {
    return games
      .filter((g) => {
        // Platform filter
        if (selectedPlatform !== 'All' && g.platform !== selectedPlatform) {
          return false;
        }

        // Search Query
        if (searchQuery.trim()) {
          const q = searchQuery.toLowerCase();
          const matchTitle = g.title.toLowerCase().includes(q);
          const matchPlatform = g.platform.toLowerCase().includes(q);
          const matchSub = (g.subcategory || '').toLowerCase().includes(q);
          const matchGenre = g.genre.toLowerCase().includes(q);
          if (!matchTitle && !matchPlatform && !matchSub && !matchGenre) {
            return false;
          }
        }

        return true;
      })
      .sort((a, b) => {
        if (sortBy === 'title') {
          return a.title.localeCompare(b.title);
        } else if (sortBy === 'hours') {
          return (b.timeToBeatMain || 0) - (a.timeToBeatMain || 0);
        }
        return 0;
      });
  }, [games, selectedPlatform, searchQuery, sortBy]);

  // Platform clusters to display
  const platformsToDisplay: Platform[] =
    selectedPlatform === 'All' ? PLATFORMS : [selectedPlatform];

  const currentWallpaper = WALLPAPERS[wallpaperIndex];

  return (
    <div className="relative min-h-screen text-[var(--fg-primary)]">
      
      {/* 16-bit Retro Pixel Art Wallpaper in Background (z-0) */}
      <div 
        className="fixed inset-0 z-0 bg-cover bg-center bg-no-repeat pointer-events-none transition-all duration-500"
        style={{
          backgroundImage: `url('${currentWallpaper.url}')`,
          imageRendering: 'pixelated',
        }}
      />

      {/* Subtle Ambient Theme Overlay for transparent components to blur over (z-0) */}
      <div 
        className="fixed inset-0 z-0 pointer-events-none transition-colors duration-200"
        style={{
          backgroundColor: currentTheme.mode === 'dark' 
            ? 'rgba(9, 10, 18, 0.32)' 
            : 'rgba(250, 246, 238, 0.40)',
        }}
      />

      {/* Content Layer (relative z-10) directly on top of wallpaper so backdrop-filter blurs it */}
      <div className="relative z-10 pb-16">
        {/* Clean Top Menu Bar */}
        <Waybar
          currentTheme={currentTheme}
          onThemeSelect={handleThemeSelect}
          onOpenThemeModal={() => setIsThemeModalOpen(true)}
          onOpenAddModal={() => {
            setAddPlatformPreset('PC');
            setIsAddModalOpen(true);
          }}
          onOpenSearchModal={() => setIsSearchModalOpen(true)}
          selectedPlatform={selectedPlatform}
          onSelectPlatform={setSelectedPlatform}
          totalGames={games.length}
          wallpaperName={currentWallpaper.name}
          onCycleWallpaper={handleCycleWallpaper}
        />

        {/* Main Container */}
        <main className="mx-auto w-[96%] max-w-5xl mt-6 space-y-4">
        
        {/* Minimal Header */}
        <div className="flex items-center justify-between pb-2 border-b border-[var(--border-color)]">
          <div>
            <h1 className="text-lg font-bold tracking-tight text-[var(--fg-bright)] drop-shadow-xs">
              Game Catalog
            </h1>
            <p className="text-[11px] text-[var(--fg-light)]">
              Clustered by platform • Time to beat via HowLongToBeat
            </p>
          </div>

          <div className="flex items-center gap-2">
            <span className="px-2.5 py-1 rounded-lg hypr-glass text-[11px] font-mono text-[var(--fg-light)]">
              Theme: <span className="text-[var(--accent)] font-semibold">{currentTheme.name}</span>
            </span>
          </div>
        </div>

        {/* Filter Bar */}
        <PlatformFilterBar
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          sortBy={sortBy}
          onSortChange={setSortBy}
          totalFiltered={filteredGames.length}
        />

        {/* Platform Clusters or Loading / Empty State */}
        {isLoading && games.length === 0 ? (
          <div className="py-16 text-center hypr-glass rounded-2xl p-8 border border-[var(--border-color)]">
            <div className="flex flex-col items-center justify-center gap-3">
              <div className="w-6 h-6 border-2 border-[var(--accent)] border-t-transparent rounded-full animate-spin" />
              <p className="font-mono text-xs text-[var(--fg-light)]">Connecting to catalog database...</p>
            </div>
          </div>
        ) : games.length === 0 ? (
          <div className="py-16 text-center hypr-glass rounded-2xl p-8 border border-[var(--border-color)]">
            <p className="font-bold text-sm text-[var(--fg-primary)]">Catalog database is empty</p>
            <p className="font-mono text-xs text-[var(--fg-light)] mt-1">Add your first game using the "+ Add Game" button above.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {platformsToDisplay.map((platform) => {
              const gamesForPlatform = filteredGames.filter(
                (g) => g.platform === platform
              );

              return (
                <PlatformCluster
                  key={platform}
                  platform={platform}
                  games={gamesForPlatform}
                  onSelectGame={(g) => setInspectedGame(g)}
                  onDeleteGame={handleDeleteGame}
                  onAddGameToPlatform={handleOpenAddForPlatform}
                />
              );
            })}
          </div>
        )}

      </main>
      </div>

      {/* Omarchy Theme Selector Modal */}
      <ThemeSelectorModal
        isOpen={isThemeModalOpen}
        onClose={() => setIsThemeModalOpen(false)}
        currentTheme={currentTheme}
        onSelectTheme={handleThemeSelect}
      />

      {/* Add New Game Modal */}
      <AddGameModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onAddGame={handleAddGame}
        initialPlatform={addPlatformPreset}
      />

      {/* Game Detail Modal */}
      <GameDetailModal
        game={inspectedGame}
        isOpen={!!inspectedGame}
        onClose={() => setInspectedGame(null)}
        onUpdateGame={handleUpdateGame}
        onDeleteGame={handleDeleteGame}
      />

      {/* Fast Search Modal */}
      <SearchOmnibarModal
        isOpen={isSearchModalOpen}
        onClose={() => setIsSearchModalOpen(false)}
        games={games}
        onSelectGame={(g) => setInspectedGame(g)}
      />

    </div>
  );
}
