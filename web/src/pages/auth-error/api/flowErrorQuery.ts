import { queryOptions } from "@tanstack/react-query";
import { getFlowError } from "@/shared/api/auth";

export const flowErrorQuery = (id: string) =>
  queryOptions({
    queryKey: ["flow-error", id],
    queryFn: () => getFlowError(id),
    retry: false,
  });
