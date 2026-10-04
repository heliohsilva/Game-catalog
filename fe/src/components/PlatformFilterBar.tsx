'use client';

import React from 'react';
import { 
  ArrowUpDown, 
  Search 
} from 'lucide-react';

interface PlatformFilterBarProps {
  searchQuery: string;
  onSearchChange: (q: string) => void;
  sortBy: 'title' | 'hours';
  onSortChange: (sort: 'title' | 'hours') => void;
  totalFiltered: number;
}

export const PlatformFilterBar: React.FC<PlatformFilterBarProps> = ({
  searchQuery,
  onSearchChange,
  sortBy,
  onSortChange,
  totalFiltered,
}) => {
  return (
    <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 text-xs mb-3">
      
      {/* Left: Count */}
      <div className="flex items-center gap-2">
        <span className="text-[11px] font-mono text-[var(--fg-light)]">
          Catalog ({totalFiltered} games)
        </span>
      </div>

      {/* Right: Search & Sort */}
      <div className="flex items-center gap-2">
        <div className="relative flex-1 sm:w-60">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-[var(--fg-light)]" />
          <input
            type="text"
            placeholder="Filter title, store, or genre..."
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-full pl-8 pr-3 py-1.5 rounded-xl hypr-glass border border-[var(--border-color)] focus:border-[var(--accent)] focus:outline-none text-xs text-[var(--fg-primary)] placeholder-[var(--fg-light)]/40"
          />
        </div>

        <div className="flex items-center gap-1 hypr-glass px-2.5 py-1.5 rounded-xl border border-[var(--border-color)] text-xs text-[var(--fg-light)]">
          <ArrowUpDown className="w-3 h-3 text-[var(--accent)]" />
          <select
            value={sortBy}
            onChange={(e) => onSortChange(e.target.value as 'title' | 'hours')}
            className="bg-transparent text-xs text-[var(--fg-primary)] focus:outline-none cursor-pointer"
          >
            <option value="title" className="bg-[var(--bg-primary)]">Title (A-Z)</option>
            <option value="hours" className="bg-[var(--bg-primary)]">Time to Beat</option>
          </select>
        </div>
      </div>

    </div>
  );
};
