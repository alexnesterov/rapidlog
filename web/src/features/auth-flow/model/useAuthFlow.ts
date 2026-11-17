import { useState } from "react";
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { useRouter } from "@tanstack/react-router";
import { submitFlow, type AuthFlowKind } from "@/shared/api/auth";
import { authFlowQuery } from "../api/authFlowQuery";
import { safeReturnPath } from "../lib/returnTo";

interface Options {
  flowId?: string;
  returnTo?: string;
}

export function useAuthFlow(kind: AuthFlowKind, { flowId, returnTo }: Options) {
  const queryClient = useQueryClient();
  const router = useRouter();
  const query = authFlowQuery(kind, flowId);
  const { data: flow } = useSuspenseQuery(query);
  // вход/регистрация прошли, грузим журнал — на это время форму закрывает сплэш
  const [entering, setEntering] = useState(false);

  const submit = useMutation({
    mutationFn: (values: Record<string, string>) => submitFlow(kind, flow.id, values),
    onSuccess: async (result) => {
      if (!result.ok) {
        queryClient.setQueryData(query.queryKey, result.flow);
        return;
      }

      // Новая сессия — кэш прошлого юзера не нужен; flow текущей формы оставляем, чтобы она не мигнула.
      queryClient.removeQueries({ predicate: (q) => q.queryKey[0] !== "auth-flow" });

      // Под сплэшем грузим журнал (сессия, записи, шрифты) и только потом переходим.
      // finally: если переход не увёл со страницы входа, сплэш не должен остаться висеть
      // (после успешного перехода страница уже размонтирована и setState ничего не делает).
      setEntering(true);
      try {
        const to = safeReturnPath(returnTo);
        // return_to — произвольный локальный путь, не типизированный маршрут; в рантайме preloadRoute href понимает
        const preload = { href: to } as Parameters<typeof router.preloadRoute>[0];
        await router.preloadRoute(preload).catch(() => undefined);
        await router.navigate({ href: to });
      } finally {
        setEntering(false);
      }
    },
  });

  return { flow, submit, entering };
}
