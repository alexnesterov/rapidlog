const WEEKDAY_FORMAT = new Intl.DateTimeFormat("ru-RU", { weekday: "long" });
const MONTH_FORMAT = new Intl.DateTimeFormat("ru-RU", { month: "long" });

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

export function dayNumber(iso: string): number {
  try {
    return parseIsoDate(iso).getDate();
  } catch {
    return 0;
  }
}

export function weekday(iso: string): string {
  try {
    return WEEKDAY_FORMAT.format(parseIsoDate(iso));
  } catch {
    return "";
  }
}

export function month(iso: string): string {
  try {
    return MONTH_FORMAT.format(parseIsoDate(iso));
  } catch {
    return "";
  }
}
