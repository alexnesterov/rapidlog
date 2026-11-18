import "./Logo.css";

// Рукописный логотип с маркерным выделением «Log» — один и тот же в журнале и на входе.
// Бейдж «preview» — статус продукта: пока он сырой. Убрать, когда выйдет из превью.
export function Logo() {
  return (
    <span className="brand">
      <span className="logo">
        Rapid<span className="logo__accent">Log</span>
      </span>
      <span className="brand__badge">preview</span>
    </span>
  );
}
