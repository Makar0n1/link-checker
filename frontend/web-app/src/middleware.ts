import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const token = request.cookies.get("auth-token")?.value;

  const isAuthPage =
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register");

  const isDashboard =
    request.nextUrl.pathname.startsWith("/backlinks") ||
    request.nextUrl.pathname.startsWith("/index-checker") ||
    request.nextUrl.pathname.startsWith("/site-health");

  // Если нет токена и пытается зайти в dashboard — редирект на login
  if (!token && isDashboard) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // Если есть токен и на auth странице — редирект в dashboard
  if (token && isAuthPage) {
    return NextResponse.redirect(new URL("/backlinks", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/backlinks/:path*",
    "/index-checker/:path*",
    "/site-health/:path*",
    "/login",
    "/register",
  ],
};
