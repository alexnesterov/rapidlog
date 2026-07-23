import type { BulletStatus } from "../types/bullet";

interface StatusBadgeProps {
  status: BulletStatus;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const isDone = status === "DONE";
  return (
    <span className={`status-badge ${isDone ? "is-done" : "is-open"}`}>
      <span className="status-badge__mark" aria-hidden="true">
        {isDone ? "✕" : "•"}
      </span>
      {isDone ? "готово" : "открыто"}
    </span>
  );
}
