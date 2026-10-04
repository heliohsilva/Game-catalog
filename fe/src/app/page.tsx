'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { Game, Platform, PLATFORMS, OmarchyTheme, PlatformInfo } from '../types/game';
import { 
  DEFAULT_DARK_THEME_ID, 
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
import { ManagePlatformsModal } from '../components/ManagePlatformsModal';

const STORAGE_KEY_THEME = 'game_catalog_theme_id_v6';

export default function GameCatalogPage() {
  const [games, setGames] = useState<Game[]>([]);
  const [platforms, setPlatforms] = useState<PlatformInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentTheme, setCurrentTheme] = useState<OmarchyTheme>(() => 
    getThemeById(DEFAULT_DARK_THEME_ID)
  );
  const [selectedPlatform, setSelectedPlatform] = useState<Platform | 'All'>('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'title' | 'hours'>('title');

  // Modals state
  const [isThemeModalOpen, setIsThemeModalOpen] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isSearchModalOpen, setIsSearchModalOpen] = useState(false);
  const [isPlatformsModalOpen, setIsPlatformsModalOpen] = useState(false);
  const [inspectedGame, setInspectedGame] = useState<Game | null>(null);
  const [addPlatformPreset, setAddPlatformPreset] = useState<Platform>('PC');

  // Fetch platforms from backend API
  const fetchPlatforms = useCallback(async () => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      const res = await fetch(`${apiUrl}/platforms`);
      if (res.ok) {
        const data = await res.json();
        if (data && Array.isArray(data.platforms)) {
          setPlatforms(data.platforms);
          return data.platforms as PlatformInfo[];
        }
      }
    } catch (err) {
      console.error('Failed to fetch platforms from API:', err);
    }
    return [];
  }, []);

  // Fetch games from backend API
  const fetchGames = useCallback(async () => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      const res = await fetch(`${apiUrl}/games?limit=200`);
      if (res.ok) {
        const data = await res.json();
        if (data && Array.isArray(data.data)) {
          setGames(data.data);
          return;
        }
      }
      setGames([]);
    } catch (err) {
      console.error('Failed to fetch catalog from SQLite API:', err);
      setGames([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Load saved theme and fetch platforms & games on mount
  useEffect(() => {
    try {
      // Clean up legacy localStorage cached games and wallpapers
      localStorage.removeItem('game_catalog_games');
      localStorage.removeItem('game_catalog_games_v5');
      localStorage.removeItem('game_catalog_games_v6');
      localStorage.removeItem('game_catalog_wallpaper_v5');
      localStorage.removeItem('game_catalog_wallpaper_v6');

      const savedThemeId = localStorage.getItem(STORAGE_KEY_THEME);
      if (savedThemeId) {
        const t = getThemeById(savedThemeId);
        setCurrentTheme(t);
        applyThemeToCss(t);
      } else {
        applyThemeToCss(getThemeById(DEFAULT_DARK_THEME_ID));
      }

      fetchPlatforms();
      fetchGames();
    } catch {
      applyThemeToCss(getThemeById(DEFAULT_DARK_THEME_ID));
      setIsLoading(false);
    }
  }, [fetchPlatforms, fetchGames]);

  // Handle platform/subcategory mutations
  const handlePlatformsChanged = async () => {
    const updated = await fetchPlatforms();
    fetchGames();
    if (selectedPlatform !== 'All' && updated.length > 0) {
      if (!updated.some((p) => p.name === selectedPlatform)) {
        setSelectedPlatform('All');
      }
    }
  };

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
        fetchPlatforms(); // Update game counts per platform
        return;
      }
    } catch (err) {
      console.error('Failed to save game to API:', err);
    }
    // Optimistic fallback
    const optimistic: Game = {
      ...newGameData,
      id: 'local-' + Date.now(),
      addedAt: new Date().toISOString(),
    };
    setGames((prev) => [optimistic, ...prev]);
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
      console.error('Failed to update game on API:', err);
    }
    setGames((prev) => prev.map((g) => (g.id === updated.id ? updated : g)));
  };

  // Delete game
  const handleDeleteGame = async (id: string, e?: React.MouseEvent) => {
    if (e) e.stopPropagation();
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    try {
      await fetch(`${apiUrl}/games/${id}`, {
        method: 'DELETE',
      });
      fetchPlatforms(); // Update game counts
    } catch (err) {
      console.error('Failed to delete game on API:', err);
    }
    setGames((prev) => prev.filter((g) => g.id !== id));
    if (inspectedGame && inspectedGame.id === id) {
      setInspectedGame(null);
    }
  };

  // Open add game modal with platform pre-selected
  const handleOpenAddForPlatform = (plat: Platform) => {
    setAddPlatformPreset(plat);
    setIsAddModalOpen(true);
  };

  // Filtered and sorted games
  const filteredGames = useMemo(() => {
    return games
      .filter((game) => {
        if (selectedPlatform !== 'All' && game.platform !== selectedPlatform) {
          return false;
        }
        if (searchQuery.trim()) {
          const q = searchQuery.toLowerCase();
          const matchTitle = game.title.toLowerCase().includes(q);
          const matchGenre = game.genre.toLowerCase().includes(q);
          const matchSub = (game.subcategory || '').toLowerCase().includes(q);
          const matchPlatform = game.platform.toLowerCase().includes(q);
          if (!matchTitle && !matchGenre && !matchSub && !matchPlatform) {
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

  // Dynamic platforms list
  const availablePlatformNames = useMemo(() => {
    return platforms.length > 0 ? platforms.map((p) => p.name) : PLATFORMS;
  }, [platforms]);

  // Platform clusters to display
  const platformsToDisplay: Platform[] = useMemo(() => {
    if (selectedPlatform === 'All') {
      return availablePlatformNames;
    }
    return [selectedPlatform];
  }, [selectedPlatform, availablePlatformNames]);

  return (
    <div className="relative min-h-screen text-[var(--fg-primary)]">
      {/* Static Textured Background (Zero Distractions, Tactile Matte Micro-Texture) */}
      <div 
        className="fixed inset-0 z-0 pointer-events-none select-none bg-texture"
        aria-hidden="true"
      />

      {/* Content Layer (relative z-10) with frosted glass blurring over the static texture */}
      <div className="relative z-10 pb-16">
        {/* Clean Top Menu Bar */}
        <Waybar
          currentTheme={currentTheme}
          onThemeSelect={handleThemeSelect}
          onOpenThemeModal={() => setIsThemeModalOpen(true)}
          onOpenAddModal={() => {
            const firstPlat = availablePlatformNames[0] || 'PC';
            setAddPlatformPreset(firstPlat);
            setIsAddModalOpen(true);
          }}
          onOpenSearchModal={() => setIsSearchModalOpen(true)}
          onOpenPlatformsModal={() => setIsPlatformsModalOpen(true)}
          selectedPlatform={selectedPlatform}
          onSelectPlatform={setSelectedPlatform}
          totalGames={games.length}
          platforms={platforms}
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
            <p className="font-mono text-xs text-[var(--fg-light)] mt-1">Add your first game using the "+ Add" button above, or configure platforms using the "Platforms" menu.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {platformsToDisplay.map((platform) => {
              const gamesForPlatform = filteredGames.filter(
                (g) => g.platform.toLowerCase() === platform.toLowerCase()
              );
              const platInfo = platforms.find(
                (p) => p.name.toLowerCase() === platform.toLowerCase()
              );

              return (
                <PlatformCluster
                  key={platform}
                  platform={platform}
                  subcategories={platInfo?.subcategories}
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
        platforms={platforms}
      />

      {/* Manage Platforms & Subplatforms Modal */}
      <ManagePlatformsModal
        isOpen={isPlatformsModalOpen}
        onClose={() => setIsPlatformsModalOpen(false)}
        platforms={platforms}
        onPlatformsChanged={handlePlatformsChanged}
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
