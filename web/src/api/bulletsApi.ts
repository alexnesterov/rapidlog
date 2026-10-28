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
  const bullets = await request<Bullet[]>("/bullets");

  // Group bullets by day
  const dayGroups = new Map<string, Bullet[]>();
  for (const bullet of bullets) {
    const date = new Date(bullet.createdAt);
    if (isNaN(date.getTime())) {
      continue; // Skip invalid dates
    }
    // Extract YYYY-MM-DD in local timezone
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const dayKey = `${year}-${month}-${day}`;

    if (!dayGroups.has(dayKey)) {
      dayGroups.set(dayKey, []);
    }
    dayGroups.get(dayKey)!.push(bullet);
  }

  // Sort days descending (newest first)
  const sortedDays = Array.from(dayGroups.keys()).sort((a, b) => b.localeCompare(a));

  // Sort bullets within each day by createdAt ascending
  return sortedDays.map((day) => ({
    day,
    bullets: dayGroups.get(day)!.sort((a, b) =>
      new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
    ),
  }));
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

export async function cancelBullet(id: string): Promise<Bullet> {
  return request<Bullet>(`/bullets/${id}/cancel`, {
    method: "POST",
  });
}
