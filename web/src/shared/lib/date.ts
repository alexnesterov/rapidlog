// bullet-journal style day head: "7.11 сб"
const WEEKDAY_FORMAT = new Intl.DateTimeFormat("ru-RU", { weekday: "short" });

export function todayIsoDate(): string {
  return new Date().toISOString().slice(0, 10);
}

function parseIsoDate(iso: string): Date {
  const date = new Date(iso + "T00:00:00");
  if (isNaN(date.getTime())) {
    throw new Error(`Invalid date: ${iso}`);
  }
  return date;
}

export function dayMonth(iso: string): string {
  try {
    const date = parseIsoDate(iso);
    return `${date.getDate()}.${String(date.getMonth() + 1).padStart(2, "0")}`;
  } catch {
    return "";
  }
}

export function weekday(iso: string): string {
  try {
    return WEEKDAY_FORMAT.format(parseIsoDate(iso));
  } catch {
    return "";
  }
}
