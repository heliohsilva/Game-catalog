import { NextRequest, NextResponse } from 'next/server';

interface HLTBGameResult {
  game_id: number;
  game_name: string;
  comp_main: number; // in seconds
  comp_plus: number;
  comp_100: number;
  comp_all: number;
  release_world: number;
}

const HLTB_HEADERS = {
  'User-Agent':
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
  Referer: 'https://howlongtobeat.com/',
  Origin: 'https://howlongtobeat.com',
  Accept: '*/*',
  'Accept-Language': 'en-US,en;q=0.9',
};

// In-memory token cache to avoid hammering /init on every keystroke
let cachedToken: { token: string; expiresAt: number } | null = null;

async function getHLTBToken(forceRefresh = false): Promise<string | null> {
  const now = Date.now();
  if (!forceRefresh && cachedToken && cachedToken.expiresAt > now) {
    return cachedToken.token;
  }

  try {
    const initRes = await fetch(`https://howlongtobeat.com/api/search/site/init?t=${now}`, {
      headers: HLTB_HEADERS,
      cache: 'no-store',
    });

    if (!initRes.ok) {
      return null;
    }

    const data = await initRes.json();
    if (data?.token) {
      // Cache token for 5 minutes
      cachedToken = {
        token: data.token,
        expiresAt: now + 5 * 60 * 1000,
      };
      return data.token;
    }
  } catch {
    // Network or parse error
  }

  return null;
}

async function performSearch(searchTerms: string[], token: string): Promise<HLTBGameResult[]> {
  const payload = {
    searchType: 'games',
    searchTerms,
    searchPage: 1,
    size: 5,
    searchOptions: {
      games: {
        userId: 0,
        platform: '',
        sortCategory: 'popular',
        rangeCategory: 'main',
        rangeTime: { min: 0, max: 0 },
        gameplay: { perspective: '', flow: '', genre: '' },
        modifier: '',
      },
    },
    useCache: true,
  };

  const searchRes = await fetch('https://howlongtobeat.com/api/search/site', {
    method: 'POST',
    headers: {
      ...HLTB_HEADERS,
      'Content-Type': 'application/json',
      'x-auth-token': token,
    },
    body: JSON.stringify(payload),
    cache: 'no-store',
  });

  if (!searchRes.ok) {
    throw new Error(`HLTB search HTTP error: ${searchRes.status}`);
  }

  const searchData = await searchRes.json();
  return searchData.data || [];
}

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const query = searchParams.get('q');

  if (!query || !query.trim()) {
    return NextResponse.json({ error: 'Missing search query parameter q' }, { status: 400 });
  }

  // Sanitize search query: remove punctuation that confuses HLTB backend (e.g. ':', '-', ''')
  const cleanTerms = query
    .replace(/[^a-zA-Z0-9\s]/g, ' ')
    .trim()
    .split(/\s+/)
    .filter(Boolean);

  if (cleanTerms.length === 0) {
    return NextResponse.json({ found: false, message: 'Invalid search query' });
  }

  try {
    let token = await getHLTBToken(false);
    if (!token) {
      token = await getHLTBToken(true);
    }

    if (!token) {
      return NextResponse.json(
        { error: 'Failed to obtain HowLongToBeat session token' },
        { status: 502 }
      );
    }

    let games: HLTBGameResult[] = [];
    try {
      games = await performSearch(cleanTerms, token);
    } catch {
      // If failed, token might have been invalidated; retry once with a fresh token
      token = await getHLTBToken(true);
      if (token) {
        games = await performSearch(cleanTerms, token);
      }
    }

    // If no results and multiple words, attempt searching with first 2 words (e.g. omitting subtitle)
    if (games.length === 0 && cleanTerms.length > 2 && token) {
      try {
        const fallbackTerms = cleanTerms.slice(0, 2);
        games = await performSearch(fallbackTerms, token);
      } catch {
        // Ignore fallback error
      }
    }

    if (games.length === 0) {
      return NextResponse.json({
        found: false,
        message: 'No results found on HowLongToBeat',
      });
    }

    // Top match
    const top = games[0];
    const mainHours = Math.round(top.comp_main / 3600);
    const extraHours = Math.round(top.comp_plus / 3600);
    const compHours = Math.round(top.comp_100 / 3600);

    const formatted =
      mainHours > 0
        ? `${mainHours}h`
        : extraHours > 0
        ? `${extraHours}h`
        : compHours > 0
        ? `${compHours}h`
        : 'N/A';

    return NextResponse.json({
      found: true,
      gameId: top.game_id,
      title: top.game_name,
      mainHours: mainHours > 0 ? mainHours : undefined,
      extraHours: extraHours > 0 ? extraHours : undefined,
      completionistHours: compHours > 0 ? compHours : undefined,
      timeToBeat: formatted,
      hltbUrl: `https://howlongtobeat.com/game/${top.game_id}`,
      matches: games.map((g) => ({
        id: g.game_id,
        name: g.game_name,
        mainHours: Math.round(g.comp_main / 3600),
      })),
    });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    return NextResponse.json({ error: message }, { status: 500 });
  }
}
