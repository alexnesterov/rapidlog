import { useState, type FormEvent } from "react";
import { createBullet } from "../api/bulletsApi";
import { ApiError } from "../types/bullet";

interface BulletFormProps {
  onCreated: () => void;
}

export function BulletForm({ onCreated }: BulletFormProps) {
  const [title, setTitle] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (title.trim().length === 0) {
      setError("title is required");
      return;
    }

    setSubmitting(true);
    try {
      await createBullet({ title: title.trim() });
      setTitle("");
      onCreated();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "не удалось создать запись");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className="bullet-form" onSubmit={handleSubmit}>
      <div className="bullet-form__row">
        <input
          className="bullet-form__title"
          type="text"
          placeholder="Новая запись…"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          maxLength={200}
          disabled={submitting}
          aria-label="Название"
        />
        <button className="bullet-form__submit" type="submit" disabled={submitting}>
          {submitting ? "..." : "Добавить"}
        </button>
      </div>
      {error && <p className="bullet-form__error">{error}</p>}
    </form>
  );
}
