import "./Splash.css";

// Первый экран, пока нечего показать: холодный старт и долгие переходы роутера.
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
