import type { ReactElement } from "react";
import type { Kind, MigrateTarget, Signifier } from "../model/bullet";
import { CompletedMark, EventMark, NoteMark, TaskMark } from "./bulletIcons";

export const KIND_MARKS: Record<Kind, () => ReactElement> = {
  KIND_TASK: TaskMark,
  KIND_EVENT: EventMark,
  KIND_NOTE: NoteMark,
};

export const SIGNIFIER_MARKS: Partial<Record<Signifier, () => ReactElement>> = {
  SIGNIFIER_COMPLETED: CompletedMark,
};

export const KIND_LABELS: Record<Kind, string> = {
  KIND_TASK: "задача",
  KIND_EVENT: "событие",
  KIND_NOTE: "заметка",
};

export const MIGRATE_TARGETS: { id: MigrateTarget; label: string }[] = [
  { id: "today", label: "Сегодня" },
];
