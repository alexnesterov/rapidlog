import type { Bullet } from "../types/bullet";

function isoDate(daysOffset: number): string {
  const d = new Date();
  d.setDate(d.getDate() + daysOffset);
  return d.toISOString().slice(0, 10);
}

function timestamp(daysOffset: number): string {
  const d = new Date();
  d.setDate(d.getDate() + daysOffset);
  return d.toISOString();
}

export const mockBullets: Bullet[] = [
  {
    id: crypto.randomUUID(),
    title: "Оплатить хостинг",
    date: isoDate(0),
    status: "OPEN",
    created_at: timestamp(-1),
    updated_at: timestamp(-1),
  },
  {
    id: crypto.randomUUID(),
    title: "Накидать план по API-эндпоинтам",
    date: isoDate(-1),
    status: "DONE",
    created_at: timestamp(-2),
    updated_at: timestamp(-1),
  },
  {
    id: crypto.randomUUID(),
    title: "Прогнать миграции на staging",
    date: isoDate(1),
    status: "OPEN",
    created_at: timestamp(-1),
    updated_at: timestamp(-1),
  },
];
