import {
  ApiError,
  type ApiEnvelope,
  type Bullet,
  type BulletDayGroup,
  type CreateBulletRequest,
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

function applyDoneOverride(bullet: Bullet): Bullet {
  return doneOverrides.has(bullet.id) ? { ...bullet, signifier: "completed" } : bullet;
}

export async function listBullets(): Promise<BulletDayGroup[]> {
  const groups = await request<BulletDayGroup[]>("/bullets");
  return groups.map((group) => ({ ...group, bullets: group.bullets.map(applyDoneOverride) }));
}

export async function createBullet(req: CreateBulletRequest): Promise<Bullet> {
  return request<Bullet>("/bullets", {
    method: "POST",
    body: JSON.stringify(req),
  });
}

export async function markBulletDone(id: string): Promise<void> {
  doneOverrides.add(id);
}
