'use client';

import React, { useState, useEffect } from 'react';
import { X, Plus, Loader2 } from 'lucide-react';
import { Game, Platform, PLATFORMS, PLATFORM_SUBCATEGORIES } from '../types/game';

interface AddGameModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAddGame: (newGame: Omit<Game, 'id' | 'addedAt'>) => void;
  initialPlatform?: Platform;
}

export const AddGameModal: React.FC<AddGameModalProps> = ({
  isOpen,
  onClose,
  onAddGame,
  initialPlatform = 'PC',
}) => {
  const [title, setTitle] = useState('');
  const [platform, setPlatform] = useState<Platform>(initialPlatform);
  const [subcategory, setSubcategory] = useState<string>('Steam');
  const [genre, setGenre] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (initialPlatform) {
      setPlatform(initialPlatform);
      const subs = PLATFORM_SUBCATEGORIES[initialPlatform] || [];
      if (subs.length > 0) {
        setSubcategory(subs[0]);
      }
    }
  }, [initialPlatform]);

  // Update default subcategory when platform changes
  const handlePlatformChange = (newPlatform: Platform) => {
    setPlatform(newPlatform);
    const subs = PLATFORM_SUBCATEGORIES[newPlatform] || [];
    if (subs.length > 0) {
      setSubcategory(subs[0]);
    }
  };

  if (!isOpen) return null;

  const currentSubcategories = PLATFORM_SUBCATEGORIES[platform] || [];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    let hltbTime = '';
    let hltbMain: number | undefined = undefined;

    // Automatically retrieve time to beat from HowLongToBeat
    try {
      const res = await fetch(`/api/hltb?q=${encodeURIComponent(title.trim())}`);
      if (res.ok) {
        const data = await res.json();
        if (data.found) {
          hltbTime = data.timeToBeat || `${data.mainHours}h`;
          hltbMain = data.mainHours;
        }
      }
    } catch {
      // Graceful fallback if offline
    }

    onAddGame({
      title: title.trim(),
      platform,
      subcategory: currentSubcategories.length > 1 ? subcategory : undefined,
      genre: genre.trim() || 'Gaming',
      timeToBeat: hltbTime || undefined,
      timeToBeatMain: hltbMain,
    });

    setIsSubmitting(false);
    setTitle('');
    setGenre('');
    onClose();
  };

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs"
      onClick={onClose}
    >
      <div 
        className="hypr-window-in w-full max-w-md rounded-2xl hypr-glass border border-[var(--border-color)] shadow-2xl overflow-hidden"
        style={{ background: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Title Bar */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--border-color)] bg-[var(--bg-darker)]/60">
          <span className="font-bold text-xs text-[var(--fg-primary)] tracking-wide">
            Add Game to Catalog
          </span>
          <button
            onClick={onClose}
            className="p-1 rounded-md text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Form Body: Only Title, Platform, Subcategory, Genre */}
        <form onSubmit={handleSubmit} className="p-5 space-y-3.5 text-xs">
          
          {/* Title */}
          <div>
            <label htmlFor="game-title" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
              GAME TITLE *
            </label>
            <input
              id="game-title"
              type="text"
              required
              placeholder="e.g. Chrono Trigger, Metroid Dread..."
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3 py-2 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
            />
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
              {PLATFORMS.map((p) => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
          </div>

          {/* Subcategory (e.g. Steam / GOG / Epic or PS2 / PS3 / PS4 / PS5 or Retro Consoles) */}
          {currentSubcategories.length > 1 && (
            <div>
              <label htmlFor="game-subcategory" className="block text-[11px] font-mono text-[var(--fg-light)] mb-1">
                {platform === 'PC' ? 'STORE / LAUNCHER *' : 'SUBSECTION / CONSOLE *'}
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

          {/* Form Actions */}
          <div className="pt-3 border-t border-[var(--border-color)] flex items-center justify-end gap-2">
            <button
              type="button"
              onClick={onClose}
              className="px-3.5 py-1.5 rounded-xl border border-[var(--border-color)] text-xs text-[var(--fg-light)] hover:text-[var(--fg-primary)] transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-[var(--accent)] text-[var(--bg-primary)] font-bold text-xs hover:opacity-90 disabled:opacity-60 transition-opacity cursor-pointer"
            >
              {isSubmitting && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
              <span>Add to Catalog</span>
            </button>
          </div>

        </form>
      </div>
    </div>
  );
};
