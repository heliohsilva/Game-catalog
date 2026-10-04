import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ThemeSelectorModal } from '../components/ThemeSelectorModal';
import { OMARCHY_THEMES } from '../data/themes';

describe('ThemeSelectorModal Component', () => {
  const currentTheme = OMARCHY_THEMES[0];

  it('renders theme selector with quick defaults and Omarchy theme catalog', () => {
    render(
      <ThemeSelectorModal
        isOpen={true}
        onClose={vi.fn()}
        currentTheme={currentTheme}
        onSelectTheme={vi.fn()}
      />
    );

    expect(screen.getByText('Theme Selector')).toBeInTheDocument();
    expect(screen.getByText('Default Options')).toBeInTheDocument();
    expect(screen.getByText('Cyber Blue (Default Dark)')).toBeInTheDocument();
    expect(screen.getByText('Cream Latte (Default Light)')).toBeInTheDocument();
    expect(screen.getByText(/All 22 Omarchy Themes/i)).toBeInTheDocument();
  });

  it('filters themes by text search', () => {
    render(
      <ThemeSelectorModal
        isOpen={true}
        onClose={vi.fn()}
        currentTheme={currentTheme}
        onSelectTheme={vi.fn()}
      />
    );

    const searchInput = screen.getByPlaceholderText(/Search theme name/i);
    fireEvent.change(searchInput, { target: { value: 'Gruvbox' } });

    expect(screen.getByText('Gruvbox')).toBeInTheDocument();
    expect(screen.queryByText('Everforest')).not.toBeInTheDocument();
  });

  it('filters themes by Dark / Light mode toggle', () => {
    render(
      <ThemeSelectorModal
        isOpen={true}
        onClose={vi.fn()}
        currentTheme={currentTheme}
        onSelectTheme={vi.fn()}
      />
    );

    const lightBtn = screen.getByRole('button', { name: /^Light$/i });
    fireEvent.click(lightBtn);

    // Should display light themes like Catppuccin Latte
    expect(screen.getByText('Catppuccin Latte')).toBeInTheDocument();
  });

  it('calls onSelectTheme when clicking a theme option', () => {
    const handleSelectTheme = vi.fn();
    render(
      <ThemeSelectorModal
        isOpen={true}
        onClose={vi.fn()}
        currentTheme={currentTheme}
        onSelectTheme={handleSelectTheme}
      />
    );

    const creamLatteOption = screen.getByText('Cream Latte (Default Light)');
    fireEvent.click(creamLatteOption);

    expect(handleSelectTheme).toHaveBeenCalledWith('default-light');
  });

  it('closes modal when clicking outside on the backdrop (click-outside-to-dismiss)', () => {
    const handleClose = vi.fn();
    const { container } = render(
      <ThemeSelectorModal
        isOpen={true}
        onClose={handleClose}
        currentTheme={currentTheme}
        onSelectTheme={vi.fn()}
      />
    );

    // The backdrop is the outer fixed container
    const backdrop = container.firstChild as HTMLElement;
    fireEvent.click(backdrop);

    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
