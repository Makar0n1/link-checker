import { authApi, setTokens, clearTokens, getRefreshToken } from "./client";
import type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
} from "@/types/api";

export async function login(data: LoginRequest): Promise<AuthResponse> {
  const response = await authApi.post<AuthResponse>(
    "/api/v1/auth/login",
    data,
    false
  );
  setTokens(response);
  return response;
}

export async function register(data: RegisterRequest): Promise<User> {
  return authApi.post<User>("/api/v1/auth/register", data, false);
}

export async function logout(): Promise<void> {
  const refreshToken = getRefreshToken();
  if (refreshToken) {
    try {
      await authApi.post("/api/v1/auth/logout", { refresh_token: refreshToken });
    } catch {
      // Ignore logout errors, clear tokens anyway
    }
  }
  clearTokens();
}

export async function getCurrentUser(): Promise<User> {
  return authApi.get<User>("/api/v1/auth/me");
}
