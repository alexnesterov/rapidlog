import { queryOptions } from "@tanstack/react-query";
import { createFlow, type AuthFlowKind } from "@/shared/api/auth";

// Flow грузит loader маршрута (и preload по наведению на ссылку), страница получает его готовым.
// gcTime > 0: между loader'ом и монтированием у запроса нет подписчиков.
export const authFlowQuery = (kind: AuthFlowKind, flowId?: string) =>
  queryOptions({
    queryKey: ["auth-flow", kind, flowId ?? "new"],
    queryFn: () => createFlow(kind),
    staleTime: Infinity,
    gcTime: 30_000,
    retry: false,
  });
