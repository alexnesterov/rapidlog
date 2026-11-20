import { LogoMark } from "@/shared/ui/logo";
import { THEMES, type Theme } from "../model/theme";
import "./ThemePicker.css";

interface ThemePickerProps {
  theme: Theme;
  onChange: (theme: Theme) => void;
}

// Образец темы — знак логотипа в её цветах. Нативные radio: стрелки и Tab работают сами.
export function ThemePicker({ theme, onChange }: ThemePickerProps) {
  const current = THEMES.find((t) => t.id === theme) ?? THEMES[0];

  return (
    <fieldset className="theme-picker">
      <legend className="theme-picker__legend">Тема</legend>
      <div className="theme-picker__row">
        {THEMES.map((t) => (
          <label key={t.id} className={`theme-picker__option theme-picker__option--${t.id}`} title={t.label}>
            <input
              className="theme-picker__input"
              type="radio"
              name="theme"
              value={t.id}
              checked={theme === t.id}
              onChange={() => onChange(t.id)}
              aria-label={t.label}
            />
            <span className="theme-picker__swatch" aria-hidden="true">
              <LogoMark className="theme-picker__mark" />
            </span>
          </label>
        ))}
      </div>
      <p className="theme-picker__current" aria-hidden="true">
        {current.label}
      </p>
    </fieldset>
  );
}
