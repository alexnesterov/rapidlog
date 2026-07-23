export type BulletStatus = "OPEN" | "DONE";

export interface Bullet {
  id: string;
  title: string;
  date: string; // YYYY-MM-DD
  status: BulletStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateBulletRequest {
  title: string;
  date?: string;
}

export interface ListBulletsParams {
  search?: string;
  limit?: number;
  offset?: number;
}

export interface ListBulletsResponse {
  bullets: Bullet[];
  total: number;
  limit: number;
  offset: number;
}

export interface ApiErrorBody {
  error: {
    code: string;
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
