import type { BulletType, MigrateTarget, Signifier } from "../types/bullet";

export const TYPE_MARKS: Record<BulletType, string> = {
  task: "•",
  event: "○",
  note: "–",
};

export const SIGNIFIER_MARKS: Partial<Record<Signifier, string>> = {
  completed: "X",
  migrated: ">",
};

export const TYPE_LABELS: Record<BulletType, string> = {
  task: "задача",
  event: "событие",
  note: "заметка",
};

export const MIGRATE_TARGETS: { id: MigrateTarget; label: string }[] = [
  { id: "today", label: "сегодня" },
];
