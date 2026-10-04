import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { PlatformFilterBar } from '../components/PlatformFilterBar';

describe('PlatformFilterBar Component', () => {
  it('renders total game count and search input', () => {
    render(
      <PlatformFilterBar
        searchQuery=""
        onSearchChange={vi.fn()}
        sortBy="title"
        onSortChange={vi.fn()}
        totalFiltered={12}
      />
    );

    expect(screen.getByText('Catalog (12 games)')).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Filter title, store, or genre/i)).toBeInTheDocument();
  });

  it('triggers onSearchChange when typing into input', () => {
    const handleSearch = vi.fn();
    render(
      <PlatformFilterBar
        searchQuery=""
        onSearchChange={handleSearch}
        sortBy="title"
        onSortChange={vi.fn()}
        totalFiltered={12}
      />
    );

    const input = screen.getByPlaceholderText(/Filter title, store, or genre/i);
    fireEvent.change(input, { target: { value: 'Zelda' } });
    expect(handleSearch).toHaveBeenCalledWith('Zelda');
  });

  it('triggers onSortChange when changing the sort dropdown', () => {
    const handleSort = vi.fn();
    render(
      <PlatformFilterBar
        searchQuery=""
        onSearchChange={vi.fn()}
        sortBy="title"
        onSortChange={handleSort}
        totalFiltered={12}
      />
    );

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: 'hours' } });
    expect(handleSort).toHaveBeenCalledWith('hours');
  });
});
