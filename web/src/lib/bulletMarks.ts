import type { ReactElement } from "react";
import type { BulletType, MigrateTarget, Signifier } from "../types/bullet";
import { CompletedMark, EventMark, MigratedMark, NoteMark, TaskMark } from "./bulletIcons";

export const TYPE_MARKS: Record<BulletType, () => ReactElement> = {
  task: TaskMark,
  event: EventMark,
  note: NoteMark,
};

export const SIGNIFIER_MARKS: Partial<Record<Signifier, () => ReactElement>> = {
  completed: CompletedMark,
  migrated: MigratedMark,
};

export const TYPE_LABELS: Record<BulletType, string> = {
  task: "задача",
  event: "событие",
  note: "заметка",
};

export const MIGRATE_TARGETS: { id: MigrateTarget; label: string }[] = [
  { id: "today", label: "сегодня" },
];
