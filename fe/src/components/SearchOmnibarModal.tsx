'use client';

import React, { useState, useEffect, useRef } from 'react';
import { Search, CornerDownLeft, Clock } from 'lucide-react';
import { Game } from '../types/game';

interface SearchOmnibarModalProps {
  isOpen: boolean;
  onClose: () => void;
  games: Game[];
  onSelectGame: (game: Game) => void;
}

export const SearchOmnibarModal: React.FC<SearchOmnibarModalProps> = ({
  isOpen,
  onClose,
  games,
  onSelectGame,
}) => {
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 40);
      setSelectedIndex(0);
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const filtered = games.filter((g) => {
    const q = query.toLowerCase();
    return (
      g.title.toLowerCase().includes(q) ||
      g.platform.toLowerCase().includes(q) ||
      (g.subcategory && g.subcategory.toLowerCase().includes(q)) ||
      g.genre.toLowerCase().includes(q)
    );
  });

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSelectedIndex((prev) => Math.min(prev + 1, filtered.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSelectedIndex((prev) => Math.max(prev - 1, 0));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (filtered[selectedIndex]) {
        onSelectGame(filtered[selectedIndex]);
        onClose();
      }
    } else if (e.key === 'Escape') {
      onClose();
    }
  };

  return (
    <div 
      className="fixed inset-0 z-50 flex items-start justify-center pt-24 px-4 bg-black/50 backdrop-blur-xs"
      onClick={onClose}
    >
      <div 
        className="hypr-window-in w-full max-w-lg rounded-2xl hypr-glass border border-[var(--border-color)] shadow-2xl overflow-hidden"
        style={{ background: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Search Input Bar */}
        <div className="flex items-center gap-3 px-4 py-3 border-b border-[var(--border-color)] bg-[var(--bg-darker)]/60">
          <Search className="w-4 h-4 text-[var(--accent)]" />
          <input
            ref={inputRef}
            type="text"
            placeholder="Search game title, platform, store, or genre..."
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              setSelectedIndex(0);
            }}
            onKeyDown={handleKeyDown}
            className="flex-1 bg-transparent text-xs text-[var(--fg-primary)] focus:outline-none placeholder-[var(--fg-light)]/40 font-mono"
          />
          <kbd className="px-1.5 py-0.5 rounded text-[10px] font-mono bg-[var(--selection)] text-[var(--fg-light)]">
            ESC
          </kbd>
        </div>

        {/* Results List */}
        <div className="max-h-72 overflow-y-auto p-2 space-y-1">
          {filtered.length > 0 ? (
            filtered.map((game, idx) => {
              const isSelected = idx === selectedIndex;
              return (
                <div
                  key={game.id}
                  onClick={() => {
                    onSelectGame(game);
                    onClose();
                  }}
                  onMouseEnter={() => setSelectedIndex(idx)}
                  className={`flex items-center justify-between p-2.5 rounded-xl cursor-pointer transition-colors ${
                    isSelected
                      ? 'bg-[var(--accent)] text-[var(--bg-primary)] font-semibold'
                      : 'hover:bg-[var(--selection)] text-[var(--fg-primary)]'
                  }`}
                >
                  <div className="truncate min-w-0">
                    <div className="text-xs truncate font-bold">{game.title}</div>
                    <div className={`text-[10px] font-mono ${isSelected ? 'text-[var(--bg-primary)]/80' : 'text-[var(--fg-light)]'}`}>
                      {game.platform} {game.subcategory ? `[${game.subcategory}]` : ''} • {game.genre}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    {game.timeToBeat && (
                      <span className={`flex items-center gap-1 text-[10px] font-mono px-2 py-0.5 rounded ${isSelected ? 'bg-black/20' : 'bg-sky-500/10 text-sky-300'}`}>
                        <Clock className="w-2.5 h-2.5" />
                        <span>{game.timeToBeat}</span>
                      </span>
                    )}
                    {isSelected && <CornerDownLeft className="w-3.5 h-3.5" />}
                  </div>
                </div>
              );
            })
          ) : (
            <div className="p-6 text-center text-xs text-[var(--fg-light)] font-mono opacity-60">
              No games found for "{query}".
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
