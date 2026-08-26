import {
  ApiError,
  type ApiEnvelope,
  type Bullet,
  type BulletDayGroup,
  type CreateBulletRequest,
} from "../types/bullet";

const API_BASE = "/api";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
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

export async function listBullets(): Promise<BulletDayGroup[]> {
  return request<BulletDayGroup[]>("/bullets");
}

export async function createBullet(req: CreateBulletRequest): Promise<Bullet> {
  return request<Bullet>("/bullets", {
    method: "POST",
    body: JSON.stringify(req),
  });
}

export async function markBulletDone(id: string): Promise<Bullet> {
  return request<Bullet>(`/bullets/${id}/complete`, {
    method: "POST",
  });
}

export async function migrateBullet(id: string): Promise<Bullet> {
  return request<Bullet>(`/bullets/${id}/migrate`, {
    method: "POST",
  });
}
