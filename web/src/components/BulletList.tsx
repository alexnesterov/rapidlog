import { useEffect, useLayoutEffect, useRef, useState, type ReactElement } from "react";
import type { Bullet, MigrateTarget } from "../types/bullet";
import { MIGRATE_TARGETS, SIGNIFIER_MARKS, TYPE_MARKS } from "../lib/bulletMarks";
import { StrikeLine } from "../lib/bulletIcons";
import { linkifyText } from "../lib/linkify";

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

type MarkElement = HTMLSpanElement | HTMLButtonElement;

// Bullet content wraps to an unknown number of visual lines. The mark stays
// a normal (in-flow) flex sibling of the title — same as every other row —
// so the row's height/baseline math is identical across rows and doesn't
// shift neighbouring rows. Flex (`align-self: center` on .log-line__mark)
// centers it on the *whole* (multi-line) entry though, so it's nudged back
// up to the first line via a paint-only translateY, which is a no-op for
// single-line content (first line === whole entry).
//
// onMeasure gets the wrapped line rects for callers (CancelledEntry) that
// need to draw something per line.
function useAlignMarkToFirstLine(content: string, onMeasure?: (lineRects: DOMRect[], entryRect: DOMRect) => void) {
  const entryRef = useRef<HTMLSpanElement>(null);
  const markRef = useRef<MarkElement>(null);
  const textRef = useRef<HTMLSpanElement>(null);
  // Read via a ref inside the effect below instead of listing onMeasure as a
  // dependency, so re-measuring stays tied to content changing, not to the
  // caller re-creating its callback every render.
  const onMeasureRef = useRef(onMeasure);
  onMeasureRef.current = onMeasure;

  useLayoutEffect(() => {
    const entry = entryRef.current;
    const mark = markRef.current;
    const text = textRef.current;
    if (!entry || !mark || !text) return;

    function measure() {
      if (!entry || !mark || !text) return;
      const entryRect = entry.getBoundingClientRect();
      const markRect = mark.getBoundingClientRect();

      const range = document.createRange();
      range.selectNodeContents(text);
      const lineRects = Array.from(range.getClientRects());
      if (lineRects.length === 0) return;

      const firstLine = lineRects[0];
      const firstLineCenter = firstLine.top + firstLine.height / 2;
      const markCenter = markRect.top + markRect.height / 2;
      mark.style.transform = `translateY(${firstLineCenter - markCenter}px)`;

      onMeasureRef.current?.(lineRects, entryRect);
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

  return { entryRef, markRef, textRef };
}

// Open/completed/migrated rows: a mark (icon or complete-button) plus title,
// with the mark aligned to the title's first line regardless of how many
// lines the content wraps to (see useAlignMarkToFirstLine).
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
  const { entryRef, markRef, textRef } = useAlignMarkToFirstLine(content);
  const setMarkRef = (el: MarkElement | null) => {
    markRef.current = el;
  };

  return (
    <span className="log-line__entry" ref={entryRef}>
      {interactive ? (
        <button type="button" className="log-line__mark" ref={setMarkRef} onClick={onMarkClick} aria-label={markAriaLabel}>
          <Mark />
        </button>
      ) : (
        <span className="log-line__mark" aria-hidden="true" ref={setMarkRef}>
          <Mark />
        </span>
      )}
      <span className="log-line__title" ref={textRef}>
        {linkifyText(content)}
      </span>
    </span>
  );
}

function CancelledEntry({ content, TypeMark }: { content: string; TypeMark: () => ReactElement }) {
  const [strikes, setStrikes] = useState<LineStrike[]>([]);

  const { entryRef, markRef, textRef } = useAlignMarkToFirstLine(content, (lineRects, entryRect) => {
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
  });

  return (
    <span className="log-line__entry" ref={entryRef}>
      <span className="log-line__mark" aria-hidden="true" ref={markRef}>
        <TypeMark />
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
        const closed = bullet.signifier === "completed" || bullet.signifier === "migrated" || bullet.signifier === "cancelled";
        const canComplete = bullet.type === "task" && bullet.signifier === "open";
        const showMigrate = canMigrate && canComplete;
        const canCancel = bullet.signifier === "open";
        const pickerOpen = openId === bullet.id;
        const ClosedMark = SIGNIFIER_MARKS[bullet.signifier];
        const TypeMark = TYPE_MARKS[bullet.type];

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
          <li className={`log-line ${closed ? "is-closed" : ""} ${pickerOpen ? "is-picking" : ""}`} key={bullet.id}>
            {bullet.signifier === "cancelled" ? (
              <CancelledEntry content={bullet.content} TypeMark={TypeMark} />
            ) : (
              <OpenEntry
                content={bullet.content}
                Mark={ClosedMark ?? TypeMark}
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
