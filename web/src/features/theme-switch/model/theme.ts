// Тема — атрибут data-theme на <html>. «Paper» — без атрибута, «Catppuccin» — Latte или Mocha;
// обе светлые или тёмные по prefers-color-scheme.
// Ключ и список тем дублируются inline-скриптом в index.html (он ставит тему до отрисовки).
export const THEMES = [
  { id: "paper", label: "Paper" },
  { id: "catppuccin", label: "Catppuccin" },
  { id: "latte", label: "Latte" },
  { id: "frappe", label: "Frappé" },
  { id: "macchiato", label: "Macchiato" },
  { id: "mocha", label: "Mocha" },
] as const;

export type Theme = (typeof THEMES)[number]["id"];

const STORAGE_KEY = "rapidlog-theme";

function isTheme(value: string | undefined): value is Theme {
  return THEMES.some((t) => t.id === value);
}

export function readTheme(): Theme {
  const value = document.documentElement.dataset.theme;
  return isTheme(value) ? value : "paper";
}

// Выбор живёт только в этом браузере; недоступный localStorage (приватный режим) не мешает сменить тему.
export function applyTheme(theme: Theme) {
  const root = document.documentElement;
  try {
    if (theme === "paper") {
      delete root.dataset.theme;
      localStorage.removeItem(STORAGE_KEY);
    } else {
      root.dataset.theme = theme;
      localStorage.setItem(STORAGE_KEY, theme);
    }
  } catch {
    // тема уже применена к <html> выше, не сохранилась только на следующий раз
  }
}
