import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Simple JWT validation - check if token is a valid JWT structure
function isValidJwtFormat(token: string): boolean {
  if (!token || token === "mock-token") return false;
  const parts = token.split(".");
  if (parts.length !== 3) return false;
  try {
    // Check if parts are valid base64
    for (const part of parts) {
      atob(part.replace(/-/g, "+").replace(/_/g, "/"));
    }
    return true;
  } catch {
    return false;
  }
}

export function middleware(request: NextRequest) {
  const token = request.cookies.get("auth-token")?.value;
  const hasValidToken = token && isValidJwtFormat(token);

  const isAuthPage =
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register");

  const isDashboard =
    request.nextUrl.pathname.startsWith("/backlinks") ||
    request.nextUrl.pathname.startsWith("/index-checker") ||
    request.nextUrl.pathname.startsWith("/site-health");

  // Если нет валидного токена и пытается зайти в dashboard — редирект на login
  if (!hasValidToken && isDashboard) {
    const response = NextResponse.redirect(new URL("/login", request.url));
    // Clear invalid token if exists
    if (token) {
      response.cookies.set("auth-token", "", { path: "/", maxAge: 0 });
    }
    return response;
  }

  // Если есть валидный токен и на auth странице — редирект в dashboard
  if (hasValidToken && isAuthPage) {
    return NextResponse.redirect(new URL("/backlinks", request.url));
  }

  // Clear invalid token on auth pages
  if (!hasValidToken && token && isAuthPage) {
    const response = NextResponse.next();
    response.cookies.set("auth-token", "", { path: "/", maxAge: 0 });
    return response;
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
