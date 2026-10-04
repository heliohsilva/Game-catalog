'use client';

import React, { useState, useEffect } from 'react';
import { 
  X, 
  Trash2, 
  Search, 
  Loader2 
} from 'lucide-react';
import { Game } from '../types/game';

interface GameDetailModalProps {
  game: Game | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdateGame: (updated: Game) => void;
  onDeleteGame: (id: string) => void;
}

export const GameDetailModal: React.FC<GameDetailModalProps> = ({
  game,
  isOpen,
  onClose,
  onUpdateGame,
  onDeleteGame,
}) => {
  const [editedGame, setEditedGame] = useState<Game | null>(null);
  const [isFetchingHLTB, setIsFetchingHLTB] = useState(false);
  const [hltbStatus, setHltbStatus] = useState('');

  useEffect(() => {
    if (game) {
      setEditedGame({ ...game });
      setHltbStatus('');
    }
  }, [game]);

  if (!isOpen || !game || !editedGame) return null;

  const handleFetchHLTB = async () => {
    setIsFetchingHLTB(true);
    setHltbStatus('Checking HowLongToBeat...');
    try {
      let res = await fetch(`/api/hltb?q=${encodeURIComponent(editedGame.title)}`);
      if (!res.ok) {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
        res = await fetch(`${apiUrl}/hltb?q=${encodeURIComponent(editedGame.title)}`);
      }
      const data = await res.json();
      if (data.found) {
        const formatted = data.timeToBeat || (data.mainHours ? `${data.mainHours}h` : 'N/A');
        const updated = {
          ...editedGame,
          timeToBeat: formatted,
          timeToBeatMain: data.mainHours,
        };
        setEditedGame(updated);
        onUpdateGame(updated);
        setHltbStatus(`Updated: ~${formatted}`);
      } else {
        setHltbStatus('No matches found on HowLongToBeat');
      }
    } catch {
      setHltbStatus('Error reaching HowLongToBeat');
    } finally {
      setIsFetchingHLTB(false);
    }
  };

  const handleDelete = () => {
    if (window.confirm(`Remove "${game.title}" from catalog?`)) {
      onDeleteGame(game.id);
      onClose();
    }
  };

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs"
      onClick={onClose}
    >
      <div 
        className="hypr-window-in w-full max-w-sm rounded-2xl hypr-glass border border-[var(--border-color)] shadow-xl overflow-hidden"
        style={{ background: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--border-color)] bg-[var(--bg-darker)]/60">
          <span className="font-bold text-xs text-[var(--fg-primary)] tracking-wide truncate max-w-[240px]">
            {editedGame.title}
          </span>
          <button
            onClick={onClose}
            className="p-1 rounded-md text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Body */}
        <div className="p-4 space-y-3.5 text-xs">
          
          {/* Title & Platform Header */}
          <div>
            <div className="flex items-center gap-1.5 mb-1.5">
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-[var(--selection)] text-[var(--accent)] border border-[var(--border-color)]">
                {editedGame.platform}
              </span>
              {editedGame.subcategory && (
                <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-[var(--selection)] text-[var(--fg-primary)] border border-[var(--border-color)]">
                  {editedGame.subcategory}
                </span>
              )}
            </div>
            <h2 className="text-lg sm:text-xl font-extrabold text-[var(--fg-bright)] tracking-tight">
              {editedGame.title}
            </h2>
          </div>

          {/* Meta details */}
          <div className="p-3 rounded-xl bg-[var(--bg-lighter)]/40 border border-[var(--border-color)] space-y-2">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[10px] text-[var(--fg-light)] opacity-70">GENRE</span>
              <span className="font-semibold text-xs text-[var(--fg-primary)]">{editedGame.genre}</span>
            </div>

            <div className="h-px bg-[var(--border-color)]/50" />

            {/* HowLongToBeat Time */}
            <div className="flex items-center justify-between">
              <span className="font-mono text-[10px] text-[var(--fg-light)] opacity-70">TIME TO BEAT</span>
              <span className="font-mono font-bold text-xs text-[var(--fg-primary)]">
                {editedGame.timeToBeat ? `~${editedGame.timeToBeat}` : 'Not retrieved'}
              </span>
            </div>
          </div>

          {/* Refresh HLTB button */}
          <div className="flex items-center justify-between">
            <button
              type="button"
              onClick={handleFetchHLTB}
              disabled={isFetchingHLTB}
              className="flex items-center gap-1 text-[11px] text-[var(--accent)] hover:underline disabled:opacity-50 cursor-pointer"
            >
              {isFetchingHLTB ? (
                <Loader2 className="w-3 h-3 animate-spin" />
              ) : (
                <Search className="w-3 h-3" />
              )}
              <span>Refresh from HowLongToBeat</span>
            </button>

            {hltbStatus && (
              <span className="text-[10px] font-mono text-[var(--accent)]">
                {hltbStatus}
              </span>
            )}
          </div>

        </div>

        {/* Footer */}
        <div className="flex items-center justify-between px-4 py-3 border-t border-[var(--border-color)] bg-[var(--bg-darker)]/60 text-xs">
          <button
            onClick={handleDelete}
            className="flex items-center gap-1 px-2.5 py-1 rounded-lg text-rose-400 hover:bg-rose-500/15 transition-colors cursor-pointer"
          >
            <Trash2 className="w-3.5 h-3.5" />
            <span>Remove</span>
          </button>

          <button
            onClick={onClose}
            className="px-3.5 py-1 rounded-lg bg-[var(--selection)] hover:bg-[var(--accent)] hover:text-[var(--bg-primary)] font-medium transition-colors cursor-pointer text-[var(--fg-primary)]"
          >
            Close
          </button>
        </div>

      </div>
    </div>
  );
};
