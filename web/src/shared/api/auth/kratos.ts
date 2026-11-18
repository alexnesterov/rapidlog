// Ory Kratos через same-origin прокси gateway (/auth/*). Browser-flow в JSON-режиме:
// Kratos не редиректит, а отвечает JSON, cookie сессии и CSRF ставит на наш origin.
import {
  Configuration,
  FrontendApi,
  isResponseError,
  isSessionAlreadyAvailable,
  isUiNodeInputAttributes,
  type LoginFlow,
  type RegistrationFlow,
  type UiNode,
  type UiText,
} from "@ory/client-fetch";
import type { AuthFlow, AuthFlowKind, FlowError, FlowField, FlowFieldType, FlowMessage, Session, SubmitResult } from "./types";

const ory = new FrontendApi(
  new Configuration({
    basePath: "/auth",
    credentials: "include",
    // без Accept: application/json Kratos отвечает на browser-flow редиректом, а не JSON
    headers: { Accept: "application/json" },
  }),
);

// группы узлов формы, которые мы рендерим: общий CSRF и метод «пароль»
const RENDERED_GROUPS = new Set(["default", "password"]);

// сообщения о просроченном flow (Kratos их шлёт сам в этих случаях; дублируем при пересоздании)
const FLOW_EXPIRED: Record<AuthFlowKind, FlowMessage> = {
  login: { id: 4010001, text: "The login flow expired", type: "error" },
  registration: { id: 4040001, text: "The registration flow expired", type: "error" },
};

const status = (err: unknown) => (isResponseError(err) ? err.response.status : undefined);

async function responseBody(err: unknown): Promise<unknown> {
  if (!isResponseError(err)) return undefined;
  try {
    return await err.response.clone().json();
  } catch {
    return undefined;
  }
}

function toMessage(text: UiText): FlowMessage {
  return { id: text.id, text: text.text, type: text.type as FlowMessage["type"] };
}

function toFieldType(name: string, type: string): FlowFieldType {
  if (name === "identifier" || name === "traits.email") return "email";
  if (type === "password" || type === "hidden" || type === "submit") return type;
  return "text";
}

function toField(node: UiNode): FlowField | null {
  const attrs = node.attributes;
  if (!isUiNodeInputAttributes(attrs) || !RENDERED_GROUPS.has(node.group)) return null;

  return {
    name: attrs.name,
    type: toFieldType(attrs.name, attrs.type),
    value: attrs.value == null ? undefined : String(attrs.value),
    required: attrs.required,
    autocomplete: attrs.autocomplete,
    messages: node.messages.map(toMessage),
  };
}

function toAuthFlow(kind: AuthFlowKind, flow: LoginFlow | RegistrationFlow): AuthFlow {
  return {
    id: flow.id,
    kind,
    fields: flow.ui.nodes.map(toField).filter((f): f is FlowField => f !== null),
    messages: (flow.ui.messages ?? []).map(toMessage),
  };
}

function isFlow(body: unknown): body is LoginFlow | RegistrationFlow {
  return typeof body === "object" && body !== null && "id" in body && "ui" in body;
}

export async function getSession(): Promise<Session | null> {
  try {
    const session = await ory.toSession();
    const identity = session.identity;
    if (!session.active || !identity) return null;
    return { identityId: identity.id, email: String(identity.traits?.email ?? "") };
  } catch (err) {
    if (status(err) === 401 || status(err) === 403) return null;
    throw err;
  }
}

export async function logout(): Promise<void> {
  try {
    const { logout_token } = await ory.createBrowserLogoutFlow();
    await ory.updateLogoutFlow({ token: logout_token });
  } catch (err) {
    // сессии уже нет — считаем, что вышли
    if (status(err) === 401) return;
    throw err;
  }
}

async function createNewFlow(kind: AuthFlowKind): Promise<AuthFlow> {
  const flow = kind === "login" ? await ory.createBrowserLoginFlow() : await ory.createBrowserRegistrationFlow();
  return toAuthFlow(kind, flow);
}

// flowId — когда Kratos сам отправил браузер на ui_url с ?flow=<id>; просроченный или чужой flow
// заменяем новым.
export async function createFlow(kind: AuthFlowKind, flowId?: string): Promise<AuthFlow> {
  if (flowId) {
    try {
      const flow = kind === "login" ? await ory.getLoginFlow({ id: flowId }) : await ory.getRegistrationFlow({ id: flowId });
      return toAuthFlow(kind, flow);
    } catch (err) {
      const code = status(err);
      if (code !== 403 && code !== 404 && code !== 410) throw err;
    }
  }
  return createNewFlow(kind);
}

export async function submitFlow(kind: AuthFlowKind, flowId: string, values: Record<string, string>): Promise<SubmitResult> {
  const csrf_token = values.csrf_token;
  const password = values.password ?? "";

  try {
    if (kind === "login") {
      await ory.updateLoginFlow({
        flow: flowId,
        updateLoginFlowBody: { method: "password", csrf_token, identifier: values.identifier ?? "", password },
      });
    } else {
      await ory.updateRegistrationFlow({
        flow: flowId,
        updateRegistrationFlowBody: { method: "password", csrf_token, traits: { email: values["traits.email"] ?? "" }, password },
      });
    }
    return { ok: true };
  } catch (err) {
    const body = await responseBody(err);

    // 400 с flow внутри — ошибки валидации/неверный пароль, показываем их в форме
    if (status(err) === 400 && isFlow(body)) return { ok: false, flow: toAuthFlow(kind, body) };
    // уже вошли (например, в соседней вкладке) — для пользователя это успех
    if (isSessionAlreadyAvailable(body)) return { ok: true };
    // flow истёк (410) или сломан CSRF (403) — начинаем заново с понятным сообщением
    if (status(err) === 410 || status(err) === 403) {
      const fresh = await createNewFlow(kind);
      return { ok: false, flow: { ...fresh, messages: [FLOW_EXPIRED[kind]] } };
    }
    throw err;
  }
}

export async function getFlowError(id: string): Promise<FlowError> {
  const res = await ory.getFlowError({ id });
  const error = (res.error ?? {}) as { message?: string; reason?: string };
  return { message: error.reason || error.message || "Не удалось завершить вход. Попробуйте ещё раз." };
}
