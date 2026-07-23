import { useRef, useState, type FormEvent } from "react";
import { createBullet } from "../api/bulletsApi";
import { ApiError } from "../types/bullet";

interface BulletFormProps {
  onCreated: () => void;
}

export function BulletForm({ onCreated }: BulletFormProps) {
  const [title, setTitle] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (title.trim().length === 0) {
      setError("нужно название");
      return;
    }

    setSubmitting(true);
    try {
      await createBullet({ title: title.trim() });
      setTitle("");
      onCreated();
      inputRef.current?.focus();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "не удалось создать запись");
    } finally {
      setSubmitting(false);
    }
  }

  const hasText = title.trim().length > 0;

  return (
    <form className="scribble" onSubmit={handleSubmit}>
      <span className={`scribble__mark ${hasText ? "scribble__mark--active" : ""}`} aria-hidden="true">
        {hasText ? "•" : "+"}
      </span>
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
