import { useEffect, useRef, useState } from "react";
import type { PocketItem } from "../model/pocket";
import { MoreIcon, PocketMark } from "./bulletIcons";
import { linkifyText } from "@/shared/lib";
import "./PocketSection.css";

function tasksWord(n: number): string {
  const mod10 = n % 10;
  const mod100 = n % 100;
  if (mod10 === 1 && mod100 !== 11) return "задача";
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return "задачи";
  return "задач";
}

interface PocketSectionProps {
  items: PocketItem[];
  onTakeToday: (item: PocketItem) => void;
  onCancel: (item: PocketItem) => void;
}

// Карман — одна тихая строка под сегодняшним днём; раскрывается на месте
// обычным списком записей. Пустой карман не показываем вовсе.
export function PocketSection({ items, onTakeToday, onCancel }: PocketSectionProps) {
  const [open, setOpen] = useState(false);
  const [menuId, setMenuId] = useState<string | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const count = items.length;

  // меню «…» закрывается кликом мимо и Escape — как у строк дня в BulletList
  useEffect(() => {
    if (!menuId) return;

    function onPointerDown(e: PointerEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setMenuId(null);
    }

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setMenuId(null);
    }

    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [menuId]);

  if (count === 0) return null;

  return (
    <div className="pocket">
      <div className="pocket__head">
        <span className="pocket__rule" aria-hidden="true" />
        <button type="button" className="pocket__toggle" onClick={() => setOpen((v) => !v)} aria-expanded={open}>
          в кармане {count} {tasksWord(count)}
          <svg className={`pocket__chevron ${open ? "is-open" : ""}`} viewBox="0 0 16 16" aria-hidden="true">
            <path d="M6 4l4 4-4 4" />
          </svg>
        </button>
      </div>

      {open && (
        <ul className="log-lines pocket__list">
          {items.map((item) => {
            const menuOpen = menuId === item.id;
            const actions = [
              { key: "take", label: "Достать", danger: false, onSelect: () => onTakeToday(item) },
              { key: "cancel", label: "Отменить", danger: true, onSelect: () => onCancel(item) },
            ];

            return (
              <li key={item.id} className={`log-line pocket__line ${menuOpen ? "is-picking" : ""}`}>
                <span className="log-line__entry">
                  <span className="log-line__mark" aria-hidden="true">
                    <PocketMark />
                  </span>
                  <span className="log-line__title">{linkifyText(item.content)}</span>
                </span>
                <div className="log-line__menu" ref={menuOpen ? menuRef : undefined}>
                  <button
                    type="button"
                    className="log-line__menu-trigger"
                    onClick={() => setMenuId(menuOpen ? null : item.id)}
                    aria-haspopup="menu"
                    aria-expanded={menuOpen}
                    aria-label="Действия с задачей в кармане"
                    title="Действия с задачей в кармане"
                  >
                    <MoreIcon />
                  </button>
                  {menuOpen && (
                    <ul className="log-line__menu-list" role="menu" aria-label="Действия с задачей в кармане">
                      {actions.map((action) => (
                        <li
                          key={action.key}
                          role="menuitem"
                          className={`log-line__menu-option ${action.danger ? "log-line__menu-option--danger" : ""}`}
                          onClick={() => {
                            setMenuId(null);
                            action.onSelect();
                          }}
                        >
                          {action.label}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
