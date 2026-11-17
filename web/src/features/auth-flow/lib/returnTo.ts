// return_to приходит из URL — пускаем только локальные пути, иначе это open redirect.
export function safeReturnPath(returnTo: string | undefined): string {
  if (!returnTo || !returnTo.startsWith("/") || returnTo.startsWith("//")) return "/";
  return returnTo;
}
