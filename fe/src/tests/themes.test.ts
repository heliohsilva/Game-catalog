import { describe, it, expect, beforeEach } from 'vitest';
import { 
  OMARCHY_THEMES, 
  DEFAULT_DARK_THEME_ID, 
  DEFAULT_LIGHT_THEME_ID, 
  applyThemeToCss,
  hexToRgba
} from '../data/themes';

describe('Omarchy Themes & CSS Customization', () => {
  beforeEach(() => {
    document.documentElement.style.cssText = '';
  });

  it('contains all 22 official Omarchy themes plus 2 quick defaults', () => {
    // 2 quick defaults (Cyber Blue, Cream Latte) + 22 official Omarchy themes = 24 themes
    expect(OMARCHY_THEMES.length).toBeGreaterThanOrEqual(24);
    
    const defaultDark = OMARCHY_THEMES.find(t => t.id === DEFAULT_DARK_THEME_ID);
    const defaultLight = OMARCHY_THEMES.find(t => t.id === DEFAULT_LIGHT_THEME_ID);

    expect(defaultDark).toBeDefined();
    expect(defaultDark?.name).toContain('Cyber Blue');
    expect(defaultDark?.mode).toBe('dark');

    expect(defaultLight).toBeDefined();
    expect(defaultLight?.name).toContain('Cream Latte');
    expect(defaultLight?.mode).toBe('light');
  });

  it('correctly converts 3-digit and 6-digit hex colors to rgba', () => {
    expect(hexToRgba('#fff', 0.5)).toBe('rgba(255, 255, 255, 0.5)');
    expect(hexToRgba('#000000', 0.8)).toBe('rgba(0, 0, 0, 0.8)');
    expect(hexToRgba('#38bdf8', 0.45)).toBe('rgba(56, 189, 248, 0.45)');
    expect(hexToRgba('invalid', 0.5)).toBe('rgba(255, 255, 255, 0.5)');
  });

  it('applies theme variables to document root element', () => {
    const cyberBlue = OMARCHY_THEMES.find(t => t.id === DEFAULT_DARK_THEME_ID)!;
    applyThemeToCss(cyberBlue);

    const root = document.documentElement;
    expect(root.getAttribute('data-theme')).toBe('default-dark');
    expect(root.getAttribute('data-mode')).toBe('dark');
    expect(root.style.getPropertyValue('--accent')).toBe(cyberBlue.accent);
    expect(root.style.getPropertyValue('--bg-primary')).toBe(cyberBlue.background);
    expect(root.style.getPropertyValue('--bg-card')).toContain('rgba');
    expect(root.style.getPropertyValue('--bg-pill')).toContain('rgba');
    expect(root.style.getPropertyValue('--border-color')).toContain('rgba');
  });

  it('applies light mode theme with warm cream latte properties', () => {
    const creamLatte = OMARCHY_THEMES.find(t => t.id === DEFAULT_LIGHT_THEME_ID)!;
    applyThemeToCss(creamLatte);

    const root = document.documentElement;
    expect(root.getAttribute('data-theme')).toBe('default-light');
    expect(root.getAttribute('data-mode')).toBe('light');
    expect(root.style.getPropertyValue('--bg-primary')).toBe('#FAF6EE');
  });
});
