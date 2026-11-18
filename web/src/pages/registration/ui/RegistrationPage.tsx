import { Link, getRouteApi } from "@tanstack/react-router";
import { AuthFlowForm } from "@/features/auth-flow";
import { AuthLayout } from "@/shared/ui/auth-layout";

const route = getRouteApi("/registration");

export function RegistrationPage() {
  const { flow: flowId, return_to: returnTo } = route.useSearch();
  return (
    <AuthLayout
      title="создание журнала"
      footer={
        <>
          Уже есть журнал?{" "}
          <Link to="/login" search={{ return_to: returnTo }}>
            Войти
          </Link>
        </>
      }
    >
      <AuthFlowForm key={flowId} kind="registration" flowId={flowId} returnTo={returnTo} submitLabel="Создать журнал" />
    </AuthLayout>
  );
}
