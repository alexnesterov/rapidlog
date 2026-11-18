// Упрощённая модель self-service flow Kratos: UI не зависит от формата @ory/client-fetch,
// реальные flow маппятся в неё в kratos.ts.

export type AuthFlowKind = "login" | "registration";

export type FlowMessageType = "error" | "info" | "success";

export interface FlowMessage {
  // id сообщения Kratos (https://www.ory.sh/docs/kratos/concepts/ui-messages) — по нему переводим
  id: number;
  text: string;
  type: FlowMessageType;
}

export type FlowFieldType = "email" | "text" | "password" | "hidden" | "submit";

export interface FlowField {
  name: string;
  type: FlowFieldType;
  value?: string;
  required?: boolean;
  autocomplete?: string;
  messages: FlowMessage[];
}

export interface AuthFlow {
  id: string;
  kind: AuthFlowKind;
  fields: FlowField[];
  messages: FlowMessage[];
}

export type SubmitResult = { ok: true } | { ok: false; flow: AuthFlow };

export interface Session {
  identityId: string;
  email: string;
}

export interface FlowError {
  message: string;
}
