// UI-прототип: имитация Kratos в памяти браузера. Заменяется реальными вызовами Kratos
// (шаг 5 плана) с теми же сигнатурами.
import type { AuthFlow, AuthFlowKind, FlowError, FlowMessage, Session, SubmitResult } from "./types";

const STORAGE_KEY = "rapidlog-auth-mock";

interface MockState {
  users: Record<string, { id: string; password: string }>;
  session: Session | null;
}

function load(): MockState {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY);
    if (raw) return JSON.parse(raw) as MockState;
  } catch {
    // sessionStorage недоступен — живём в памяти
  }
  return { users: { "demo@rapidlog.dev": { id: crypto.randomUUID(), password: "password123" } }, session: null };
}

let state = load();

function save() {
  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    // игнорируем
  }
}

const delay = (ms = 350) => new Promise((resolve) => setTimeout(resolve, ms));

const msg = (id: number, text: string): FlowMessage => ({ id, text, type: "error" });

function newFlow(kind: AuthFlowKind, email = ""): AuthFlow {
  return {
    id: crypto.randomUUID(),
    kind,
    messages: [],
    fields: [
      { name: "csrf_token", type: "hidden", value: "mock-csrf", required: true, messages: [] },
      {
        name: kind === "login" ? "identifier" : "traits.email",
        type: "email",
        value: email,
        required: true,
        autocomplete: "email",
        messages: [],
      },
      {
        name: "password",
        type: "password",
        required: true,
        autocomplete: kind === "login" ? "current-password" : "new-password",
        messages: [],
      },
      { name: "method", type: "submit", value: "password", messages: [] },
    ],
  };
}

function withFieldError(flow: AuthFlow, name: string, message: FlowMessage): AuthFlow {
  return {
    ...flow,
    fields: flow.fields.map((f) => (f.name === name ? { ...f, messages: [message] } : f)),
  };
}

export async function getSession(): Promise<Session | null> {
  await delay(150);
  return state.session;
}

export async function logout(): Promise<void> {
  await delay();
  state = { ...state, session: null };
  save();
}

export async function createFlow(kind: AuthFlowKind): Promise<AuthFlow> {
  await delay(200);
  return newFlow(kind);
}

export async function submitFlow(kind: AuthFlowKind, values: Record<string, string>): Promise<SubmitResult> {
  await delay();

  const emailField = kind === "login" ? "identifier" : "traits.email";
  const email = (values[emailField] ?? "").trim().toLowerCase();
  const password = values.password ?? "";
  const flow = newFlow(kind, email);

  if (!email) return { ok: false, flow: withFieldError(flow, emailField, msg(4000002, "Property email is missing.")) };
  if (!password) {
    return { ok: false, flow: withFieldError(flow, "password", msg(4000002, "Property password is missing.")) };
  }

  if (kind === "login") {
    const user = state.users[email];
    if (!user || user.password !== password) {
      return {
        ok: false,
        flow: {
          ...flow,
          messages: [msg(4000006, "The provided credentials are invalid, check for spelling mistakes in your password or username, email address, or phone number.")],
        },
      };
    }
    state = { ...state, session: { identityId: user.id, email } };
    save();
    return { ok: true };
  }

  if (password.length < 8) {
    return {
      ok: false,
      flow: withFieldError(flow, "password", msg(4000032, "The password must be at least 8 characters long, but got 0.")),
    };
  }
  if (state.users[email]) {
    return {
      ok: false,
      flow: { ...flow, messages: [msg(4000007, "An account with the same identifier (email, phone, username, ...) exists already.")] },
    };
  }

  const id = crypto.randomUUID();
  state = { users: { ...state.users, [email]: { id, password } }, session: { identityId: id, email } };
  save();
  return { ok: true };
}

export async function getFlowError(id: string): Promise<FlowError> {
  await delay(150);
  return { message: `Ошибка self-service flow (${id})` };
}
