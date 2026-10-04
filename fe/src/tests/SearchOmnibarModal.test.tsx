import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SearchOmnibarModal } from '../components/SearchOmnibarModal';
import { Game } from '../types/game';

describe('SearchOmnibarModal Component', () => {
  const mockGames: Game[] = [
    { id: '1', title: 'Hollow Knight: Silksong', platform: 'PC', subcategory: 'Steam', genre: 'Metroidvania', timeToBeat: '35h' },
    { id: '2', title: 'Bloodborne', platform: 'PlayStation', subcategory: 'PS4', genre: 'Action RPG', timeToBeat: '34h' },
    { id: '3', title: 'Super Mario World', platform: 'Retro / Emulation', subcategory: 'SNES', genre: 'Platformer', timeToBeat: '5h' },
  ];

  it('renders omnibar input and lists matching games', () => {
    render(
      <SearchOmnibarModal
        isOpen={true}
        onClose={vi.fn()}
        games={mockGames}
        onSelectGame={vi.fn()}
      />
    );

    expect(screen.getByPlaceholderText(/Search game title, platform, store, or genre/i)).toBeInTheDocument();
    expect(screen.getByText('Hollow Knight: Silksong')).toBeInTheDocument();
    expect(screen.getByText('Bloodborne')).toBeInTheDocument();
  });

  it('filters games dynamically as user types', () => {
    render(
      <SearchOmnibarModal
        isOpen={true}
        onClose={vi.fn()}
        games={mockGames}
        onSelectGame={vi.fn()}
      />
    );

    const input = screen.getByPlaceholderText(/Search game title, platform, store, or genre/i);
    fireEvent.change(input, { target: { value: 'Mario' } });

    expect(screen.getByText('Super Mario World')).toBeInTheDocument();
    expect(screen.queryByText('Bloodborne')).not.toBeInTheDocument();
  });

  it('calls onSelectGame and closes when a game item is clicked', () => {
    const handleSelect = vi.fn();
    const handleClose = vi.fn();

    render(
      <SearchOmnibarModal
        isOpen={true}
        onClose={handleClose}
        games={mockGames}
        onSelectGame={handleSelect}
      />
    );

    fireEvent.click(screen.getByText('Bloodborne'));
    expect(handleSelect).toHaveBeenCalledWith(mockGames[1]);
    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('navigates with keyboard and selects game with Enter', () => {
    const handleSelect = vi.fn();
    const handleClose = vi.fn();

    render(
      <SearchOmnibarModal
        isOpen={true}
        onClose={handleClose}
        games={mockGames}
        onSelectGame={handleSelect}
      />
    );

    const input = screen.getByPlaceholderText(/Search game title, platform, store, or genre/i);
    
    // Press ArrowDown to move to item 1 (Bloodborne)
    fireEvent.keyDown(input, { key: 'ArrowDown' });
    // Press Enter to select
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(handleSelect).toHaveBeenCalledWith(mockGames[1]);
    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('closes on Escape key press or backdrop click', () => {
    const handleClose = vi.fn();
    const { container } = render(
      <SearchOmnibarModal
        isOpen={true}
        onClose={handleClose}
        games={mockGames}
        onSelectGame={vi.fn()}
      />
    );

    const input = screen.getByPlaceholderText(/Search game title, platform, store, or genre/i);
    // Escape key
    fireEvent.keyDown(input, { key: 'Escape' });
    expect(handleClose).toHaveBeenCalledTimes(1);

    // Backdrop click
    const backdrop = container.firstChild as HTMLElement;
    fireEvent.click(backdrop);
    expect(handleClose).toHaveBeenCalledTimes(2);
  });
});
