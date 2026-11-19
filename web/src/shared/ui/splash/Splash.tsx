import { LogoMark } from "@/shared/ui/logo";
import "./Splash.css";

// Пока нечего показать: вход в журнал после логина/регистрации и долгие переходы роутера (>1 с).
// Первую загрузку закрывает копия в index.html. Знак логотипа «дописывается» по кругу.
export function Splash() {
  return (
    <div className="splash">
      <LogoMark className="splash__mark" />
    </div>
  );
}
