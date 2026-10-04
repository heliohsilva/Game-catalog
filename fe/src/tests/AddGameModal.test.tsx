import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { AddGameModal } from '../components/AddGameModal';

describe('AddGameModal Component', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders minimal form with Title, Platform, Subcategory, and Genre', () => {
    render(
      <AddGameModal
        isOpen={true}
        onClose={vi.fn()}
        onAddGame={vi.fn()}
      />
    );

    expect(screen.getByText('Add Game to Catalog')).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/e\.g\. Chrono Trigger, Metroid Dread/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Platform/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Store \/ Launcher/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/e\.g\. RPG, Action, Metroidvania/i)).toBeInTheDocument();
  });

  it('does NOT render removed fields (rating, release year, or manual time to beat)', () => {
    render(
      <AddGameModal
        isOpen={true}
        onClose={vi.fn()}
        onAddGame={vi.fn()}
      />
    );

    expect(screen.queryByLabelText(/Rating/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Release Year/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Launch Year/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Time to Beat/i)).not.toBeInTheDocument();
  });

  it('updates subcategory options when switching platforms', () => {
    render(
      <AddGameModal
        isOpen={true}
        onClose={vi.fn()}
        onAddGame={vi.fn()}
        initialPlatform="PC"
      />
    );

    const platformSelect = screen.getByLabelText(/Platform/i);
    expect(screen.getByText('Steam')).toBeInTheDocument();
    expect(screen.getByText('GOG')).toBeInTheDocument();
    expect(screen.getByText('Epic')).toBeInTheDocument();

    // Switch to Retro
    fireEvent.change(platformSelect, { target: { value: 'Retro / Emulation' } });
    expect(screen.getByText('SNES')).toBeInTheDocument();
    expect(screen.getByText('PS1')).toBeInTheDocument();
    expect(screen.getByText('Genesis')).toBeInTheDocument();
    expect(screen.getByText('Neo Geo')).toBeInTheDocument();
  });

  it('calls onClose when clicking the backdrop outside the modal dialog', () => {
    const handleClose = vi.fn();
    const { container } = render(
      <AddGameModal
        isOpen={true}
        onClose={handleClose}
        onAddGame={vi.fn()}
      />
    );

    // The backdrop is the outer fixed container
    const backdrop = container.firstChild as HTMLElement;
    fireEvent.click(backdrop);
    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('submits form with auto-fetched HLTB time', async () => {
    const handleAdd = vi.fn();
    const handleClose = vi.fn();

    // Mock global fetch for HLTB route
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        found: true,
        title: 'Hades',
        mainHours: 22,
        timeToBeat: '22h',
      }),
    } as Response);

    render(
      <AddGameModal
        isOpen={true}
        onClose={handleClose}
        onAddGame={handleAdd}
        initialPlatform="PC"
      />
    );

    const titleInput = screen.getByPlaceholderText(/e\.g\. Chrono Trigger, Metroid Dread/i);
    const genreInput = screen.getByPlaceholderText(/e\.g\. RPG, Action, Metroidvania/i);

    fireEvent.change(titleInput, { target: { value: 'Hades' } });
    fireEvent.change(genreInput, { target: { value: 'Roguelike' } });

    const submitBtn = screen.getByRole('button', { name: /Add to Catalog/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(handleAdd).toHaveBeenCalledTimes(1);
      expect(handleAdd).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Hades',
          platform: 'PC',
          genre: 'Roguelike',
          timeToBeat: '22h',
          timeToBeatMain: 22,
        })
      );
      expect(handleClose).toHaveBeenCalledTimes(1);
    });
  });
});
