import { LogoMark } from "./LogoMark";
import "./Logo.css";

interface LogoProps {
  // compact — колонтитул журнала, full — страницы входа, где логотип главный.
  size?: "compact" | "full";
}

// Знак стоит перед словом, как bullet перед записью.
// Пометка «preview» — статус продукта: пока он сырой. Убрать, когда выйдет из превью.
export function Logo({ size = "full" }: LogoProps) {
  return (
    <span className={`logo logo--${size}`}>
      <LogoMark className="logo__mark" />
      <span className="logo__word">rapidlog</span>
      <span className="logo__note">preview</span>
    </span>
  );
}
