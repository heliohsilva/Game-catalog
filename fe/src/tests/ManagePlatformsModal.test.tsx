import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ManagePlatformsModal } from '../components/ManagePlatformsModal';
import { PlatformInfo } from '../types/game';

describe('ManagePlatformsModal Component', () => {
  const mockPlatforms: PlatformInfo[] = [
    {
      name: 'PC',
      subcategories: ['Steam', 'GOG', 'Epic'],
      gameCount: 12,
    },
    {
      name: 'Nintendo Switch',
      subcategories: ['Switch'],
      gameCount: 5,
    },
  ];

  beforeEach(() => {
    vi.restoreAllMocks();
    window.confirm = vi.fn().mockReturnValue(true);
  });

  it('renders modal header, add platform form, and platform list', () => {
    render(
      <ManagePlatformsModal
        isOpen={true}
        onClose={vi.fn()}
        platforms={mockPlatforms}
        onPlatformsChanged={vi.fn()}
      />
    );

    expect(screen.getByText('Manage Platforms & Subplatforms')).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/e\.g\. Nintendo DS/i)).toBeInTheDocument();
    expect(screen.getByText('PC')).toBeInTheDocument();
    expect(screen.getByText('12 games')).toBeInTheDocument();
    expect(screen.getByText('Steam')).toBeInTheDocument();
    expect(screen.getByText('GOG')).toBeInTheDocument();
    expect(screen.getByText('Epic')).toBeInTheDocument();
    expect(screen.getByText('Nintendo Switch')).toBeInTheDocument();
    expect(screen.getByText('5 games')).toBeInTheDocument();
  });

  it('submits new platform form and triggers onPlatformsChanged', async () => {
    const handlePlatformsChanged = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        platform: { name: 'Nintendo DS', subcategories: ['Cartridge'] },
      }),
    } as Response);

    render(
      <ManagePlatformsModal
        isOpen={true}
        onClose={vi.fn()}
        platforms={mockPlatforms}
        onPlatformsChanged={handlePlatformsChanged}
      />
    );

    const nameInput = screen.getByPlaceholderText(/e\.g\. Nintendo DS/i);
    const subsInput = screen.getByPlaceholderText(/e\.g\. Cartridge, Homebrew/i);
    const submitBtn = screen.getByRole('button', { name: /Add Platform/i });

    fireEvent.change(nameInput, { target: { value: 'Nintendo DS' } });
    fireEvent.change(subsInput, { target: { value: 'Cartridge, Homebrew' } });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/platforms'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            name: 'Nintendo DS',
            subcategories: ['Cartridge', 'Homebrew'],
          }),
        })
      );
      expect(handlePlatformsChanged).toHaveBeenCalledTimes(1);
    });
  });

  it('deletes a platform after confirmation', async () => {
    const handlePlatformsChanged = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ message: 'Platform deleted' }),
    } as Response);

    render(
      <ManagePlatformsModal
        isOpen={true}
        onClose={vi.fn()}
        platforms={mockPlatforms}
        onPlatformsChanged={handlePlatformsChanged}
      />
    );

    const removeBtns = screen.getAllByRole('button', { name: /Remove/i });
    // First remove button is for PC
    fireEvent.click(removeBtns[0]);

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/platforms/PC?force=true'),
        expect.objectContaining({
          method: 'DELETE',
        })
      );
      expect(handlePlatformsChanged).toHaveBeenCalledTimes(1);
    });
  });

  it('adds subcategory to an existing platform', async () => {
    const handlePlatformsChanged = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ message: 'Subcategory added' }),
    } as Response);

    render(
      <ManagePlatformsModal
        isOpen={true}
        onClose={vi.fn()}
        platforms={mockPlatforms}
        onPlatformsChanged={handlePlatformsChanged}
      />
    );

    const addSubBtns = screen.getAllByRole('button', { name: /Add Subplatform/i });
    // Click on PC's "Add Subplatform"
    fireEvent.click(addSubBtns[0]);

    const subInput = screen.getByPlaceholderText(/e\.g\. Ubisoft Connect/i);
    fireEvent.change(subInput, { target: { value: 'Ubisoft Connect' } });

    const saveBtn = screen.getByRole('button', { name: /Save/i });
    fireEvent.click(saveBtn);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/platforms/PC/subcategories'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ name: 'Ubisoft Connect' }),
        })
      );
      expect(handlePlatformsChanged).toHaveBeenCalledTimes(1);
    });
  });

  it('deletes subcategory from an existing platform', async () => {
    const handlePlatformsChanged = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ message: 'Subcategory deleted' }),
    } as Response);

    render(
      <ManagePlatformsModal
        isOpen={true}
        onClose={vi.fn()}
        platforms={mockPlatforms}
        onPlatformsChanged={handlePlatformsChanged}
      />
    );

    const removeSubBtn = screen.getByTitle('Remove Steam from PC');
    fireEvent.click(removeSubBtn);

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/platforms/PC/subcategories/Steam?force=true'),
        expect.objectContaining({
          method: 'DELETE',
        })
      );
      expect(handlePlatformsChanged).toHaveBeenCalledTimes(1);
    });
  });
});
