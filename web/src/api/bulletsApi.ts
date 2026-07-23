import { mockBullets } from "./mockData";
import { todayIsoDate } from "../lib/date";
import {
  ApiError,
  type Bullet,
  type CreateBulletRequest,
  type ListBulletsParams,
  type ListBulletsResponse,
} from "../types/bullet";

const MOCK_LATENCY_MS = 300;

function delay<T>(value: T): Promise<T> {
  return new Promise((resolve) => setTimeout(() => resolve(value), MOCK_LATENCY_MS));
}

function validateTitle(title: string): void {
  if (title.trim().length === 0) {
    throw new ApiError("validation_error", "title is required");
  }
  if (Array.from(title).length > 200) {
    throw new ApiError("validation_error", "title is too long");
  }
}

export async function listBullets(
  params: ListBulletsParams = {},
): Promise<ListBulletsResponse> {
  const { search, limit = 20, offset = 0 } = params;

  let filtered = mockBullets;
  if (search) {
    const needle = search.toLowerCase();
    filtered = filtered.filter((b) => b.title.toLowerCase().includes(needle));
  }

  const sorted = [...filtered].sort((a, b) => b.created_at.localeCompare(a.created_at));
  const page = sorted.slice(offset, offset + limit);

  return delay({
    bullets: page,
    total: filtered.length,
    limit,
    offset,
  });
}

export async function createBullet(req: CreateBulletRequest): Promise<Bullet> {
  validateTitle(req.title);

  const now = new Date().toISOString();
  const bullet: Bullet = {
    id: crypto.randomUUID(),
    title: req.title,
    date: req.date ?? todayIsoDate(),
    status: "OPEN",
    created_at: now,
    updated_at: now,
  };

  mockBullets.push(bullet);
  return delay(bullet);
}

export async function markBulletDone(id: string): Promise<Bullet> {
  const bullet = mockBullets.find((b) => b.id === id);
  if (!bullet) {
    throw new ApiError("not_found", "bullet not found");
  }

  if (bullet.status !== "DONE") {
    bullet.status = "DONE";
    bullet.updated_at = new Date().toISOString();
  }

  return delay(bullet);
}
