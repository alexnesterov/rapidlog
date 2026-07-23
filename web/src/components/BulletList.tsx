import type { Bullet } from "../types/bullet";
import { StatusBadge } from "./StatusBadge";

interface BulletListProps {
  bullets: Bullet[];
  loading: boolean;
  error: string | null;
}

function formatDate(iso: string): string {
  const d = new Date(iso + "T00:00:00");
  return d.toLocaleDateString("ru-RU", { day: "2-digit", month: "short" });
}

export function BulletList({ bullets, loading, error }: BulletListProps) {
  if (loading) {
    return <p className="bullet-list__state">загружаю…</p>;
  }

  if (error) {
    return <p className="bullet-list__state bullet-list__state--error">{error}</p>;
  }

  if (bullets.length === 0) {
    return <p className="bullet-list__state">пока пусто — добавь первую запись</p>;
  }

  return (
    <ul className="bullet-list">
      {bullets.map((bullet) => (
        <li className={`bullet-item ${bullet.status === "DONE" ? "is-done" : ""}`} key={bullet.id}>
          <span className="bullet-item__title">{bullet.title}</span>
          <span className="bullet-item__meta">
            <time className="bullet-item__date">{formatDate(bullet.date)}</time>
            <StatusBadge status={bullet.status} />
          </span>
        </li>
      ))}
    </ul>
  );
}
