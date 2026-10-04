import { NextRequest, NextResponse } from 'next/server';

async function handleProxy(
  request: NextRequest,
  context: { params: Promise<{ path: string[] }> | { path: string[] } }
) {
  try {
    const resolvedParams = await Promise.resolve(context.params);
    const pathSegments = Array.isArray(resolvedParams?.path) ? resolvedParams.path : [];
    const targetPath = pathSegments.join('/');
    const search = request.nextUrl.search || '';

    // In Docker, defaults to internal service 'http://api:8080/api/v1'.
    // Outside Docker, defaults to 'http://localhost:8080/api/v1'.
    const baseUrl = process.env.INTERNAL_API_URL || 'http://localhost:8080/api/v1';
    const targetUrl = `${baseUrl.replace(/\/+$/, '')}/${targetPath}${search}`;

    const headers = new Headers();
    request.headers.forEach((val, key) => {
      // Avoid forwarding host or connection headers
      if (!['host', 'connection', 'content-length'].includes(key.toLowerCase())) {
        headers.set(key, val);
      }
    });

    const isBodyAllowed = !['GET', 'HEAD'].includes(request.method.toUpperCase());
    let body: ArrayBuffer | undefined = undefined;
    if (isBodyAllowed) {
      body = await request.arrayBuffer();
    }

    const response = await fetch(targetUrl, {
      method: request.method,
      headers,
      body,
      cache: 'no-store',
    });

    const resHeaders = new Headers();
    response.headers.forEach((val, key) => {
      if (!['content-encoding', 'transfer-encoding', 'connection'].includes(key.toLowerCase())) {
        resHeaders.set(key, val);
      }
    });

    const responseData = await response.arrayBuffer();
    return new NextResponse(responseData, {
      status: response.status,
      statusText: response.statusText,
      headers: resHeaders,
    });
  } catch (err: unknown) {
    const errMsg = err instanceof Error ? err.message : 'Unknown proxy error';
    return NextResponse.json(
      { error: { code: 'PROXY_ERROR', message: `Internal API proxy error: ${errMsg}` } },
      { status: 502 }
    );
  }
}

export const GET = handleProxy;
export const POST = handleProxy;
export const PUT = handleProxy;
export const DELETE = handleProxy;
export const PATCH = handleProxy;
