import type { AuthFlowKind } from "@/shared/api/auth";
import { FlowForm } from "@/shared/ui/flow-form";
import { useAuthFlow } from "../model/useAuthFlow";

interface AuthFlowFormProps {
  kind: AuthFlowKind;
  flowId?: string;
  returnTo?: string;
  submitLabel: string;
}

export function AuthFlowForm({ kind, flowId, returnTo, submitLabel }: AuthFlowFormProps) {
  const { flow, submit } = useAuthFlow(kind, { flowId, returnTo });

  return (
    <FlowForm
      flow={flow}
      submitLabel={submitLabel}
      pending={submit.isPending}
      onSubmit={(values) => submit.mutate(values)}
    />
  );
}
