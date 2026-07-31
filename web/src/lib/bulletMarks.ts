import type { BulletType } from "../types/bullet";

export const TYPE_MARKS: Record<BulletType, string> = {
  task: "•",
  event: "○",
  note: "–",
};

export const TYPE_LABELS: Record<BulletType, string> = {
  task: "задача",
  event: "событие",
  note: "заметка",
};
