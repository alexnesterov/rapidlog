import type { Bullet, MigrateTarget } from "../model/bullet";
import { dayMonth, weekday } from "@/shared/lib";
import { BulletForm } from "./BulletForm";
import { BulletList } from "./BulletList";

interface DaySectionProps {
  date: string;
  bullets: Bullet[];
  isToday: boolean;
  onComplete: (bullet: Bullet) => void;
  onMigrate: (bullet: Bullet, target: MigrateTarget) => void;
  onCancel: (bullet: Bullet) => void;
  animationDelay: number;
}

export function DaySection({
  date,
  bullets,
  isToday,
  onComplete,
  onMigrate,
  onCancel,
  animationDelay,
}: DaySectionProps) {
  return (
    <section className="day" style={{ animationDelay: `${animationDelay}ms` }}>
      <header className="day__head">
        <span className="day__head-day">{dayMonth(date)}</span>
        <span className="day__head-weekday">{weekday(date)}</span>
      </header>
      <BulletList
        bullets={bullets}
        canMigrate={!isToday}
        onComplete={onComplete}
        onMigrate={onMigrate}
        onCancel={onCancel}
      />
      {isToday && <BulletForm />}
    </section>
  );
}
