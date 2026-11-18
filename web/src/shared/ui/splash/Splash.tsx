import "./Splash.css";

// Пока нечего показать: вход в журнал после логина/регистрации и долгие переходы роутера (>1 с).
// Первую загрузку закрывает копия в index.html.
export function Splash() {
  return (
    <div className="splash">
      <div className="splash__mark" aria-hidden="true">
        <span className="splash__dot" />
      </div>
      <div className="splash__spinner" aria-hidden="true" />
    </div>
  );
}
