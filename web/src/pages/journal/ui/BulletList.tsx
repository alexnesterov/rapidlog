import { useEffect, useLayoutEffect, useRef, useState, type ReactElement } from "react";
import type { Bullet, MigrateTarget } from "../model/bullet";
import { MIGRATE_TARGETS, SIGNIFIER_MARKS, KIND_MARKS } from "./bulletMarks";
import { StrikeLine } from "./bulletIcons";
import { linkifyText } from "@/shared/lib";

function MoreIcon() {
  return (
    <svg className="log-line__menu-icon" viewBox="0 0 16 16" aria-hidden="true">
      <circle cx="3" cy="8" r="1.3" />
      <circle cx="8" cy="8" r="1.3" />
      <circle cx="13" cy="8" r="1.3" />
    </svg>
  );
}

interface LineAction {
  key: string;
  label: string;
  danger?: boolean;
  onSelect: () => void;
}

interface LineStrike {
  left: number;
  top: number;
  width: number;
}

// Open/completed rows: a mark (icon or complete-button) plus title.
// The mark is a plain flex sibling of the title, aligned to the top of the
// row via .log-line__entry's `align-items: flex-start` — which is the first
// line's top for any number of wrapped lines, no measurement needed.
function OpenEntry({
  content,
  Mark,
  interactive,
  onMarkClick,
  markAriaLabel,
}: {
  content: string;
  Mark: () => ReactElement;
  interactive: boolean;
  onMarkClick?: () => void;
  markAriaLabel?: string;
}) {
  return (
    <span className="log-line__entry">
      {interactive ? (
        <button type="button" className="log-line__mark" onClick={onMarkClick} aria-label={markAriaLabel}>
          <Mark />
        </button>
      ) : (
        <span className="log-line__mark" aria-hidden="true">
          <Mark />
        </span>
      )}
      <span className="log-line__title">{linkifyText(content)}</span>
    </span>
  );
}

// Cancelled rows additionally need a strike drawn over each wrapped line of
// the title (a single absolutely-positioned line can't follow wrapped
// text), so — unlike OpenEntry — this one does need to measure the title's
// line boxes via Range.getClientRects().
function CancelledEntry({ content, KindMark }: { content: string; KindMark: () => ReactElement }) {
  const entryRef = useRef<HTMLSpanElement>(null);
  const textRef = useRef<HTMLSpanElement>(null);
  const [strikes, setStrikes] = useState<LineStrike[]>([]);

  useLayoutEffect(() => {
    const entry = entryRef.current;
    const text = textRef.current;
    if (!entry || !text) return;

    function measure() {
      if (!entry || !text) return;
      const entryRect = entry.getBoundingClientRect();

      const range = document.createRange();
      range.selectNodeContents(text);
      const lineRects = Array.from(range.getClientRects());
      if (lineRects.length === 0) return;

      setStrikes(
        lineRects.map((rect, i) => {
          // First line starts at the signifier mark; wrapped continuation
          // lines start where the text itself resumes.
          const left = i === 0 ? 0 : rect.left - entryRect.left;
          const right = rect.right - entryRect.left;
          return {
            left: left - 3,
            top: rect.top - entryRect.top + rect.height / 2,
            width: right - left + 6,
          };
        }),
      );
    }

    measure();
    window.addEventListener("resize", measure);
    const observer = new ResizeObserver(measure);
    observer.observe(entry);
    return () => {
      window.removeEventListener("resize", measure);
      observer.disconnect();
    };
  }, [content]);

  return (
    <span className="log-line__entry" ref={entryRef}>
      <span className="log-line__mark" aria-hidden="true">
        <KindMark />
      </span>
      <span className="log-line__title" ref={textRef}>
        {linkifyText(content)}
      </span>
      {strikes.map((strike, i) => (
        <StrikeLine key={i} style={{ left: strike.left, top: strike.top, width: strike.width }} />
      ))}
    </span>
  );
}

interface BulletListProps {
  bullets: Bullet[];
  canMigrate: boolean;
  onComplete: (bullet: Bullet) => void;
  onMigrate: (bullet: Bullet, target: MigrateTarget) => void;
  onCancel: (bullet: Bullet) => void;
}

export function BulletList({ bullets, canMigrate, onComplete, onMigrate, onCancel }: BulletListProps) {
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
        // only cancelled rows are dimmed (under their strike); completed ones
        // stay full ink — the cross mark alone says they're done
        const cancelled = bullet.signifier === "SIGNIFIER_CANCELLED";
        const canComplete = bullet.kind === "KIND_TASK" && bullet.signifier === "SIGNIFIER_OPEN";
        const showMigrate = canMigrate && canComplete;
        const canCancel = bullet.signifier === "SIGNIFIER_OPEN";
        const pickerOpen = openId === bullet.id;
        const ClosedMark = SIGNIFIER_MARKS[bullet.signifier];
        const KindMark = KIND_MARKS[bullet.kind];

        const actions: LineAction[] = [];
        if (showMigrate) {
          for (const target of MIGRATE_TARGETS) {
            actions.push({
              key: `migrate-${target.id}`,
              label: target.label,
              onSelect: () => onMigrate(bullet, target.id),
            });
          }
        }
        if (canCancel) {
          actions.push({ key: "cancel", label: "Отменить", danger: true, onSelect: () => onCancel(bullet) });
        }

        return (
          <li className={`log-line ${cancelled ? "is-cancelled" : ""} ${pickerOpen ? "is-picking" : ""}`} key={bullet.id}>
            {cancelled ? (
              <CancelledEntry content={bullet.content} KindMark={KindMark} />
            ) : (
              <OpenEntry
                content={bullet.content}
                Mark={ClosedMark ?? KindMark}
                interactive={canComplete}
                onMarkClick={() => onComplete(bullet)}
                markAriaLabel="Отметить выполненным"
              />
            )}
            <div className="log-line__menu" ref={pickerOpen ? pickerRef : undefined}>
              {actions.length > 0 && (
                <button
                  type="button"
                  className="log-line__menu-trigger"
                  onClick={() => setOpenId(pickerOpen ? null : bullet.id)}
                  aria-haspopup="menu"
                  aria-expanded={pickerOpen}
                  aria-label="Действия с записью"
                  title="Действия с записью"
                >
                  <MoreIcon />
                </button>
              )}
              {pickerOpen && (
                <ul className="log-line__menu-list" role="menu" aria-label="Действия с записью">
                  {actions.map((action) => (
                    <li
                      key={action.key}
                      role="menuitem"
                      className={`log-line__menu-option ${action.danger ? "log-line__menu-option--danger" : ""}`}
                      onClick={() => {
                        setOpenId(null);
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
  );
}
