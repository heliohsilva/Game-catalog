'use client';

import React, { useState } from 'react';
import { 
  ChevronDown, 
  ChevronRight, 
  Plus 
} from 'lucide-react';
import { Game, Platform, PLATFORM_SUBCATEGORIES } from '../types/game';
import { GameCard } from './GameCard';

interface PlatformClusterProps {
  platform: Platform;
  games: Game[];
  onSelectGame: (game: Game) => void;
  onDeleteGame: (id: string, e: React.MouseEvent) => void;
  onAddGameToPlatform: (platform: Platform) => void;
}

export const PlatformCluster: React.FC<PlatformClusterProps> = ({
  platform,
  games,
  onSelectGame,
  onDeleteGame,
  onAddGameToPlatform,
}) => {
  const [isExpanded, setIsExpanded] = useState(true);
  const [activeSubcategory, setActiveSubcategory] = useState<string>('All');

  const subcategories = PLATFORM_SUBCATEGORIES[platform] || [];
  const hasSubcategories = subcategories.length > 1;

  // Filter games if subcategory is selected
  const visibleGames = hasSubcategories && activeSubcategory !== 'All'
    ? games.filter((g) => g.subcategory === activeSubcategory)
    : games;

  return (
    <section 
      id={`platform-${platform.toLowerCase().replace(/[^a-z0-9]/g, '-')}`}
      className="hypr-glass rounded-2xl p-4 transition-colors"
    >
      {/* Cluster Header */}
      <div className="flex items-center justify-between pb-3 border-b border-[var(--border-color)]">
        <div className="flex items-center gap-2.5">
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="p-1 rounded-lg text-[var(--fg-light)] hover:text-[var(--fg-primary)] hover:bg-[var(--selection)] transition-colors cursor-pointer"
            title={isExpanded ? 'Collapse' : 'Expand'}
          >
            {isExpanded ? (
              <ChevronDown className="w-4 h-4 text-[var(--accent)]" />
            ) : (
              <ChevronRight className="w-4 h-4" />
            )}
          </button>

          <h2 className="font-bold text-sm text-[var(--fg-primary)] tracking-wide">
            {platform}
          </h2>

          <span className="px-2 py-0.5 rounded-full text-[10px] font-mono bg-[var(--selection)] text-[var(--fg-light)] border border-[var(--border-color)]">
            {games.length}
          </span>
        </div>

        <button
          onClick={() => onAddGameToPlatform(platform)}
          className="flex items-center gap-1 px-2.5 py-1 rounded-lg border border-[var(--border-color)] hover:border-[var(--accent)] hover:bg-[var(--selection)] text-xs text-[var(--fg-light)] hover:text-[var(--fg-primary)] transition-colors cursor-pointer"
        >
          <Plus className="w-3 h-3 text-[var(--accent)]" />
          <span className="text-[11px]">Add</span>
        </button>
      </div>

      {/* Subcategories Filter Pills (e.g. Steam, GOG, Epic / PS2, PS3, PS4, PS5 / Retro Consoles) */}
      {isExpanded && hasSubcategories && (
        <div className="flex items-center gap-1.5 overflow-x-auto py-2.5 border-b border-[var(--border-color)]/60 text-xs">
          <button
            onClick={() => setActiveSubcategory('All')}
            className={`px-2 py-0.5 rounded text-[10px] font-mono transition-colors cursor-pointer ${
              activeSubcategory === 'All'
                ? 'bg-[var(--accent)] text-[var(--bg-primary)] font-bold'
                : 'text-[var(--fg-light)] hover:text-[var(--fg-primary)] bg-[var(--selection)]/40'
            }`}
          >
            ALL ({games.length})
          </button>

          {subcategories.map((sub) => {
            const count = games.filter((g) => g.subcategory === sub).length;
            const isActive = activeSubcategory === sub;
            return (
              <button
                key={sub}
                onClick={() => setActiveSubcategory(sub)}
                className={`px-2 py-0.5 rounded text-[10px] font-mono uppercase transition-colors cursor-pointer whitespace-nowrap ${
                  isActive
                    ? 'bg-[var(--accent)] text-[var(--bg-primary)] font-bold'
                    : 'text-[var(--fg-light)] hover:text-[var(--fg-primary)] bg-[var(--selection)]/30 hover:bg-[var(--selection)]'
                }`}
              >
                {sub} {count > 0 ? `(${count})` : ''}
              </button>
            );
          })}
        </div>
      )}

      {/* Cluster Body: List of games */}
      {isExpanded && (
        <div className="mt-3 space-y-2">
          {visibleGames.length > 0 ? (
            // Group by subcategory if viewing All and platform has multiple subcategories
            hasSubcategories && activeSubcategory === 'All' ? (
              subcategories.map((sub) => {
                const subGames = visibleGames.filter((g) => g.subcategory === sub);
                if (subGames.length === 0) return null;
                return (
                  <div key={sub} className="space-y-1.5 pt-1.5 first:pt-0">
                    <div className="flex items-center gap-2 px-1">
                      <span className="font-mono text-[10px] font-bold tracking-wider text-[var(--accent)] uppercase">
                        {sub}
                      </span>
                      <div className="flex-1 h-px bg-[var(--border-color)]/50" />
                      <span className="text-[10px] font-mono text-[var(--fg-light)] opacity-60">
                        {subGames.length}
                      </span>
                    </div>

                    <div className="space-y-1.5">
                      {subGames.map((game) => (
                        <GameCard
                          key={game.id}
                          game={game}
                          onSelect={onSelectGame}
                          onDeleteGame={onDeleteGame}
                        />
                      ))}
                    </div>
                  </div>
                );
              })
            ) : (
              visibleGames.map((game) => (
                <GameCard
                  key={game.id}
                  game={game}
                  onSelect={onSelectGame}
                  onDeleteGame={onDeleteGame}
                />
              ))
            )
          ) : (
            <div className="py-5 text-center text-xs text-[var(--fg-light)] font-mono opacity-60">
              No games in this cluster.
            </div>
          )}
        </div>
      )}
    </section>
  );
};
