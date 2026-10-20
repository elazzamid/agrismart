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


export type Plot = { id: string; farm_id: string; name: string; area_m2: number };
export type Crop = { id: string; code: string; name: string; description?: string; is_active: boolean };
export type CropCycle = {
  id: string;
  plot_id: string;
  crop_id: string;
  variety_id?: string | null;
  planting_date: string;
  status: string;
};

export async function listPlots(farmId: string): Promise<Plot[]> {
  return request<Plot[]>(`/api/v1/farms/${encodeURIComponent(farmId)}/plots`);
}

export async function createPlot(farmId: string, input: { name: string; area_m2: number }): Promise<Plot> {
  return request<Plot>(`/api/v1/farms/${encodeURIComponent(farmId)}/plots`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function listCrops(): Promise<Crop[]> {
  return request<Crop[]>("/api/v1/crops");
}

export async function listCropCycles(plotId: string): Promise<CropCycle[]> {
  return request<CropCycle[]>(`/api/v1/plots/${encodeURIComponent(plotId)}/crop-cycles`);
}

export async function createCropCycle(plotId: string, input: { crop_id: string; planting_date: string }): Promise<CropCycle> {
  return request<CropCycle>(`/api/v1/plots/${encodeURIComponent(plotId)}/crop-cycles`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}


export type FarmActivity = { id: string; crop_cycle_id: string; activity_type: string; activity_date: string; notes?: string };
export type FarmInput = { id: string; crop_cycle_id: string; input_type: string; product_name: string; quantity: number; unit: string; applied_date: string; notes?: string };
export type FarmExpense = { id: string; crop_cycle_id: string; expense_type: string; description: string; amount: number; expense_date: string; notes?: string };
export type FarmHarvest = { id: string; crop_cycle_id: string; quantity: number; unit: string; harvest_date: string; quality_grade?: string; notes?: string };

export const listActivities = (id: string) => request<FarmActivity[]>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/activities`);
export const createActivity = (id: string, input: Omit<FarmActivity, "id" | "crop_cycle_id">) =>
  request<FarmActivity>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/activities`, { method: "POST", body: JSON.stringify(input) });
export const listInputs = (id: string) => request<FarmInput[]>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/inputs`);
export const createInput = (id: string, input: Omit<FarmInput, "id" | "crop_cycle_id">) =>
  request<FarmInput>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/inputs`, { method: "POST", body: JSON.stringify(input) });
export const listExpenses = (id: string) => request<FarmExpense[]>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/expenses`);
export const createExpense = (id: string, input: Omit<FarmExpense, "id" | "crop_cycle_id">) =>
  request<FarmExpense>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/expenses`, { method: "POST", body: JSON.stringify(input) });
export const listHarvests = (id: string) => request<FarmHarvest[]>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/harvests`);
export const createHarvest = (id: string, input: Omit<FarmHarvest, "id" | "crop_cycle_id">) =>
  request<FarmHarvest>(`/api/v1/crop-cycles/${encodeURIComponent(id)}/harvests`, { method: "POST", body: JSON.stringify(input) });
