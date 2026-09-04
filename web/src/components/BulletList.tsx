import { useEffect, useRef, useState } from "react";
import type { Bullet, MigrateTarget } from "../types/bullet";
import { MIGRATE_TARGETS, SIGNIFIER_MARKS, TYPE_MARKS } from "../lib/bulletMarks";

function MigrateIcon() {
  return (
    <svg className="log-line__migrate-icon" viewBox="0 0 16 16" aria-hidden="true">
      <rect x="2.5" y="2.5" width="11" height="11" rx="1.6" />
      <line x1="2.5" y1="6" x2="13.5" y2="6" />
    </svg>
  );
}

interface BulletListProps {
  bullets: Bullet[];
  canMigrate: boolean;
  onComplete: (bullet: Bullet) => void;
  onMigrate: (bullet: Bullet, target: MigrateTarget) => void;
}

export function BulletList({ bullets, canMigrate, onComplete, onMigrate }: BulletListProps) {
  const [openId, setOpenId] = useState<string | null>(null);
  const pickerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!openId) return;

    function onPointerDown(e: PointerEvent) {
      if (pickerRef.current && !pickerRef.current.contains(e.target as Node)) {
        setOpenId(null);
      }
    }

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setOpenId(null);
    }

    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [openId]);

  if (bullets.length === 0) {
    return <p className="log-state">новый день — с чего начнёшь?</p>;
  }

  return (
    <ul className="log-lines">
      {bullets.map((bullet) => {
        const closed = bullet.signifier === "completed" || bullet.signifier === "migrated";
        const canComplete = bullet.type === "task" && bullet.signifier === "open";
        const showMigrate = canMigrate && canComplete;
        const pickerOpen = openId === bullet.id;
        const ClosedMark = SIGNIFIER_MARKS[bullet.signifier];
        const TypeMark = TYPE_MARKS[bullet.type];

        return (
          <li className={`log-line ${closed ? "is-closed" : ""} ${pickerOpen ? "is-picking" : ""}`} key={bullet.id}>
            {ClosedMark ? (
              <span className="log-line__mark" aria-hidden="true">
                <ClosedMark />
              </span>
            ) : canComplete ? (
              <button
                type="button"
                className="log-line__mark"
                onClick={() => onComplete(bullet)}
                aria-label="Отметить выполненным"
              >
                <TypeMark />
              </button>
            ) : (
              <span className="log-line__mark" aria-hidden="true">
                <TypeMark />
              </span>
            )}
            <span className="log-line__title">{bullet.content}</span>
            {showMigrate && (
              <div className="log-line__migrate" ref={pickerOpen ? pickerRef : undefined}>
                <button
                  type="button"
                  className="log-line__migrate-trigger"
                  onClick={() => setOpenId(pickerOpen ? null : bullet.id)}
                  aria-haspopup="menu"
                  aria-expanded={pickerOpen}
                  aria-label="Перенести задачу"
                >
                  <MigrateIcon />
                </button>
                {pickerOpen && (
                  <ul className="log-line__migrate-menu" role="menu" aria-label="Куда перенести">
                    {MIGRATE_TARGETS.map((target) => (
                      <li
                        key={target.id}
                        role="menuitem"
                        className="log-line__migrate-option"
                        onClick={() => {
                          setOpenId(null);
                          onMigrate(bullet, target.id);
                        }}
                      >
                        {target.label}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            )}
          </li>
        );
      })}
    </ul>
  );
}
