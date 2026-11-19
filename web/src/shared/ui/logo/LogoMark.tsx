// Знак RapidLog — рукописная точка-задача и быстрый росчерк следом: первая строка любого дня.
// Точка — красная, росчерк — currentColor. pathLength=1 у росчерка — для анимации дописывания.
interface LogoMarkProps {
  className?: string;
}

export function LogoMark({ className }: LogoMarkProps) {
  return (
    <svg className={className ? `logo-mark ${className}` : "logo-mark"} viewBox="0 0 40 16" aria-hidden="true">
      <path
        className="logo-mark__dot"
        d="M5.3 4.3c2.2-.1 3.9 1.5 3.8 3.6-.1 2.2-1.9 3.7-4 3.5-1.9-.1-3.4-1.7-3.2-3.7.1-1.8 1.6-3.3 3.4-3.4z"
      />
      <path
        className="logo-mark__stroke"
        d="M13.2 9.4c1.6-2.3 3.1-2.4 4.2-.4s2.4 2 3.8-.1 2.6-2.3 3.9-.2 2.5 1.9 3.9-.1 2.8-1.6 4.4-.6c1.1.6 2.4.6 3.9.1"
        pathLength={1}
      />
    </svg>
  );
}
