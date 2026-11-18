import type { ReactNode } from "react";
import { Logo } from "@/shared/ui/logo";
import "./AuthLayout.css";

interface AuthLayoutProps {
  // Видимого заголовка нет — какая это страница, говорит кнопка; title остаётся для скринридеров.
  title: string;
  children: ReactNode;
  footer?: ReactNode;
}

export function AuthLayout({ title, children, footer }: AuthLayoutProps) {
  return (
    <main className="auth">
      <h1 className="auth__brand">
        <Logo />
        <span className="visually-hidden">: {title}</span>
      </h1>
      {children}
      {footer && <p className="auth__footer">{footer}</p>}
    </main>
  );
}
