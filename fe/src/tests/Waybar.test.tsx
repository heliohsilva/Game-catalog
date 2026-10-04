import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Waybar } from '../components/Waybar';
import { PlatformInfo } from '../types/game';
import { OMARCHY_THEMES, DEFAULT_DARK_THEME_ID, DEFAULT_LIGHT_THEME_ID } from '../data/themes';

describe('Waybar Component', () => {
  const currentDarkTheme = OMARCHY_THEMES[0]; // Cyber Blue (Default Dark)
  const currentLightTheme = OMARCHY_THEMES[1]; // Cream Latte (Default Light)

  it('renders catalog brand, platform switchers, and controls', () => {
    render(
      <Waybar
        currentTheme={currentDarkTheme}
        onThemeSelect={vi.fn()}
        onOpenThemeModal={vi.fn()}
        onOpenAddModal={vi.fn()}
        onOpenSearchModal={vi.fn()}
        selectedPlatform="All"
        onSelectPlatform={vi.fn()}
        totalGames={32}
      />
    );

    expect(screen.getByText('CATALOG')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'All' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'PC' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Switch' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'PlayStation' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Xbox' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Retro' })).toBeInTheDocument();
  });

  it('selects platform when clicking a platform tab', () => {
    const handleSelectPlatform = vi.fn();
    render(
      <Waybar
        currentTheme={currentDarkTheme}
        onThemeSelect={vi.fn()}
        onOpenThemeModal={vi.fn()}
        onOpenAddModal={vi.fn()}
        onOpenSearchModal={vi.fn()}
        selectedPlatform="All"
        onSelectPlatform={handleSelectPlatform}
        totalGames={32}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'Retro' }));
    expect(handleSelectPlatform).toHaveBeenCalledWith('Retro / Emulation');
  });

  it('toggles quick dark/light theme between Cyber Blue and Cream Latte', () => {
    const handleThemeSelect = vi.fn();
    const { rerender } = render(
      <Waybar
        currentTheme={currentDarkTheme}
        onThemeSelect={handleThemeSelect}
        onOpenThemeModal={vi.fn()}
        onOpenAddModal={vi.fn()}
        onOpenSearchModal={vi.fn()}
        selectedPlatform="All"
        onSelectPlatform={vi.fn()}
        totalGames={32}
      />
    );

    const toggleBtn = screen.getByTitle(/Switch to Cream Latte/i);
    fireEvent.click(toggleBtn);
    expect(handleThemeSelect).toHaveBeenCalledWith(DEFAULT_LIGHT_THEME_ID);

    // When currently in light mode, should switch to dark
    rerender(
      <Waybar
        currentTheme={currentLightTheme}
        onThemeSelect={handleThemeSelect}
        onOpenThemeModal={vi.fn()}
        onOpenAddModal={vi.fn()}
        onOpenSearchModal={vi.fn()}
        selectedPlatform="All"
        onSelectPlatform={vi.fn()}
        totalGames={32}
      />
    );

    const toggleDarkBtn = screen.getByTitle(/Switch to Cyber Blue/i);
    fireEvent.click(toggleDarkBtn);
    expect(handleThemeSelect).toHaveBeenCalledWith(DEFAULT_DARK_THEME_ID);
  });

  it('triggers search, add, and theme modal actions', () => {
    const handleSearch = vi.fn();
    const handleAdd = vi.fn();
    const handleThemeModal = vi.fn();

    render(
      <Waybar
        currentTheme={currentDarkTheme}
        onThemeSelect={vi.fn()}
        onOpenThemeModal={handleThemeModal}
        onOpenAddModal={handleAdd}
        onOpenSearchModal={handleSearch}
        selectedPlatform="All"
        onSelectPlatform={vi.fn()}
        totalGames={32}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: /Search/i }));
    expect(handleSearch).toHaveBeenCalledTimes(1);

    fireEvent.click(screen.getByRole('button', { name: /Add/i }));
    expect(handleAdd).toHaveBeenCalledTimes(1);

    fireEvent.click(screen.getByTitle(/Browse all 22 Omarchy Themes/i));
    expect(handleThemeModal).toHaveBeenCalledTimes(1);
  });

  it('renders custom added platforms with their real names rather than PC', () => {
    const customPlatforms: PlatformInfo[] = [
      { name: 'PC', subcategories: ['Steam'] },
      { name: 'PlayStation', subcategories: ['PS5'] },
      { name: 'Nintendo Switch', subcategories: ['Switch'] },
      { name: 'Xbox', subcategories: ['Xbox Series X/S'] },
      { name: 'Retro / Emulation', subcategories: ['PS1'] },
      { name: 'Steam Deck', subcategories: ['Steam'] },
      { name: 'Nintendo 64', subcategories: ['N64'] },
    ];

    render(
      <Waybar
        currentTheme={currentDarkTheme}
        onThemeSelect={vi.fn()}
        onOpenThemeModal={vi.fn()}
        onOpenAddModal={vi.fn()}
        onOpenSearchModal={vi.fn()}
        selectedPlatform="All"
        onSelectPlatform={vi.fn()}
        totalGames={10}
        platforms={customPlatforms}
      />
    );

    expect(screen.getByRole('button', { name: 'Steam Deck' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Nintendo 64' })).toBeInTheDocument();
    const pcButtons = screen.getAllByRole('button', { name: 'PC' });
    expect(pcButtons).toHaveLength(1);
  });
});
