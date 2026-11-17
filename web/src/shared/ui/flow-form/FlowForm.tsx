import type { FormEvent } from "react";
import type { AuthFlow, AuthFlowKind, FlowField } from "@/shared/api/auth";
import { messageText } from "./messages";
import "./FlowForm.css";

// Подписей над полями нет: текст-подсказка внутри строки, имя поля — для скринридеров.
const names: Record<string, string> = {
  identifier: "Email",
  "traits.email": "Email",
  password: "Пароль",
};

const placeholders: Record<AuthFlowKind, Record<string, string>> = {
  login: { identifier: "email", password: "пароль" },
  registration: { "traits.email": "email", password: "пароль, от 8 символов" },
};

interface FlowFormProps {
  flow: AuthFlow;
  submitLabel: string;
  pending: boolean;
  onSubmit: (values: Record<string, string>) => void;
}

export function FlowForm({ flow, submitLabel, pending, onSubmit }: FlowFormProps) {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const values: Record<string, string> = {};
    new FormData(e.currentTarget).forEach((value, key) => {
      if (typeof value === "string") values[key] = value;
    });
    const submit = flow.fields.find((f) => f.type === "submit");
    if (submit?.value) values[submit.name] = submit.value;
    onSubmit(values);
  };

  return (
    <form className="flow-form" onSubmit={handleSubmit} noValidate>
      {flow.messages.length > 0 && (
        <div className="flow-form__alert" role="alert">
          {flow.messages.map((m) => (
            <p key={m.id}>{messageText(m)}</p>
          ))}
        </div>
      )}

      <div className="flow-form__fields">
        {flow.fields.filter((f) => f.type !== "submit").map((field) => (
          <Field key={field.name} field={field} placeholder={placeholders[flow.kind][field.name]} />
        ))}
      </div>

      <button className="flow-form__submit" type="submit" disabled={pending}>
        {pending ? "…" : submitLabel}
      </button>
    </form>
  );
}

function Field({ field, placeholder }: { field: FlowField; placeholder?: string }) {
  if (field.type === "hidden") {
    return <input type="hidden" name={field.name} value={field.value ?? ""} />;
  }

  const id = `flow-${field.name}`;
  const errorId = `${id}-error`;
  const invalid = field.messages.some((m) => m.type === "error");

  return (
    <div className="flow-form__field">
      <label className="visually-hidden" htmlFor={id}>
        {names[field.name] ?? field.name}
      </label>
      <input
        id={id}
        className="flow-form__input"
        name={field.name}
        type={field.type}
        defaultValue={field.value}
        placeholder={placeholder}
        required={field.required}
        autoComplete={field.autocomplete}
        aria-invalid={invalid || undefined}
        aria-describedby={field.messages.length > 0 ? errorId : undefined}
      />
      {field.messages.length > 0 && (
        <span id={errorId} className={`flow-form__hint flow-form__hint--${field.messages[0].type}`}>
          {field.messages.map(messageText).join(". ")}
        </span>
      )}
    </div>
  );
}
