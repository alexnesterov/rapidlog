import { useState } from "react";
import { applyTheme, readTheme, type Theme } from "./theme";

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(readTheme);

  const setTheme = (next: Theme) => {
    applyTheme(next);
    setThemeState(next);
  };

  return { theme, setTheme };
}
