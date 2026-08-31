import { useEffect, useRef, useState, type FormEvent } from "react";
import { createBullet } from "../api/bulletsApi";
import { ApiError, type BulletType } from "../types/bullet";
import { TYPE_LABELS, TYPE_MARKS } from "../lib/bulletMarks";

const TYPES: BulletType[] = ["task", "event", "note"];

interface BulletFormProps {
  onCreated: () => void;
}

export function BulletForm({ onCreated }: BulletFormProps) {
  const [expanded, setExpanded] = useState(false);
  const [title, setTitle] = useState("");
  const [type, setType] = useState<BulletType>("task");
  const [typeOpen, setTypeOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const typeRef = useRef<HTMLDivElement>(null);
  const formRef = useRef<HTMLFormElement>(null);

  useEffect(() => {
    if (!typeOpen) return;

    function onPointerDown(e: PointerEvent) {
      if (typeRef.current && !typeRef.current.contains(e.target as Node)) {
        setTypeOpen(false);
      }
    }

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setTypeOpen(false);
    }

    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [typeOpen]);

  useEffect(() => {
    if (!expanded) return;
    inputRef.current?.focus();
  }, [expanded]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (title.trim().length === 0) {
      setError("нужно название");
      return;
    }

    setSubmitting(true);
    try {
      await createBullet({ content: title.trim(), type });
      setTitle("");
      setType("task");
      onCreated();
      inputRef.current?.focus();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "не удалось создать запись");
    } finally {
      setSubmitting(false);
    }
  }

  function handleCancel() {
    setExpanded(false);
    setTitle("");
    setError(null);
  }

  if (!expanded) {
    return (
      <button type="button" className="scribble-trigger" onClick={() => setExpanded(true)}>
        <span className="scribble-trigger__mark">+</span>
        Добавить запись
      </button>
    );
  }

  return (
    <form
      ref={formRef}
      className={`scribble-card ${typeOpen ? "scribble-card--type-open" : ""}`}
      onSubmit={handleSubmit}
      onKeyDown={(e) => {
        if (e.key === "Escape") handleCancel();
      }}
    >
      <input
        ref={inputRef}
        className="scribble-card__input"
        type="text"
        placeholder="записать новое…"
        value={title}
        onChange={(e) => {
          setTitle(e.target.value);
          if (error) setError(null);
        }}
        maxLength={200}
        disabled={submitting}
        aria-label="Название новой записи"
      />
      <div className="scribble-card__footer">
        <div className="scribble-card__type" ref={typeRef}>
          <button
            type="button"
            className="scribble-card__type-trigger"
            onClick={() => setTypeOpen((v) => !v)}
            disabled={submitting}
            aria-haspopup="listbox"
            aria-expanded={typeOpen}
            aria-label={`Тип записи: ${TYPE_LABELS[type]}`}
          >
            <span className="scribble-card__type-mark">{TYPE_MARKS[type]}</span>
            {TYPE_LABELS[type]}
          </button>
          {typeOpen && (
            <ul className="scribble__type-menu" role="listbox" aria-label="Тип записи">
              {TYPES.map((t) => (
                <li
                  key={t}
                  role="option"
                  aria-selected={t === type}
                  className={`scribble__type-option ${t === type ? "scribble__type-option--selected" : ""}`}
                  onClick={() => {
                    setType(t);
                    setTypeOpen(false);
                  }}
                >
                  <span className="scribble__type-option-mark">{TYPE_MARKS[t]}</span>
                  <span className="scribble__type-option-label">{TYPE_LABELS[t]}</span>
                </li>
              ))}
            </ul>
          )}
        </div>
        <div className="scribble-card__actions">
          <button
            type="button"
            className="scribble-card__cancel"
            onClick={handleCancel}
            aria-label="Отменить"
            data-tooltip="Отменить"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M4 4L12 12M12 4L4 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </button>
          <button
            className="scribble-card__submit"
            type="submit"
            disabled={submitting || title.trim().length === 0}
            aria-label="Добавить запись"
            data-tooltip="Добавить запись"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M8 13V3M8 3L3 8M8 3L13 8" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </button>
        </div>
      </div>
      {error && <span className="scribble-card__error">{error}</span>}
    </form>
  );
}
