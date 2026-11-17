import { useQuery } from "@tanstack/react-query";
import { Link, getRouteApi } from "@tanstack/react-router";
import { AuthLayout } from "@/shared/ui/auth-layout";
import { flowErrorQuery } from "../api/flowErrorQuery";

const route = getRouteApi("/error");

export function AuthErrorPage() {
  const { id: errorId } = route.useSearch();
  const { data, isPending } = useQuery({ ...flowErrorQuery(errorId ?? ""), enabled: Boolean(errorId) });

  const message = !errorId ? "Не удалось завершить вход. Попробуйте ещё раз." : isPending ? "загрузка…" : (data?.message ?? "Не удалось завершить вход. Попробуйте ещё раз.");

  return (
    <AuthLayout
      title="ошибка входа"
      footer={
        <Link to="/login" search={{}}>
          Вернуться ко входу
        </Link>
      }
    >
      <p className="flow-form__state flow-form__state--error">{message}</p>
    </AuthLayout>
  );
}
