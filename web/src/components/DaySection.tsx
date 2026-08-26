import type { Bullet, MigrateTarget } from "../types/bullet";
import { dayNumber, month, weekday } from "../lib/date";
import { BulletForm } from "./BulletForm";
import { BulletList } from "./BulletList";

interface DaySectionProps {
  date: string;
  bullets: Bullet[];
  isToday: boolean;
  onCreated: () => void;
  onComplete: (bullet: Bullet) => void;
  onMigrate: (bullet: Bullet, target: MigrateTarget) => void;
  animationDelay: number;
}

export function DaySection({
  date,
  bullets,
  isToday,
  onCreated,
  onComplete,
  onMigrate,
  animationDelay,
}: DaySectionProps) {
  return (
    <section className="day" style={{ animationDelay: `${animationDelay}ms` }}>
      <header className="day__head">
        <span className="day__head-month">{month(date)}</span>
        <span className="day__head-day">{dayNumber(date)}</span>
        <span className="day__head-weekday">{weekday(date)}</span>
      </header>
      <BulletList bullets={bullets} canMigrate={!isToday} onComplete={onComplete} onMigrate={onMigrate} />
      {isToday && <BulletForm onCreated={onCreated} />}
    </section>
  );
}
