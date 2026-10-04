export type Platform =
  | 'PC'
  | 'Nintendo Switch'
  | 'PlayStation'
  | 'Xbox'
  | 'Retro / Emulation';

export const PLATFORMS: Platform[] = [
  'PC',
  'Nintendo Switch',
  'PlayStation',
  'Xbox',
  'Retro / Emulation',
];

export const PLATFORM_SUBCATEGORIES: Record<Platform, string[]> = {
  'PC': ['Steam', 'GOG', 'Epic'],
  'PlayStation': ['PS2', 'PS3', 'PS4', 'PS5'],
  'Nintendo Switch': ['Switch'],
  'Xbox': ['Xbox Series X/S', 'Xbox One', 'Xbox 360'],
  'Retro / Emulation': ['PS1', 'N64', 'SNES', 'Genesis', 'NES', 'Master System', 'Neo Geo'],
};

export interface Game {
  id: string;
  title: string;
  platform: Platform;
  subcategory?: string;
  genre: string;
  timeToBeat?: string; // Retrieved automatically from HowLongToBeat
  timeToBeatMain?: number;
  addedAt?: string;
}

export interface OmarchyTheme {
  id: string;
  name: string;
  mode: 'dark' | 'light';
  description?: string;
  isQuickDefault?: boolean;
  accent: string;
  selection: string;
  muted: string;
  background: string;
  dark_background: string;
  darker_background: string;
  lighter_background: string;
  foreground: string;
  dark_foreground: string;
  light_foreground: string;
  bright_foreground: string;
  hyprland_active_border?: string;
  red?: string;
  yellow?: string;
  orange?: string;
  green?: string;
  cyan?: string;
  blue?: string;
  magenta?: string;
  brown?: string;
  bright_red?: string;
  bright_yellow?: string;
  bright_green?: string;
  bright_cyan?: string;
  bright_blue?: string;
  bright_magenta?: string;
}
