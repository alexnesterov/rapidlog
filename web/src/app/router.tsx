import type { QueryClient } from "@tanstack/react-query";
import { Navigate, Outlet, createRootRouteWithContext, createRoute, createRouter, redirect } from "@tanstack/react-router";
import { sessionQuery } from "@/entities/session";
import { AuthFlowFailed, authFlowQuery } from "@/features/auth-flow";
import { AuthErrorPage } from "@/pages/auth-error";
import { JournalPage, bulletsQuery } from "@/pages/journal";
import { LoginPage } from "@/pages/login";
import { RegistrationPage } from "@/pages/registration";
import { waitForFonts } from "@/shared/lib";
import { Splash } from "@/shared/ui/splash";

interface RouterContext {
  queryClient: QueryClient;
}

interface AuthSearch {
  flow?: string;
  return_to?: string;
}

const optionalString = (value: unknown) => (typeof value === "string" && value !== "" ? value : undefined);

const validateAuthSearch = (search: Record<string, unknown>): AuthSearch => ({
  flow: optionalString(search.flow),
  return_to: optionalString(search.return_to),
});

const rootRoute = createRootRouteWithContext<RouterContext>()({
  // любая страница впервые показывается уже с Caveat/Golos, без подмены системного шрифта
  // (шрифты Google подключены с display=swap); пока ждём — сплэш из index.html. Для уже загруженных
  // шрифтов document.fonts.load разрешается сразу, так что на переходах это бесплатно.
  loader: () => waitForFonts(),
  component: Outlet,
  notFoundComponent: () => <Navigate to="/" />,
});

const journalRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: async ({ context, location }) => {
    const session = await context.queryClient.ensureQueryData(sessionQuery);
    if (!session) {
      throw redirect({ to: "/login", search: { return_to: location.href === "/" ? undefined : location.href } });
    }
  },
  // prefetch, а не ensure: ошибку загрузки записей показывает сама страница
  loader: ({ context }) => context.queryClient.prefetchQuery(bulletsQuery),
  component: JournalPage,
});

// Залогиненному на страницах входа/регистрации делать нечего.
const guestOnly = async ({ context }: { context: RouterContext }) => {
  const session = await context.queryClient.ensureQueryData(sessionQuery);
  if (session) throw redirect({ to: "/" });
};

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  validateSearch: validateAuthSearch,
  beforeLoad: guestOnly,
  loaderDeps: ({ search }) => ({ flow: search.flow }),
  loader: ({ context, deps }) => context.queryClient.ensureQueryData(authFlowQuery("login", deps.flow)),
  errorComponent: AuthFlowFailed,
  component: LoginPage,
});

const registrationRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/registration",
  validateSearch: validateAuthSearch,
  beforeLoad: guestOnly,
  loaderDeps: ({ search }) => ({ flow: search.flow }),
  loader: ({ context, deps }) => context.queryClient.ensureQueryData(authFlowQuery("registration", deps.flow)),
  errorComponent: AuthFlowFailed,
  component: RegistrationPage,
});

const errorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/error",
  validateSearch: (search: Record<string, unknown>): { id?: string } => ({ id: optionalString(search.id) }),
  component: AuthErrorPage,
});

const routeTree = rootRoute.addChildren([journalRoute, loginRoute, registrationRoute, errorRoute]);

export function createAppRouter(queryClient: QueryClient) {
  return createRouter({
    routeTree,
    context: { queryClient },
    // соседняя страница грузится по наведению на ссылку — к клику она уже готова
    defaultPreload: "intent",
    // свежесть данных решает TanStack Query, роутер всегда зовёт loader (повтор берётся из кэша)
    defaultPreloadStaleTime: 0,
    // переход внутри SPA затянулся дольше секунды — тоже сплэш
    defaultPendingComponent: Splash,
  });
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof createAppRouter>;
  }
}
