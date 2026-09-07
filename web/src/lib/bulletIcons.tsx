import type { CSSProperties } from "react";

export function TaskMark() {
  return (
    <svg className="mark-icon mark-icon--task" viewBox="0 0 16 16" aria-hidden="true">
      <path d="M8 5.7c1.4-.1 2.6.9 2.5 2.3-.1 1.4-1.2 2.4-2.6 2.3-1.2-.1-2.2-1.1-2.1-2.4.1-1.1 1-2.1 2.2-2.2z" />
    </svg>
  );
}

export function EventMark() {
  return (
    <svg className="mark-icon mark-icon--event" viewBox="0 0 16 16" aria-hidden="true">
      <path d="M7.8 3.95c2.5-.4 4.6 1.3 4.6 3.9 0 2.6-2.2 4.6-4.7 4.4-2.4-.2-4.1-2.1-4-4.5.1-2.2 1.8-3.7 4.1-3.8z" />
    </svg>
  );
}

export function NoteMark() {
  return (
    <svg className="mark-icon mark-icon--note" viewBox="0 0 16 16" aria-hidden="true">
      <path d="M3.2 8.4c2.9-.7 6.1-.7 9.6-.1" />
    </svg>
  );
}

export function CompletedMark() {
  return (
    <svg className="mark-icon mark-icon--completed" viewBox="0 0 16 16" aria-hidden="true">
      <path d="M3.6 3.9c2.8 2.7 5.6 5.5 8.4 8.3" />
      <path d="M12.3 4.1c-2.9 2.7-5.7 5.5-8.5 8.2" />
    </svg>
  );
}

export function MigratedMark() {
  return (
    <svg className="mark-icon mark-icon--migrated" viewBox="0 0 16 16" aria-hidden="true">
      <path d="M5.8 3.5c1.7 1.4 3.2 2.9 4.5 4.6-1.5 1.6-3.1 3.1-4.8 4.4" />
    </svg>
  );
}

interface StrikeLineProps {
  style?: CSSProperties;
}

export function StrikeLine({ style }: StrikeLineProps) {
  return (
    <svg className="log-line__strike" style={style} viewBox="0 0 100 14" preserveAspectRatio="none" aria-hidden="true">
      <path
        className="log-line__strike-stroke"
        d="M2 7.3Q11 6.3 19 7.5T36 7T53 7.3T70 6.9T87 7.2T98 7"
        vectorEffect="non-scaling-stroke"
      />
      <path
        className="log-line__strike-stroke log-line__strike-stroke--accent"
        d="M48 7.9Q62 7 75 7.9T98 7.5"
        vectorEffect="non-scaling-stroke"
      />
    </svg>
  );
}
