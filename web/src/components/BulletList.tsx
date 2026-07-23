import type { Bullet } from "../types/bullet";

interface BulletListProps {
  bullets: Bullet[];
  onToggle: (bullet: Bullet) => void;
}

export function BulletList({ bullets, onToggle }: BulletListProps) {
  if (bullets.length === 0) {
    return <p className="log-state">страница пуста — запиши первую мысль</p>;
  }

  return (
    <ul className="log-lines">
      {bullets.map((bullet) => {
        const done = bullet.status === "DONE";
        return (
          <li className={`log-line ${done ? "is-done" : ""}`} key={bullet.id}>
            <button
              type="button"
              className="log-line__mark"
              onClick={() => onToggle(bullet)}
              aria-label={done ? "Вернуть в открытые" : "Отметить выполненным"}
            >
              {done ? "X" : "•"}
            </button>
            <span className="log-line__title">{bullet.title}</span>
          </li>
        );
      })}
    </ul>
  );
}
