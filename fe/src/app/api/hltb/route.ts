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

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const query = searchParams.get('q');

  if (!query || !query.trim()) {
    return NextResponse.json({ error: 'Missing search query parameter q' }, { status: 400 });
  }

  try {
    // 1. Fetch search session token from HowLongToBeat
    const initRes = await fetch(`https://howlongtobeat.com/api/search/site/init?t=${Date.now()}`, {
      headers: {
        'User-Agent':
          'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36',
        Referer: 'https://howlongtobeat.com',
      },
      next: { revalidate: 3600 },
    });

    if (!initRes.ok) {
      return NextResponse.json(
        { error: 'Failed to initialize HLTB session', status: initRes.status },
        { status: 502 }
      );
    }

    const { token } = await initRes.json();
    if (!token) {
      return NextResponse.json({ error: 'No token returned from HLTB' }, { status: 502 });
    }

    // 2. Perform search on HLTB
    const searchTerms = query.trim().split(/\s+/);
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
        'Content-Type': 'application/json',
        'x-auth-token': token,
        'User-Agent':
          'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36',
        Referer: 'https://howlongtobeat.com',
      },
      body: JSON.stringify(payload),
    });

    if (!searchRes.ok) {
      return NextResponse.json(
        { error: 'Search request to HLTB failed', status: searchRes.status },
        { status: 502 }
      );
    }

    const searchData = await searchRes.json();
    const games: HLTBGameResult[] = searchData.data || [];

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
        : 'N/A';

    return NextResponse.json({
      found: true,
      gameId: top.game_id,
      title: top.game_name,
      mainHours: mainHours || undefined,
      extraHours: extraHours || undefined,
      completionistHours: compHours || undefined,
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
