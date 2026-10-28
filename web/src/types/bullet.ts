export type BulletType = "BULLET_TYPE_TASK" | "BULLET_TYPE_EVENT" | "BULLET_TYPE_NOTE";

export type Signifier = "open" | "completed" | "migrated" | "scheduled" | "cancelled";

export type MigrateTarget = "today";

export interface Bullet {
  id: string;
  type: BulletType;
  signifier: Signifier;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateBulletRequest {
  content: string;
  type?: BulletType;
}

// GET /api/bullets отдаёт bullets уже сгруппированными по дню — новый день
// первым, bullets внутри дня по возрастанию createdAt
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
