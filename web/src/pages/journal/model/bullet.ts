export type BulletType = "BULLET_TYPE_TASK" | "BULLET_TYPE_EVENT" | "BULLET_TYPE_NOTE";

export type Signifier =
  | "SIGNIFIER_OPEN"
  | "SIGNIFIER_COMPLETED"
  | "SIGNIFIER_SCHEDULED"
  | "SIGNIFIER_CANCELLED";

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

// GET /api/bullets отдаёт плоский список; группировку по дню делает listBullets
// в api/bulletsApi.ts — новый день первым, bullets внутри дня по возрастанию createdAt
export interface BulletDayGroup {
  day: string; // YYYY-MM-DD
  bullets: Bullet[];
}
