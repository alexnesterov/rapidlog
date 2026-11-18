import { AuthLayout } from "@/shared/ui/auth-layout";

// errorComponent маршрутов входа/регистрации: flow не удалось создать.
export function AuthFlowFailed() {
  return (
    <AuthLayout title="ошибка входа">
      <p className="flow-form__state flow-form__state--error">Не удалось открыть форму. Обновите страницу.</p>
    </AuthLayout>
  );
}
