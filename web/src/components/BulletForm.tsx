import { useEffect, useRef, useState, type FormEvent } from "react";
import { createBullet } from "../api/bulletsApi";
import { ApiError, type BulletType } from "../types/bullet";
import { TYPE_LABELS, TYPE_MARKS } from "../lib/bulletMarks";

const TYPES: BulletType[] = ["task", "event", "note"];

interface BulletFormProps {
  onCreated: () => void;
}

export function BulletForm({ onCreated }: BulletFormProps) {
  const [title, setTitle] = useState("");
  const [type, setType] = useState<BulletType>("task");
  const [typeOpen, setTypeOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const typeRef = useRef<HTMLDivElement>(null);

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

  return (
    <form className={`scribble ${typeOpen ? "scribble--type-open" : ""}`} onSubmit={handleSubmit}>
      <div className="scribble__type" ref={typeRef}>
        <button
          type="button"
          className="scribble__type-trigger"
          onClick={() => setTypeOpen((v) => !v)}
          disabled={submitting}
          aria-haspopup="listbox"
          aria-expanded={typeOpen}
          aria-label={`Тип записи: ${TYPE_LABELS[type]}`}
        >
          {TYPE_MARKS[type]}
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
      <input
        ref={inputRef}
        className="scribble__input"
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
      <button
        className="scribble__submit"
        type="submit"
        disabled={submitting || title.trim().length === 0}
        aria-label="Добавить запись"
      >
        ↵
      </button>
      {error && <span className="scribble__error">{error}</span>}
    </form>
  );
}
