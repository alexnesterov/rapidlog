import type { ReactNode } from "react";
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
  onPocket: (bullet: Bullet) => void;
  onAddToPocket: (bullet: Bullet) => void;
  animationDelay: number;
  // у сегодняшнего дня — строка кармана, между задачами и формой
  children?: ReactNode;
}

export function DaySection({
  date,
  bullets,
  isToday,
  onComplete,
  onMigrate,
  onCancel,
  onPocket,
  onAddToPocket,
  animationDelay,
  children,
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
        onPocket={onPocket}
      />
      {isToday && children}
      {isToday && <BulletForm onAddToPocket={onAddToPocket} />}
    </section>
  );
}
