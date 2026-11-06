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
