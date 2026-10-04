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
      className="group flex items-center justify-between gap-3 px-3.5 py-2.5 rounded-xl hypr-glass border border-[var(--border-color)] hover:border-[var(--accent)] hover:bg-[var(--selection)]/40 transition-all cursor-pointer"
    >
      {/* Left: Title & Subcategory Tag */}
      <div className="flex items-center gap-2 min-w-0">
        {game.subcategory && (
          <span className="px-1.5 py-0.5 rounded text-[10px] font-mono font-bold bg-[var(--selection)] text-[var(--accent)] border border-[var(--border-color)] shrink-0">
            {game.subcategory}
          </span>
        )}

        <span className="font-semibold text-xs sm:text-sm text-[var(--fg-primary)] group-hover:text-[var(--accent)] transition-colors truncate">
          {game.title}
        </span>
      </div>

      {/* Right: Genre, Time to Beat (HLTB) & Delete */}
      <div className="flex items-center gap-2.5 sm:gap-3 text-xs shrink-0">
        {/* Single Genre */}
        <span className="px-2 py-0.5 rounded-md text-[10px] font-mono bg-[var(--selection)]/70 text-[var(--fg-light)] border border-[var(--border-color)]">
          {game.genre}
        </span>

        {/* Time to Beat (HLTB) */}
        {game.timeToBeat && (
          <span 
            className="flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-mono bg-sky-500/10 text-sky-300 border border-sky-500/20"
            title="Time to Beat retrieved from HowLongToBeat"
          >
            <Clock className="w-2.5 h-2.5 opacity-80" />
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
            className="opacity-0 group-hover:opacity-100 p-1 rounded text-[var(--fg-light)] hover:text-rose-400 transition-opacity cursor-pointer"
            title="Remove game"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        )}
      </div>
    </div>
  );
};
