import { describe, it, expect, vi, beforeEach } from 'vitest';
import { NextRequest } from 'next/server';
import { GET } from '../app/api/hltb/route';

describe('HowLongToBeat API Route Handler (/api/hltb)', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('returns 400 when search query q is missing', async () => {
    const req = new NextRequest('http://localhost:3000/api/hltb');
    const res = await GET(req);
    expect(res.status).toBe(400);

    const data = await res.json();
    expect(data.error).toMatch(/Missing search query/i);
  });

  it('returns found: false when query contains only punctuation or whitespace', async () => {
    const req = new NextRequest('http://localhost:3000/api/hltb?q=---:::');
    const res = await GET(req);
    expect(res.status).toBe(200);

    const data = await res.json();
    expect(data.found).toBe(false);
  });

  it('returns found: true and formatted timeToBeat on successful HLTB lookup', async () => {
    const mockInitResponse = { token: 'mock-token-abc' };
    const mockSearchResponse = {
      data: [
        {
          game_id: 1705,
          game_name: 'Chrono Trigger',
          comp_main: 82800, // 23 hours in seconds
          comp_plus: 97200, // 27 hours
          comp_100: 151200, // 42 hours
          comp_all: 95000,
          release_world: 1995,
        },
      ],
    };

    vi.spyOn(global, 'fetch')
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockInitResponse,
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockSearchResponse,
      } as Response);

    const req = new NextRequest('http://localhost:3000/api/hltb?q=Chrono+Trigger');
    const res = await GET(req);
    expect(res.status).toBe(200);

    const data = await res.json();
    expect(data.found).toBe(true);
    expect(data.gameId).toBe(1705);
    expect(data.title).toBe('Chrono Trigger');
    expect(data.mainHours).toBe(23);
    expect(data.timeToBeat).toBe('23h');
  });

  it('returns found: false when HLTB has 0 results for the query', async () => {
    const mockInitResponse = { token: 'mock-token-xyz' };
    const mockSearchResponse = { data: [] };

    vi.spyOn(global, 'fetch')
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockInitResponse,
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockSearchResponse,
      } as Response);

    const req = new NextRequest('http://localhost:3000/api/hltb?q=NonExistentGameXYZ');
    const res = await GET(req);
    expect(res.status).toBe(200);

    const data = await res.json();
    expect(data.found).toBe(false);
  });
});
