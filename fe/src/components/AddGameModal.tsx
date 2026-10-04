'use client';

import React, { useState, useEffect } from 'react';
import { X, Plus, Loader2 } from 'lucide-react';
import { Game, Platform, PLATFORMS, PLATFORM_SUBCATEGORIES, PlatformInfo } from '../types/game';
import { getApiUrl } from '../utils/api';

interface AddGameModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAddGame: (newGame: Omit<Game, 'id' | 'addedAt'>) => void;
  initialPlatform?: Platform;
  platforms?: PlatformInfo[];
}

export const AddGameModal: React.FC<AddGameModalProps> = ({
  isOpen,
  onClose,
  onAddGame,
  initialPlatform = 'PC',
  platforms,
}) => {
  const [title, setTitle] = useState('');
  const [platform, setPlatform] = useState<Platform>(initialPlatform);
  const [subcategory, setSubcategory] = useState<string>('Steam');
  const [genre, setGenre] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const getSubsForPlatform = (platName: string): string[] => {
    if (platforms && platforms.length > 0) {
      const found = platforms.find((p) => p.name.toLowerCase() === platName.toLowerCase());
      if (found) return found.subcategories || [];
    }
    return PLATFORM_SUBCATEGORIES[platName] || [];
  };

  const availablePlatformNames = platforms && platforms.length > 0 
    ? platforms.map((p) => p.name) 
    : PLATFORMS;

  useEffect(() => {
    const available = platforms && platforms.length > 0 ? platforms.map((p) => p.name) : PLATFORMS;
    const initial = initialPlatform && available.includes(initialPlatform) 
      ? initialPlatform 
      : (available[0] || 'PC');
    setPlatform(initial);
    const subs = getSubsForPlatform(initial);
    setSubcategory(subs.length > 0 ? subs[0] : '');
  }, [initialPlatform, platforms]);

  // Update default subcategory when platform changes
  const handlePlatformChange = (newPlatform: Platform) => {
    setPlatform(newPlatform);
    const subs = getSubsForPlatform(newPlatform);
    setSubcategory(subs.length > 0 ? subs[0] : '');
  };

  if (!isOpen) return null;

  const currentSubcategories = getSubsForPlatform(platform);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    let hltbTime = '';
    let hltbMain: number | undefined = undefined;

    // Automatically retrieve time to beat from HowLongToBeat
    try {
      let res = await fetch(`/api/hltb?q=${encodeURIComponent(title.trim())}`);
      if (!res.ok) {
        const apiUrl = getApiUrl();
        res = await fetch(`${apiUrl}/hltb?q=${encodeURIComponent(title.trim())}`);
      }
      const data = await res.json();
      if (data.found) {
        hltbTime = data.timeToBeat || data.game?.timeToBeat || '';
        hltbMain = data.mainHours ?? data.timeToBeatMain ?? data.game?.timeToBeatMain ?? undefined;
      }
    } catch {
      // Fallback: continue without HLTB data if offline
    }

    const subValue = currentSubcategories.length > 0 ? (subcategory || currentSubcategories[0]) : undefined;

    onAddGame({
      title: title.trim(),
      platform,
      subcategory: subValue,
      genre: genre.trim() || 'Gaming',
      timeToBeat: hltbTime || undefined,
      timeToBeatMain: hltbMain,
    });

    // Reset and close
    setTitle('');
    const defaultSubs = getSubsForPlatform(platform);
    setSubcategory(defaultSubs.length > 0 ? defaultSubs[0] : '');
    setGenre('');
    setIsSubmitting(false);
    onClose();
  };

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md"
      role="dialog"
      aria-modal="true"
      aria-label="Add New Game"
      onClick={onClose}
    >
      <div 
        className="hypr-glass rounded-2xl w-full max-w-md border border-[var(--border-color)] hypr-window-in shadow-2xl overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--border-color)]">
          <div className="flex items-center gap-2">
            <Plus className="w-4 h-4 text-[var(--accent)]" />
            <h2 className="font-bold text-sm text-[var(--fg-bright)] tracking-wide">
              Add Game to Catalog
            </h2>
          </div>
          <button
            onClick={onClose}
            className="p-1 rounded-lg text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
            title="Close"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-5 space-y-4">
          
          {/* Title */}
          <div>
            <label htmlFor="game-title" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
              GAME TITLE *
            </label>
            <input
              id="game-title"
              type="text"
              required
              autoFocus
              placeholder="e.g. Chrono Trigger, Metroid Dread..."
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3 py-2 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
            />
            <p className="text-[10px] text-[var(--fg-light)] mt-1">
              Time to beat is queried automatically from HowLongToBeat.
            </p>
          </div>

          {/* Platform */}
          <div>
            <label htmlFor="game-platform" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
              PLATFORM *
            </label>
            <select
              id="game-platform"
              value={platform}
              onChange={(e) => handlePlatformChange(e.target.value as Platform)}
              className="w-full px-3 py-2 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] bg-[var(--bg-lighter)]"
            >
              {availablePlatformNames.map((p) => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
          </div>

          {/* Subcategory (e.g. Steam / GOG / Epic or Console subplatforms) */}
          {currentSubcategories.length > 0 && (
            <div>
              <label htmlFor="game-subcategory" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
                {platform === 'PC' ? 'STORE / LAUNCHER *' : 'SUBSECTION / LAUNCHER *'}
              </label>
              <select
                id="game-subcategory"
                value={subcategory}
                onChange={(e) => setSubcategory(e.target.value)}
                className="w-full px-3 py-2 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] bg-[var(--bg-lighter)]"
              >
                {currentSubcategories.map((sub) => (
                  <option key={sub} value={sub}>{sub}</option>
                ))}
              </select>
            </div>
          )}

          {/* Genre (Single) */}
          <div>
            <label htmlFor="game-genre" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
              GENRE (SINGLE) *
            </label>
            <input
              id="game-genre"
              type="text"
              required
              placeholder="e.g. RPG, Action, Metroidvania, Platformer..."
              value={genre}
              onChange={(e) => setGenre(e.target.value)}
              className="w-full px-3 py-2 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
            />
          </div>

          {/* Footer Actions */}
          <div className="flex items-center justify-end gap-2 pt-3 border-t border-[var(--border-color)]">
            <button
              type="button"
              onClick={onClose}
              className="px-3.5 py-1.5 rounded-xl border border-[var(--border-color)] text-xs text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-[var(--accent)] text-[var(--bg-primary)] text-xs font-bold transition-opacity hover:opacity-90 disabled:opacity-50 cursor-pointer"
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  <span>Searching HLTB...</span>
                </>
              ) : (
                <span>Add to Catalog</span>
              )}
            </button>
          </div>

        </form>

      </div>
    </div>
  );
};
