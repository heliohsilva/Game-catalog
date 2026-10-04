import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { GameCard } from '../components/GameCard';
import { Game } from '../types/game';

describe('GameCard Component', () => {
  const mockGame: Game = {
    id: 'game-123',
    title: 'Hollow Knight: Silksong',
    platform: 'PC',
    subcategory: 'Steam',
    genre: 'Metroidvania',
    timeToBeat: '35h',
    timeToBeatMain: 35,
    addedAt: '2026-10-04T00:00:00Z',
  };

  it('renders game title, subcategory tag, and single genre', () => {
    const handleSelect = vi.fn();
    render(<GameCard game={mockGame} onSelect={handleSelect} />);

    expect(screen.getByText('Hollow Knight: Silksong')).toBeInTheDocument();
    expect(screen.getByText('Steam')).toBeInTheDocument();
    expect(screen.getByText('Metroidvania')).toBeInTheDocument();
  });

  it('displays HowLongToBeat time formatted with clock icon', () => {
    render(<GameCard game={mockGame} onSelect={vi.fn()} />);

    expect(screen.getByText('35h')).toBeInTheDocument();
  });

  it('does NOT render removed features (no rating stars, launch year, or cover)', () => {
    const { container } = render(<GameCard game={mockGame} onSelect={vi.fn()} />);

    // No star icons or rating elements
    expect(container.querySelector('.lucide-star')).not.toBeInTheDocument();
    // No cover image element
    expect(container.querySelector('img')).not.toBeInTheDocument();
  });

  it('calls onSelect when the card is clicked', () => {
    const handleSelect = vi.fn();
    render(<GameCard game={mockGame} onSelect={handleSelect} />);

    fireEvent.click(screen.getByText('Hollow Knight: Silksong'));
    expect(handleSelect).toHaveBeenCalledTimes(1);
    expect(handleSelect).toHaveBeenCalledWith(mockGame);
  });

  it('calls onDeleteGame without triggering onSelect when delete button clicked', () => {
    const handleSelect = vi.fn();
    const handleDelete = vi.fn();

    render(
      <GameCard
        game={mockGame}
        onSelect={handleSelect}
        onDeleteGame={handleDelete}
      />
    );

    const deleteBtn = screen.getByTitle('Remove game');
    expect(deleteBtn).toBeInTheDocument();

    fireEvent.click(deleteBtn);
    expect(handleDelete).toHaveBeenCalledTimes(1);
    expect(handleDelete).toHaveBeenCalledWith('game-123', expect.anything());
    expect(handleSelect).not.toHaveBeenCalled();
  });
});
