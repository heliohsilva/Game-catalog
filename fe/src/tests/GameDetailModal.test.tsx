import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { GameDetailModal } from '../components/GameDetailModal';
import { Game } from '../types/game';

describe('GameDetailModal Component', () => {
  const mockGame: Game = {
    id: 'game-1',
    title: 'Metroid Dread',
    platform: 'Nintendo Switch',
    subcategory: 'Switch',
    genre: 'Metroidvania',
    timeToBeat: '10h',
    timeToBeatMain: 10,
  };

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders game details correctly', () => {
    render(
      <GameDetailModal
        game={mockGame}
        isOpen={true}
        onClose={vi.fn()}
        onUpdateGame={vi.fn()}
        onDeleteGame={vi.fn()}
      />
    );

    // Title appears in header and detail body
    expect(screen.getAllByText('Metroid Dread').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText('Nintendo Switch')).toBeInTheDocument();
    expect(screen.getByText('Metroidvania')).toBeInTheDocument();
    expect(screen.getByText('~10h')).toBeInTheDocument();
  });

  it('refreshes HowLongToBeat time on button click', async () => {
    const handleUpdate = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        found: true,
        timeToBeat: '11h',
        mainHours: 11,
      }),
    } as Response);

    render(
      <GameDetailModal
        game={mockGame}
        isOpen={true}
        onClose={vi.fn()}
        onUpdateGame={handleUpdate}
        onDeleteGame={vi.fn()}
      />
    );

    const refreshBtn = screen.getByRole('button', { name: /Refresh from HowLongToBeat/i });
    fireEvent.click(refreshBtn);

    await waitFor(() => {
      expect(handleUpdate).toHaveBeenCalledWith(
        expect.objectContaining({
          timeToBeat: '11h',
          timeToBeatMain: 11,
        })
      );
      expect(screen.getByText(/Updated: ~11h/i)).toBeInTheDocument();
    });
  });

  it('deletes game when trash button clicked and confirmed', () => {
    const handleDelete = vi.fn();
    const handleClose = vi.fn();
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    render(
      <GameDetailModal
        game={mockGame}
        isOpen={true}
        onClose={handleClose}
        onUpdateGame={vi.fn()}
        onDeleteGame={handleDelete}
      />
    );

    const deleteBtn = screen.getByRole('button', { name: /^Remove$/i });
    fireEvent.click(deleteBtn);

    expect(handleDelete).toHaveBeenCalledWith('game-1');
    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('closes on backdrop click', () => {
    const handleClose = vi.fn();
    const { container } = render(
      <GameDetailModal
        game={mockGame}
        isOpen={true}
        onClose={handleClose}
        onUpdateGame={vi.fn()}
        onDeleteGame={vi.fn()}
      />
    );

    const backdrop = container.firstChild as HTMLElement;
    fireEvent.click(backdrop);
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
