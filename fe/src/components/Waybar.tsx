'use client';

import React from 'react';
import { 
  Gamepad2, 
  Search, 
  Plus, 
  Palette, 
  Sun, 
  Moon
} from 'lucide-react';
import { OmarchyTheme, Platform, PLATFORMS } from '../types/game';
import { DEFAULT_DARK_THEME_ID, DEFAULT_LIGHT_THEME_ID } from '../data/themes';

interface WaybarProps {
  currentTheme: OmarchyTheme;
  onThemeSelect: (themeId: string) => void;
  onOpenThemeModal: () => void;
  onOpenAddModal: () => void;
  onOpenSearchModal: () => void;
  selectedPlatform: Platform | 'All';
  onSelectPlatform: (platform: Platform | 'All') => void;
  totalGames: number;
}

export const Waybar: React.FC<WaybarProps> = ({
  currentTheme,
  onThemeSelect,
  onOpenThemeModal,
  onOpenAddModal,
  onOpenSearchModal,
  selectedPlatform,
  onSelectPlatform,
  totalGames,
}) => {
  const handleQuickThemeToggle = () => {
    if (currentTheme.id === DEFAULT_DARK_THEME_ID) {
      onThemeSelect(DEFAULT_LIGHT_THEME_ID);
    } else {
      onThemeSelect(DEFAULT_DARK_THEME_ID);
    }
  };

  const navPlatforms: (Platform | 'All')[] = ['All', ...PLATFORMS];

  return (
    <header className="sticky top-3 z-40 mx-auto w-[96%] max-w-5xl">
      <nav 
        aria-label="Navigation Bar" 
        className="hypr-panel rounded-2xl px-3.5 py-2 transition-colors shadow-lg"
      >
        <div className="flex items-center justify-between gap-3">
          
          {/* Left: Title & Platform Switcher */}
          <div className="flex items-center gap-2.5">
            <div className="flex items-center gap-2 text-sm font-bold text-[var(--fg-primary)]">
              <Gamepad2 className="w-4 h-4 text-[var(--accent)]" />
              <span className="font-retro text-xs tracking-wider">CATALOG</span>
            </div>

            {/* Platform Tabs */}
            <div className="hidden sm:flex items-center gap-1 pl-2.5 border-l border-[var(--border-color)]">
              {navPlatforms.map((p) => {
                const isActive = selectedPlatform === p;
                const label =
                  p === 'All' ? 'All' :
                  p === 'Nintendo Switch' ? 'Switch' :
                  p === 'PlayStation' ? 'PlayStation' :
                  p === 'Xbox' ? 'Xbox' :
                  p === 'Retro / Emulation' ? 'Retro' : 'PC';

                return (
                  <button
                    key={p}
                    onClick={() => onSelectPlatform(p)}
                    className={`px-2.5 py-1 text-xs rounded-lg transition-colors cursor-pointer ${
                      isActive
                        ? 'bg-[var(--accent)] text-[var(--bg-primary)] font-bold'
                        : 'text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)]'
                    }`}
                  >
                    {label}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Right: Search, Add Game, Wallpaper Cycle & Theme Controls */}
          <div className="flex items-center gap-1.5 text-xs">
            
            {/* Search */}
            <button
              onClick={onOpenSearchModal}
              className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border border-[var(--border-color)] hover:border-[var(--accent)] text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
              title="Search catalog"
            >
              <Search className="w-3.5 h-3.5 text-[var(--accent)]" />
              <span className="hidden md:inline">Search</span>
            </button>

            {/* Add Game */}
            <button
              onClick={onOpenAddModal}
              className="flex items-center gap-1 px-3 py-1.5 rounded-xl bg-[var(--accent)] text-[var(--bg-primary)] font-bold transition-opacity hover:opacity-90 active:scale-98 cursor-pointer"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>Add</span>
            </button>

            <div className="h-4 w-px bg-[var(--border-color)] mx-0.5" />

            {/* Quick 1-click Theme Toggle (Cyber Blue / Cream Latte) */}
            <button
              onClick={handleQuickThemeToggle}
              className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border border-[var(--border-color)] hover:border-[var(--accent)] text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
              title={`Switch to ${currentTheme.id === DEFAULT_DARK_THEME_ID ? 'Cream Latte (Default Light)' : 'Cyber Blue (Default Dark)'}`}
            >
              {currentTheme.mode === 'dark' ? (
                <>
                  <Moon className="w-3.5 h-3.5 text-[var(--accent)]" />
                  <span className="hidden lg:inline text-[11px]">Dark</span>
                </>
              ) : (
                <>
                  <Sun className="w-3.5 h-3.5 text-amber-500" />
                  <span className="hidden lg:inline text-[11px]">Light</span>
                </>
              )}
            </button>

            {/* Omarchy Themes List Modal */}
            <button
              onClick={onOpenThemeModal}
              className="p-1.5 rounded-xl border border-[var(--border-color)] hover:border-[var(--accent)] text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
              title="Browse all 22 Omarchy Themes"
            >
              <Palette className="w-3.5 h-3.5 text-[var(--accent)]" />
            </button>

          </div>

        </div>
      </nav>
    </header>
  );
};
