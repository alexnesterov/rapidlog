export type Kind =
  | 'KIND_TASK'
  | 'KIND_EVENT'
  | 'KIND_NOTE';

export type Signifier =
  | 'SIGNIFIER_OPEN'
  | 'SIGNIFIER_COMPLETED'
  | 'SIGNIFIER_CANCELLED';

export type MigrateTarget = 'today';

export interface Bullet {
  id: string;
  kind: Kind;
  signifier: Signifier;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateBulletRequest {
  content: string;
  kind?: Kind;
}

// GET /api/bullets отдаёт плоский список; группировку по дню делает listBullets
// в api/bulletsApi.ts — новый день первым, bullets внутри дня по возрастанию createdAt
export interface BulletDayGroup {
  day: string; // YYYY-MM-DD
  bullets: Bullet[];
}
