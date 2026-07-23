import type { Bullet } from "../types/bullet";

interface BulletListProps {
  bullets: Bullet[];
}

export function BulletList({ bullets }: BulletListProps) {
  if (bullets.length === 0) {
    return <p className="log-state">страница пуста — запиши первую мысль</p>;
  }

  return (
    <ul className="log-lines">
      {bullets.map((bullet) => {
        const done = bullet.status === "DONE";
        return (
          <li className={`log-line ${done ? "is-done" : ""}`} key={bullet.id}>
            <span className="log-line__mark" aria-hidden="true">
              {done ? "X" : "•"}
            </span>
            <span className="log-line__title">{bullet.title}</span>
          </li>
        );
      })}
    </ul>
  );
}
