import { isoDateFromTimestamp } from "../lib/date";
import {
  ApiError,
  type ApiEnvelope,
  type Bullet,
  type CreateBulletRequest,
  type ListBulletsResponse,
  type RawBullet,
} from "../types/bullet";

const API_BASE = "/api";

// пометка "выполнено" не привязана к бэкенду — там ещё нет PUT /api/bullets/{id},
// поэтому статус живёт только в этой вкладке и не переживает перезагрузку страницы
const doneOverrides = new Set<string>();

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

function toBullet(raw: RawBullet): Bullet {
  return {
    ...raw,
    status: doneOverrides.has(raw.id) ? "DONE" : raw.status,
    date: isoDateFromTimestamp(raw.created_at),
  };
}

export async function listBullets(): Promise<ListBulletsResponse> {
  const bullets = await request<RawBullet[]>("/bullets");
  return { bullets: bullets.map(toBullet) };
}

export async function createBullet(req: CreateBulletRequest): Promise<Bullet> {
  const raw = await request<RawBullet>("/bullets", {
    method: "POST",
    body: JSON.stringify(req),
  });
  return toBullet(raw);
}

export async function markBulletDone(id: string): Promise<void> {
  doneOverrides.add(id);
}
