import { THEMES } from "../model/theme";
import { useTheme } from "../model/useTheme";
import { ThemePicker } from "./ThemePicker";
import "./ThemeMenu.css";

// Кнопка темы в шапке, пара к кнопке аккаунта: то же рукописное кольцо, внутри — образец
// текущей темы. Выбор — в нативном popover (Esc и клик мимо закрывают его сами).
export function ThemeMenu() {
  const { theme, setTheme } = useTheme();
  const label = (THEMES.find((t) => t.id === theme) ?? THEMES[0]).label;

  return (
    <>
      <button type="button" className="theme-menu__toggle" popoverTarget="theme-menu" aria-label={`Тема: ${label}`}>
        <svg className="theme-menu__ring" viewBox="0 0 32 32" aria-hidden="true">
          <path d="M15.6 3.6c6.9-.5 12.6 4.5 12.6 11.6 0 7-5.6 12.6-12.6 12.3C9 27.2 4 22 4.1 15.6 4.3 9.3 9.2 4 15.6 3.6z" />
        </svg>
        {/* классы образца из ThemePicker.css — цвета тем описаны там один раз */}
        <span className={`theme-menu__current theme-picker__option--${theme}`} aria-hidden="true">
          <span className="theme-picker__swatch" />
        </span>
      </button>
      <div id="theme-menu" className="theme-menu__popover" popover="auto">
        <ThemePicker theme={theme} onChange={setTheme} />
      </div>
    </>
  );
}
