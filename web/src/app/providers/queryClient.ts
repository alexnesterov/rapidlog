import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";
import { ApiError } from "@/shared/api";

const isUnauthorized = (err: unknown) => err instanceof ApiError && err.code === "401";

// onUnauthorized — /api ответил 401: сессия Kratos истекла или вышли в соседней вкладке.
export function createQueryClient(onUnauthorized: () => void) {
  const handleError = (err: unknown) => {
    if (isUnauthorized(err)) onUnauthorized();
  };

  return new QueryClient({
    queryCache: new QueryCache({ onError: handleError }),
    mutationCache: new MutationCache({ onError: handleError }),
    defaultOptions: {
      queries: {
        // данные из loader'а маршрута не перезапрашиваются сразу при монтировании страницы
        // (с дефолтным 0 был второй GET /api/bullets); мутации обновляют кэш через setQueryData
        staleTime: 60_000,
        // 401 повтором не лечится — сразу отдаём ошибку, чтобы увести на вход
        retry: (failureCount, err) => !isUnauthorized(err) && failureCount < 3,
      },
    },
  });
}
