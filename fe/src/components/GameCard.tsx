'use client';

import React from 'react';
import { Clock, Trash2 } from 'lucide-react';
import { Game } from '../types/game';

interface GameCardProps {
  game: Game;
  onSelect: (game: Game) => void;
  onDeleteGame?: (id: string, e: React.MouseEvent) => void;
}

export const GameCard: React.FC<GameCardProps> = ({
  game,
  onSelect,
  onDeleteGame,
}) => {
  return (
    <div
      onClick={() => onSelect(game)}
      className="group flex items-center justify-between gap-3 px-4 py-3 sm:py-3.5 rounded-xl hypr-glass border border-[var(--border-color)] hover:border-[var(--accent)] hover:bg-[var(--selection)]/40 hover:shadow-lg transition-all cursor-pointer"
    >
      {/* Left: Prominent Game Title & Subcategory Tag */}
      <div className="flex items-center gap-2.5 min-w-0">
        <span className="font-bold text-base sm:text-lg text-[var(--fg-bright)] group-hover:text-[var(--accent)] transition-colors truncate tracking-tight">
          {game.title}
        </span>

        {game.subcategory && (
          <span className="px-2 py-0.5 rounded text-[10px] sm:text-[11px] font-mono font-semibold bg-[var(--selection)] text-[var(--accent)] border border-[var(--border-color)] shrink-0">
            {game.subcategory}
          </span>
        )}
      </div>

      {/* Right: Single Genre, Time to Beat (HLTB) & Delete */}
      <div className="flex items-center gap-2.5 sm:gap-3 text-xs shrink-0">
        {/* Single Genre */}
        <span className="px-2.5 py-1 rounded-md text-[11px] font-mono bg-[var(--selection)]/70 text-[var(--fg-light)] border border-[var(--border-color)]">
          {game.genre}
        </span>

        {/* Time to Beat (HLTB) */}
        {game.timeToBeat && (
          <span 
            className="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-mono font-medium bg-sky-500/10 text-sky-300 border border-sky-500/20"
            title="Time to Beat retrieved from HowLongToBeat"
          >
            <Clock className="w-3 h-3 opacity-80" />
            <span>{game.timeToBeat}</span>
          </span>
        )}

        {/* Delete */}
        {onDeleteGame && (
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              onDeleteGame(game.id, e);
            }}
            className="opacity-0 group-hover:opacity-100 p-1.5 rounded-lg text-[var(--fg-light)] hover:text-rose-400 hover:bg-rose-500/10 transition-all cursor-pointer"
            title="Remove game"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        )}
      </div>
    </div>
  );
};
