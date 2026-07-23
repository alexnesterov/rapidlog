import type { Bullet } from "../types/bullet";
import { dayNumber, month, weekday } from "../lib/date";
import { BulletForm } from "./BulletForm";
import { BulletList } from "./BulletList";

interface DaySectionProps {
  date: string;
  bullets: Bullet[];
  isToday: boolean;
  onCreated: () => void;
  onComplete: (bullet: Bullet) => void;
  animationDelay: number;
}

export function DaySection({ date, bullets, isToday, onCreated, onComplete, animationDelay }: DaySectionProps) {
  return (
    <section className="day" style={{ animationDelay: `${animationDelay}ms` }}>
      <header className="day__head">
        <span className="day__head-month">{month(date)}</span>
        <span className="day__head-day">{dayNumber(date)}</span>
        <span className="day__head-weekday">{weekday(date)}</span>
      </header>
      <BulletList bullets={bullets} onComplete={onComplete} />
      {isToday && <BulletForm onCreated={onCreated} />}
    </section>
  );
}
