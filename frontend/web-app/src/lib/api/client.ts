import type { AuthResponse, ApiError } from "@/types/api";

// Use relative URLs in production (through nginx), absolute in development
const AUTH_API_URL = process.env.NEXT_PUBLIC_API_URL || "";
const BACKLINK_API_URL = process.env.NEXT_PUBLIC_BACKLINK_API_URL || "";

// Token storage keys
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const TOKEN_EXPIRY_KEY = "token_expiry";

// Token management
export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setTokens(authResponse: AuthResponse): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(ACCESS_TOKEN_KEY, authResponse.access_token);
  localStorage.setItem(REFRESH_TOKEN_KEY, authResponse.refresh_token);
  const expiry = Date.now() + authResponse.expires_in * 1000;
  localStorage.setItem(TOKEN_EXPIRY_KEY, expiry.toString());
  // Also set cookie for middleware
  document.cookie = `auth-token=${authResponse.access_token}; path=/; max-age=${authResponse.expires_in}`;
}

export function clearTokens(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(TOKEN_EXPIRY_KEY);
  document.cookie = "auth-token=; path=/; max-age=0";
}

export function isTokenExpired(): boolean {
  if (typeof window === "undefined") return true;
  const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
  if (!expiry) return true;
  // Add 30 second buffer
  return Date.now() > parseInt(expiry) - 30000;
}

// Refresh token
async function refreshAccessToken(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;

  try {
    const response = await fetch(`${AUTH_API_URL}/api/v1/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      clearTokens();
      return false;
    }

    const data: AuthResponse = await response.json();
    setTokens(data);
    return true;
  } catch {
    clearTokens();
    return false;
  }
}

// API client class
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async getAuthHeaders(): Promise<HeadersInit> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    if (isTokenExpired()) {
      const refreshed = await refreshAccessToken();
      if (!refreshed) {
        throw new ApiClientError("Session expired", 401);
      }
    }

    const token = getAccessToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    return headers;
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {},
    requireAuth = true
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (requireAuth) {
      Object.assign(headers, await this.getAuthHeaders());
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      let errorData: ApiError = { error: "Unknown error" };
      try {
        errorData = await response.json();
      } catch {
        // Response may not be JSON
      }
      throw new ApiClientError(errorData.error, response.status, errorData);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return undefined as T;
    }

    return response.json();
  }

  async get<T>(endpoint: string, requireAuth = true): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" }, requireAuth);
  }

  async post<T>(endpoint: string, data?: unknown, requireAuth = true): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: "POST",
        body: data ? JSON.stringify(data) : undefined,
      },
      requireAuth
    );
  }

  async put<T>(endpoint: string, data?: unknown, requireAuth = true): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: "PUT",
        body: data ? JSON.stringify(data) : undefined,
      },
      requireAuth
    );
  }

  async delete<T>(endpoint: string, data?: unknown, requireAuth = true): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: "DELETE",
        body: data ? JSON.stringify(data) : undefined,
      },
      requireAuth
    );
  }
}

// Custom error class
export class ApiClientError extends Error {
  status: number;
  data?: ApiError;

  constructor(message: string, status: number, data?: ApiError) {
    super(message);
    this.name = "ApiClientError";
    this.status = status;
    this.data = data;
  }
}

// Export configured clients
export const authApi = new ApiClient(AUTH_API_URL);
export const backlinkApi = new ApiClient(BACKLINK_API_URL);
