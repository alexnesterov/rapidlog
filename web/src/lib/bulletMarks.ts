import type { ReactElement } from "react";
import type { BulletType, MigrateTarget, Signifier } from "../types/bullet";
import { CompletedMark, EventMark, MigratedMark, NoteMark, TaskMark } from "./bulletIcons";

export const TYPE_MARKS: Record<BulletType, () => ReactElement> = {
  BULLET_TYPE_TASK: TaskMark,
  BULLET_TYPE_EVENT: EventMark,
  BULLET_TYPE_NOTE: NoteMark,
};

export const SIGNIFIER_MARKS: Partial<Record<Signifier, () => ReactElement>> = {
  completed: CompletedMark,
  migrated: MigratedMark,
};

export const TYPE_LABELS: Record<BulletType, string> = {
  BULLET_TYPE_TASK: "задача",
  BULLET_TYPE_EVENT: "событие",
  BULLET_TYPE_NOTE: "заметка",
};

export const MIGRATE_TARGETS: { id: MigrateTarget; label: string }[] = [
  { id: "today", label: "Сегодня" },
];
