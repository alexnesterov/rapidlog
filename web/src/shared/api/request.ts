import { ApiError, type ApiEnvelope } from "./types";

const API_BASE = "/api";

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
  });

  const body = (await res.json()) as ApiEnvelope<T>;

  if (!res.ok) {
    throw new ApiError(String(body.error?.code ?? res.status), body.error?.message ?? "request failed");
  }

  return body.data;
}
