import type { Bullet } from "../types/bullet";
import { dayNumber, month, weekday } from "../lib/date";
import { pluralSuffix } from "../lib/plural";
import { BulletForm } from "./BulletForm";
import { BulletList } from "./BulletList";

interface DaySectionProps {
  date: string;
  bullets: Bullet[];
  isToday: boolean;
  onCreated: () => void;
  animationDelay: number;
}

export function DaySection({ date, bullets, isToday, onCreated, animationDelay }: DaySectionProps) {
  const doneCount = bullets.filter((b) => b.status === "DONE").length;

  return (
    <section className="day" style={{ animationDelay: `${animationDelay}ms` }}>
      <header className="day__head">
        <span className="day__head-month">{month(date)}</span>
        <span className="day__head-day">{dayNumber(date)}</span>
        <span className="day__head-weekday">{weekday(date)}</span>
        {bullets.length > 0 && (
          <span className="day__head-stat">
            {bullets.length} запис{pluralSuffix(bullets.length)} · {doneCount} готово
          </span>
        )}
      </header>
      {isToday && <BulletForm onCreated={onCreated} />}
      <BulletList bullets={bullets} />
    </section>
  );
}
