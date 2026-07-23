import type { Bullet } from "../types/bullet";
import { todayIsoDate } from "../lib/date";

function isoDate(daysOffset: number): string {
  const d = new Date();
  d.setDate(d.getDate() + daysOffset);
  return d.toISOString().slice(0, 10);
}

function timestamp(daysOffset: number, hoursOffset: number): string {
  const d = new Date();
  d.setDate(d.getDate() + daysOffset);
  d.setHours(d.getHours() + hoursOffset);
  return d.toISOString();
}

const today = todayIsoDate();

export const mockBullets: Bullet[] = [
  // сегодня
  {
    id: crypto.randomUUID(),
    title: "Оплатить хостинг",
    date: today,
    status: "OPEN",
    created_at: timestamp(0, -3),
    updated_at: timestamp(0, -3),
  },
  {
    id: crypto.randomUUID(),
    title: "Накидать план по API-эндпоинтам",
    date: today,
    status: "DONE",
    created_at: timestamp(0, -5),
    updated_at: timestamp(0, -2),
  },
  {
    id: crypto.randomUUID(),
    title: "Прогнать миграции на staging",
    date: today,
    status: "OPEN",
    created_at: timestamp(0, -1),
    updated_at: timestamp(0, -1),
  },
  // вчера
  {
    id: crypto.randomUUID(),
    title: "Сверстать форму логина",
    date: isoDate(-1),
    status: "DONE",
    created_at: timestamp(-1, -6),
    updated_at: timestamp(-1, -4),
  },
  {
    id: crypto.randomUUID(),
    title: "Разобраться с падением тестов в CI",
    date: isoDate(-1),
    status: "DONE",
    created_at: timestamp(-1, -3),
    updated_at: timestamp(-1, -2),
  },
  {
    id: crypto.randomUUID(),
    title: "Написать черновик api-design.md",
    date: isoDate(-1),
    status: "OPEN",
    created_at: timestamp(-1, -1),
    updated_at: timestamp(-1, -1),
  },
  // позавчера
  {
    id: crypto.randomUUID(),
    title: "Поднять чистый Go-модуль проекта",
    date: isoDate(-2),
    status: "DONE",
    created_at: timestamp(-2, -5),
    updated_at: timestamp(-2, -4),
  },
  {
    id: crypto.randomUUID(),
    title: "Накидать структуру каталогов",
    date: isoDate(-2),
    status: "DONE",
    created_at: timestamp(-2, -4),
    updated_at: timestamp(-2, -3),
  },
];
