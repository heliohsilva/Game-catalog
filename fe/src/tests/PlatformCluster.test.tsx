import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { PlatformCluster } from '../components/PlatformCluster';
import { Game } from '../types/game';

describe('PlatformCluster Component', () => {
  const mockGames: Game[] = [
    {
      id: 'pc-1',
      title: 'Baldur Gate 3',
      platform: 'PC',
      subcategory: 'Steam',
      genre: 'CRPG',
      timeToBeat: '75h',
    },
    {
      id: 'pc-2',
      title: 'Cyberpunk 2077',
      platform: 'PC',
      subcategory: 'GOG',
      genre: 'Action RPG',
      timeToBeat: '25h',
    },
    {
      id: 'pc-3',
      title: 'Alan Wake 2',
      platform: 'PC',
      subcategory: 'Epic',
      genre: 'Survival Horror',
      timeToBeat: '18h',
    },
  ];

  it('renders cluster header with platform title, game count, and add button', () => {
    const handleAdd = vi.fn();
    render(
      <PlatformCluster
        platform="PC"
        games={mockGames}
        onSelectGame={vi.fn()}
        onDeleteGame={vi.fn()}
        onAddGameToPlatform={handleAdd}
      />
    );

    expect(screen.getByRole('heading', { level: 2, name: /PC/i })).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();

    const addBtn = screen.getByRole('button', { name: /Add/i });
    fireEvent.click(addBtn);
    expect(handleAdd).toHaveBeenCalledWith('PC');
  });

  it('renders subsection filter pills for PC (Steam, GOG, Epic)', () => {
    render(
      <PlatformCluster
        platform="PC"
        games={mockGames}
        onSelectGame={vi.fn()}
        onDeleteGame={vi.fn()}
        onAddGameToPlatform={vi.fn()}
      />
    );

    expect(screen.getByRole('button', { name: /ALL \(3\)/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /STEAM \(1\)/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /GOG \(1\)/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /EPIC \(1\)/i })).toBeInTheDocument();
  });

  it('filters games when clicking a subcategory filter pill', () => {
    render(
      <PlatformCluster
        platform="PC"
        games={mockGames}
        onSelectGame={vi.fn()}
        onDeleteGame={vi.fn()}
        onAddGameToPlatform={vi.fn()}
      />
    );

    // Initially all 3 games are visible
    expect(screen.getByText('Baldur Gate 3')).toBeInTheDocument();
    expect(screen.getByText('Cyberpunk 2077')).toBeInTheDocument();
    expect(screen.getByText('Alan Wake 2')).toBeInTheDocument();

    // Click GOG filter pill
    const gogBtn = screen.getByRole('button', { name: /GOG \(1\)/i });
    fireEvent.click(gogBtn);

    // Only GOG game should be visible
    expect(screen.getByText('Cyberpunk 2077')).toBeInTheDocument();
    expect(screen.queryByText('Baldur Gate 3')).not.toBeInTheDocument();
    expect(screen.queryByText('Alan Wake 2')).not.toBeInTheDocument();

    // Click ALL pill to reset
    const allBtn = screen.getByRole('button', { name: /ALL \(3\)/i });
    fireEvent.click(allBtn);
    expect(screen.getByText('Baldur Gate 3')).toBeInTheDocument();
  });

  it('allows collapsing and expanding the platform cluster', () => {
    render(
      <PlatformCluster
        platform="PC"
        games={mockGames}
        onSelectGame={vi.fn()}
        onDeleteGame={vi.fn()}
        onAddGameToPlatform={vi.fn()}
      />
    );

    const toggleBtn = screen.getByTitle('Collapse');
    fireEvent.click(toggleBtn);

    // Subcategories and game list should now be collapsed
    expect(screen.queryByText('Baldur Gate 3')).not.toBeInTheDocument();
    expect(screen.getByTitle('Expand')).toBeInTheDocument();

    // Click expand
    fireEvent.click(screen.getByTitle('Expand'));
    expect(screen.getByText('Baldur Gate 3')).toBeInTheDocument();
  });

  it('renders retro pre-PS1/PS1 subcategories for Retro platform', () => {
    const retroGames: Game[] = [
      { id: 'r1', title: 'Chrono Trigger', platform: 'Retro / Emulation', subcategory: 'SNES', genre: 'JRPG' },
      { id: 'r2', title: 'Metal Gear Solid', platform: 'Retro / Emulation', subcategory: 'PS1', genre: 'Stealth' },
    ];

    render(
      <PlatformCluster
        platform="Retro / Emulation"
        games={retroGames}
        onSelectGame={vi.fn()}
        onDeleteGame={vi.fn()}
        onAddGameToPlatform={vi.fn()}
      />
    );

    expect(screen.getByRole('button', { name: /SNES \(1\)/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /PS1 \(1\)/i })).toBeInTheDocument();
  });
});
