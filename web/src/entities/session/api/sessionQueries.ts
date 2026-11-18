import { queryOptions, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getSession, logout } from "@/shared/api/auth";

export const sessionQuery = queryOptions({
  queryKey: ["session"],
  queryFn: getSession,
  staleTime: 60_000,
});

export function useSession() {
  return useQuery(sessionQuery);
}

// Порядок важен, иначе журнал мигает пустым состоянием перед уходом:
// 1. removeQueries не оповещает подписчиков — журнал продолжает показывать свои данные,
//    а guard страницы входа перечитает сессию и пропустит;
// 2. переход на /login (onLoggedOut);
// 3. остальной кэш прошлого юзера чистим, когда журнал уже размонтирован.
//    Flow формы входа и свежую сессию (null) оставляем — их только что загрузил роутер.
export function useLogout(onLoggedOut: () => Promise<unknown>) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logout,
    onSuccess: async () => {
      queryClient.removeQueries({ queryKey: sessionQuery.queryKey });
      await onLoggedOut();
      queryClient.removeQueries({
        predicate: (q) => q.queryKey[0] !== "auth-flow" && q.queryKey[0] !== sessionQuery.queryKey[0],
      });
    },
  });
}
