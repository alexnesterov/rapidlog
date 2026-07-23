import type { Bullet } from "../types/bullet";

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
        const done = bullet.status === "DONE";
        return (
          <li className={`log-line ${done ? "is-done" : ""}`} key={bullet.id}>
            {done ? (
              <span className="log-line__mark" aria-hidden="true">
                X
              </span>
            ) : (
              <button
                type="button"
                className="log-line__mark"
                onClick={() => onComplete(bullet)}
                aria-label="Отметить выполненным"
              >
                •
              </button>
            )}
            <span className="log-line__title">{bullet.title}</span>
          </li>
        );
      })}
    </ul>
  );
}
