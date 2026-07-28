export type BulletStatus = "OPEN" | "DONE";

// форма, в которой bullet приходит по сети — как в entity.Bullet на бэкенде
export interface RawBullet {
  id: string;
  title: string;
  status: BulletStatus;
  created_at: string;
  updated_at: string;
}

// то же самое плюс date, вычисленное на клиенте из created_at — им живёт UI
export interface Bullet extends RawBullet {
  date: string; // YYYY-MM-DD
}

export interface CreateBulletRequest {
  title: string;
}

export interface ListBulletsResponse {
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
