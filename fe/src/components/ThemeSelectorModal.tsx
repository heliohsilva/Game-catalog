'use client';

import React, { useState } from 'react';
import { X, Search, Check, Palette, Moon, Sun } from 'lucide-react';
import { OmarchyTheme } from '../types/game';
import { OMARCHY_THEMES, DEFAULT_DARK_THEME_ID, DEFAULT_LIGHT_THEME_ID } from '../data/themes';

interface ThemeSelectorModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentTheme: OmarchyTheme;
  onSelectTheme: (themeId: string) => void;
}

export const ThemeSelectorModal: React.FC<ThemeSelectorModalProps> = ({
  isOpen,
  onClose,
  currentTheme,
  onSelectTheme,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [modeFilter, setModeFilter] = useState<'all' | 'dark' | 'light'>('all');

  if (!isOpen) return null;

  const quickDefaults = OMARCHY_THEMES.filter(
    (t) => t.id === DEFAULT_DARK_THEME_ID || t.id === DEFAULT_LIGHT_THEME_ID
  );

  const omarchyOfficialThemes = OMARCHY_THEMES.filter(
    (t) => t.id !== DEFAULT_DARK_THEME_ID && t.id !== DEFAULT_LIGHT_THEME_ID
  ).filter((t) => {
    const matchesSearch =
      t.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.id.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesMode = modeFilter === 'all' || t.mode === modeFilter;
    return matchesSearch && matchesMode;
  });

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs"
      onClick={onClose}
    >
      <div 
        className="hypr-window-in w-full max-w-2xl max-h-[85vh] flex flex-col rounded-2xl hypr-glass border border-[var(--border-color)] shadow-2xl overflow-hidden"
        style={{ background: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--border-color)] bg-[var(--bg-darker)]/60">
          <div className="flex items-center gap-2">
            <Palette className="w-4 h-4 text-[var(--accent)]" />
            <span className="font-bold text-xs text-[var(--fg-primary)] tracking-wide">
              Theme Selector
            </span>
          </div>

          <button
            onClick={onClose}
            className="p-1 rounded-md text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Modal Body */}
        <div className="flex-1 overflow-y-auto p-5 space-y-5">
          
          {/* Quick Defaults */}
          <div>
            <h3 className="text-xs font-mono text-[var(--fg-light)] uppercase mb-2">
              Default Options
            </h3>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              {quickDefaults.map((theme) => {
                const isSelected = currentTheme.id === theme.id;
                return (
                  <button
                    key={theme.id}
                    onClick={() => onSelectTheme(theme.id)}
                    className={`text-left p-3 rounded-xl transition-colors border cursor-pointer ${
                      isSelected
                        ? 'border-[var(--accent)] bg-[var(--selection)]'
                        : 'border-[var(--border-color)] hover:border-[var(--accent)]/60 bg-[var(--bg-lighter)]/30 hover:bg-[var(--selection)]/30'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1.5">
                      <div className="flex items-center gap-2">
                        {theme.mode === 'dark' ? (
                          <Moon className="w-3.5 h-3.5 text-sky-400" />
                        ) : (
                          <Sun className="w-3.5 h-3.5 text-amber-500" />
                        )}
                        <span className="font-bold text-xs text-[var(--fg-primary)]">
                          {theme.name}
                        </span>
                      </div>
                      {isSelected && (
                        <span className="flex items-center gap-0.5 text-[10px] font-bold text-[var(--accent)]">
                          <Check className="w-3 h-3 stroke-[3]" /> Active
                        </span>
                      )}
                    </div>

                    <p className="text-[11px] text-[var(--fg-light)] mb-2">
                      {theme.description}
                    </p>

                    {/* Color Swatches */}
                    <div className="flex items-center gap-1.5">
                      <span
                        className="w-4 h-4 rounded-full border border-white/20"
                        style={{ backgroundColor: theme.background }}
                        title="Background"
                      />
                      <span
                        className="w-4 h-4 rounded-full border border-white/20"
                        style={{ backgroundColor: theme.lighter_background }}
                        title="Card"
                      />
                      <span
                        className="w-4 h-4 rounded-full border border-white/20"
                        style={{ backgroundColor: theme.accent }}
                        title="Accent"
                      />
                      <span
                        className="w-4 h-4 rounded-full border border-white/20"
                        style={{ backgroundColor: theme.foreground }}
                        title="Foreground"
                      />
                      <span className="text-[10px] font-mono text-[var(--fg-light)] ml-auto uppercase">
                        {theme.mode}
                      </span>
                    </div>
                  </button>
                );
              })}
            </div>
          </div>

          <div className="h-px bg-[var(--border-color)]" />

          {/* All Omarchy Themes */}
          <div>
            <div className="flex items-center justify-between gap-3 mb-3">
              <h3 className="text-xs font-mono text-[var(--fg-light)] uppercase">
                All 22 Omarchy Themes
              </h3>

              {/* Mode Filters */}
              <div className="flex items-center gap-1 p-0.5 rounded-lg border border-[var(--border-color)] text-[10px]">
                {(['all', 'dark', 'light'] as const).map((m) => (
                  <button
                    key={m}
                    onClick={() => setModeFilter(m)}
                    className={`px-2 py-0.5 rounded-md uppercase font-semibold transition-colors cursor-pointer ${
                      modeFilter === m
                        ? 'bg-[var(--accent)] text-[var(--bg-primary)]'
                        : 'text-[var(--fg-light)] hover:text-[var(--fg-primary)]'
                    }`}
                  >
                    {m}
                  </button>
                ))}
              </div>
            </div>

            {/* Search Input */}
            <div className="relative mb-3">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-[var(--fg-light)]" />
              <input
                type="text"
                placeholder="Search theme name..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-9 pr-3 py-1.5 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
              />
            </div>

            {/* Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2 max-h-72 overflow-y-auto pr-1">
              {omarchyOfficialThemes.map((theme) => {
                const isSelected = currentTheme.id === theme.id;
                return (
                  <button
                    key={theme.id}
                    onClick={() => onSelectTheme(theme.id)}
                    className={`text-left p-2.5 rounded-xl transition-colors border cursor-pointer ${
                      isSelected
                        ? 'border-[var(--accent)] bg-[var(--selection)]'
                        : 'border-[var(--border-color)] hover:border-[var(--accent)]/50 bg-[var(--bg-lighter)]/20 hover:bg-[var(--selection)]/30'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="font-medium text-xs text-[var(--fg-primary)] truncate max-w-[130px]">
                        {theme.name}
                      </span>
                      {isSelected ? (
                        <Check className="w-3 h-3 text-[var(--accent)] stroke-[3]" />
                      ) : (
                        <span className="text-[9px] uppercase px-1 rounded bg-[var(--muted)]/50 text-[var(--fg-light)]">
                          {theme.mode}
                        </span>
                      )}
                    </div>

                    {/* Palette */}
                    <div className="flex items-center gap-1 mt-1.5">
                      <span
                        className="w-3.5 h-3.5 rounded-full border border-white/10"
                        style={{ backgroundColor: theme.background }}
                      />
                      <span
                        className="w-3.5 h-3.5 rounded-full border border-white/10"
                        style={{ backgroundColor: theme.accent }}
                      />
                      <span
                        className="w-3.5 h-3.5 rounded-full border border-white/10"
                        style={{ backgroundColor: theme.foreground }}
                      />
                    </div>
                  </button>
                );
              })}
            </div>
          </div>

        </div>

        {/* Footer */}
        <div className="flex items-center justify-end px-4 py-2.5 border-t border-[var(--border-color)] bg-[var(--bg-darker)]/60">
          <button
            onClick={onClose}
            className="px-3 py-1 rounded-lg bg-[var(--selection)] hover:bg-[var(--accent)] hover:text-[var(--bg-primary)] transition-colors text-xs font-medium cursor-pointer text-[var(--fg-primary)]"
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
};
