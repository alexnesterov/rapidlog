import { Link, getRouteApi } from "@tanstack/react-router";
import { AuthFlowForm } from "@/features/auth-flow";
import { AuthLayout } from "@/shared/ui/auth-layout";

const route = getRouteApi("/login");

export function LoginPage() {
  const { flow: flowId, return_to: returnTo } = route.useSearch();
  return (
    <AuthLayout
      title="вход"
      footer={
        <>
          Нет журнала?{" "}
          <Link to="/registration" search={{ return_to: returnTo }}>
            Создать
          </Link>
        </>
      }
    >
      <AuthFlowForm key={flowId} kind="login" flowId={flowId} returnTo={returnTo} submitLabel="Войти" />
    </AuthLayout>
  );
}
