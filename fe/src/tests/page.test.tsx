import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import GameCatalogPage from '../app/page';

describe('GameCatalogPage (Main Index Page)', () => {
  const mockPlatforms = [
    { name: 'PC', subcategories: ['Steam', 'GOG', 'Epic'] },
    { name: 'Nintendo Switch', subcategories: ['Switch'] },
    { name: 'PlayStation', subcategories: ['PS5'] },
    { name: 'Xbox', subcategories: ['Xbox Series X/S'] },
    { name: 'Retro / Emulation', subcategories: ['PS1'] },
    { name: 'Mega Drive', subcategories: ['Genesis'] },
  ];

  const mockGames = [
    { id: '1', title: 'Bloodborne', platform: 'PlayStation', genre: 'Action RPG', subcategory: 'PS5' },
    { id: '2', title: 'Demon Souls', platform: 'PlayStation', genre: 'Action RPG', subcategory: 'PS5' },
    { id: '3', title: 'God of War', platform: 'PlayStation', genre: 'Action', subcategory: 'PS5' },
    { id: '4', title: 'Half-Life 2', platform: 'PC', genre: 'FPS', subcategory: 'Steam' },
  ];

  beforeEach(() => {
    vi.restoreAllMocks();
    global.fetch = vi.fn().mockImplementation(async (url: string) => {
      if (url.includes('/platforms')) {
        return {
          ok: true,
          json: async () => ({ platforms: mockPlatforms }),
        } as Response;
      }
      if (url.includes('/games')) {
        return {
          ok: true,
          json: async () => ({ data: mockGames }),
        } as Response;
      }
      return { ok: true, json: async () => ({}) } as Response;
    });
  });

  it('sorts platforms on index by amount of games (descending) and hides platforms without games', async () => {
    render(<GameCatalogPage />);

    // Wait for games to load
    await waitFor(() => {
      expect(screen.getByText('Bloodborne')).toBeInTheDocument();
    });

    // Check headings (platform cluster titles)
    const clusterHeadings = screen.getAllByRole('heading', { level: 2 });
    const clusterNames = clusterHeadings.map((h) => h.textContent?.trim());

    // PlayStation has 3 games, PC has 1 game -> PlayStation first, PC second
    expect(clusterNames).toEqual(['PlayStation', 'PC']);

    // Nintendo Switch, Xbox, Retro, Mega Drive have 0 games -> should NOT be shown
    expect(screen.queryByRole('heading', { level: 2, name: 'Nintendo Switch' })).not.toBeInTheDocument();
    expect(screen.queryByRole('heading', { level: 2, name: 'Xbox' })).not.toBeInTheDocument();
    expect(screen.queryByRole('heading', { level: 2, name: 'Retro / Emulation' })).not.toBeInTheDocument();
    expect(screen.queryByRole('heading', { level: 2, name: 'Mega Drive' })).not.toBeInTheDocument();
  });

  it('renders dynamic platform in navbar with real name rather than PC', async () => {
    render(<GameCatalogPage />);

    await waitFor(() => {
      expect(screen.getByText('Bloodborne')).toBeInTheDocument();
    });

    // Navbar should have button with "Mega Drive"
    expect(screen.getByRole('button', { name: 'Mega Drive' })).toBeInTheDocument();

    // Verify there is only one "PC" button in the navbar
    const pcButtons = screen.getAllByRole('button', { name: 'PC' });
    expect(pcButtons).toHaveLength(1);
  });

  it('allows viewing a specific platform cluster with 0 games when selected directly in navbar', async () => {
    render(<GameCatalogPage />);

    await waitFor(() => {
      expect(screen.getByText('Bloodborne')).toBeInTheDocument();
    });

    // Click on Mega Drive tab in navbar
    const megaDriveTab = screen.getByRole('button', { name: 'Mega Drive' });
    fireEvent.click(megaDriveTab);

    // Now Mega Drive cluster should be visible
    expect(screen.getByRole('heading', { level: 2, name: 'Mega Drive' })).toBeInTheDocument();
    expect(screen.getByText('No games in this cluster.')).toBeInTheDocument();

    // Switch back to "All" (index)
    fireEvent.click(screen.getByRole('button', { name: 'All' }));

    // Mega Drive cluster should be hidden again on index
    expect(screen.queryByRole('heading', { level: 2, name: 'Mega Drive' })).not.toBeInTheDocument();
    // PlayStation and PC should still be there
    expect(screen.getByRole('heading', { level: 2, name: 'PlayStation' })).toBeInTheDocument();
  });

  it('hides empty clusters when searching', async () => {
    render(<GameCatalogPage />);

    await waitFor(() => {
      expect(screen.getByText('Bloodborne')).toBeInTheDocument();
    });

    const searchInput = screen.getByPlaceholderText(/Filter title, store, or genre/i);
    fireEvent.change(searchInput, { target: { value: 'Half-Life' } });

    // PC has Half-Life 2, PlayStation has none matching Half-Life
    expect(screen.getByRole('heading', { level: 2, name: 'PC' })).toBeInTheDocument();
    expect(screen.queryByRole('heading', { level: 2, name: 'PlayStation' })).not.toBeInTheDocument();
  });
});
