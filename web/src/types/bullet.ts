export type BulletType = "task" | "event" | "note";

export type Signifier = "open" | "completed" | "migrated" | "scheduled" | "cancelled";

export type MigrateTarget = "today";

export interface Bullet {
  id: string;
  type: BulletType;
  signifier: Signifier;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CreateBulletRequest {
  content: string;
  type?: BulletType;
}

// GET /api/bullets отдаёт bullets уже сгруппированными по дню — новый день
// первым, bullets внутри дня по возрастанию created_at
export interface BulletDayGroup {
  day: string; // YYYY-MM-DD
  bullets: Bullet[];
}

export interface ApiEnvelope<T> {
  data: T;
  error?: {
    code: number;
    message: string;
  };
}

export class ApiError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}
