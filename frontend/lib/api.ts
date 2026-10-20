"use client";

export type User = { id: string; email: string; role: string };
export type Farm = {
  id: string;
  name: string;
  location_name?: string;
  latitude?: number;
  longitude?: number;
};

type AuthResponse = { access_token: string; token_type: "Bearer"; user: User };
type APIError = { error?: string };

const TOKEN_KEY = "agrismart_access_token";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.sessionStorage.getItem(TOKEN_KEY);
}

export function saveToken(token: string): void {
  window.sessionStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  window.sessionStorage.removeItem(TOKEN_KEY);
}

async function request<T>(path: string, init: RequestInit = {}, token = getToken()): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (init.body) headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  const response = await fetch(path, { ...init, headers, cache: "no-store" });
  if (!response.ok) {
    let message = `Permintaan gagal (${response.status})`;
    try {
      const body = (await response.json()) as APIError;
      if (body.error) message = body.error;
    } catch {
      // Keep the safe status-based message when the response isn't JSON.
    }
    if (response.status === 401) {
      clearToken();
      if (typeof window !== "undefined") window.dispatchEvent(new Event("agrismart:unauthorized"));
    }
    throw new Error(message);
  }
  return (await response.json()) as T;
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  }, null);
}

export async function register(input: {
  email: string;
  password: string;
  full_name: string;
  phone: string;
}): Promise<User> {
  return request<User>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  }, null);
}

export async function currentUser(): Promise<User> {
  return request<User>("/api/v1/auth/me");
}

export async function listFarms(): Promise<Farm[]> {
  return request<Farm[]>("/api/v1/farms");
}

export async function createFarm(input: { name: string; location_name: string }): Promise<Farm> {
  return request<Farm>("/api/v1/farms", {
    method: "POST",
    body: JSON.stringify(input),
  });
}
