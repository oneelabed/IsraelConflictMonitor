import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const apiKey = request.cookies.get('api_key')?.value;
  const role = request.cookies.get('role')?.value;
  const { pathname } = request.nextUrl;

  if (pathname.startsWith('/admin')) {
    if (!apiKey || role !== 'admin') {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  if (pathname.startsWith('/posts') || pathname.startsWith('/feeds')) {
    if (!apiKey) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  return NextResponse.next();
}

// Only apply this middleware to these specific routes
export const config = {
  matcher: ['/posts/:path*', '/admin/:path*', '/feeds/:path*'],
};