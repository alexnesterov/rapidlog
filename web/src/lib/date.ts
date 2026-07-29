const WEEKDAY_FORMAT = new Intl.DateTimeFormat("ru-RU", { weekday: "long" });
const MONTH_FORMAT = new Intl.DateTimeFormat("ru-RU", { month: "long" });

export function todayIsoDate(): string {
  return new Date().toISOString().slice(0, 10);
}

function parseIsoDate(iso: string): Date {
  return new Date(iso + "T00:00:00");
}

export function dayNumber(iso: string): number {
  return parseIsoDate(iso).getDate();
}

export function weekday(iso: string): string {
  return WEEKDAY_FORMAT.format(parseIsoDate(iso));
}

export function month(iso: string): string {
  return MONTH_FORMAT.format(parseIsoDate(iso));
}
