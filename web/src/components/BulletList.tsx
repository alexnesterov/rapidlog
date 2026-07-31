import type { Bullet } from "../types/bullet";
import { TYPE_MARKS } from "../lib/bulletMarks";

interface BulletListProps {
  bullets: Bullet[];
  onComplete: (bullet: Bullet) => void;
}

export function BulletList({ bullets, onComplete }: BulletListProps) {
  if (bullets.length === 0) {
    return <p className="log-state">страница пуста — запиши первую мысль</p>;
  }

  return (
    <ul className="log-lines">
      {bullets.map((bullet) => {
        const done = bullet.signifier === "completed";
        const canComplete = bullet.type === "task" && !done;
        return (
          <li className={`log-line ${done ? "is-done" : ""}`} key={bullet.id}>
            {done ? (
              <span className="log-line__mark" aria-hidden="true">
                X
              </span>
            ) : canComplete ? (
              <button
                type="button"
                className="log-line__mark"
                onClick={() => onComplete(bullet)}
                aria-label="Отметить выполненным"
              >
                {TYPE_MARKS[bullet.type]}
              </button>
            ) : (
              <span className="log-line__mark" aria-hidden="true">
                {TYPE_MARKS[bullet.type]}
              </span>
            )}
            <span className="log-line__title">{bullet.content}</span>
          </li>
        );
      })}
    </ul>
  );
}
