import { useCallback, useEffect, useMemo, useState } from "react";
import type { Bullet } from "./bullet";
import { useCancelBulletMutation, useCreateBulletMutation, useMigrateBulletMutation } from "../api/bulletsQueries";

// UI-прототип pocket. Каждая задача кармана — обычный bullet на сервере (open task),
// а «лежит в кармане» — пока только пометка на клиенте (localStorage): карман прячет
// эти bullet'ы из их дней. С бэкендом пометка станет колонкой in_pocket, и этот модуль
// перейдёт на API — остальной UI не меняется.

const STORAGE_KEY = "rapidlog.pocket.v2";
// первая версия прототипа держала задачи, записанные в карман, только в браузере
const LEGACY_STORAGE_KEY = "rapidlog.pocket";

export type PocketItem = Bullet;

interface LegacyPocketItem {
  content: string;
  source: "local" | "server";
  bullet?: Bullet;
}

function load(): PocketItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    const parsed: unknown = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed) ? (parsed as PocketItem[]) : [];
  } catch {
    return [];
  }
}

function save(items: PocketItem[]) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
  } catch {
    // приватное окно / заблокированное хранилище — карман живёт до перезагрузки
  }
}

function takeLegacy(): LegacyPocketItem[] {
  try {
    const raw = localStorage.getItem(LEGACY_STORAGE_KEY);
    if (!raw) return [];
    localStorage.removeItem(LEGACY_STORAGE_KEY);
    const parsed: unknown = JSON.parse(raw);
    return Array.isArray(parsed) ? (parsed as LegacyPocketItem[]) : [];
  } catch {
    return [];
  }
}

function localIsoDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

function createdToday(bullet: Bullet): boolean {
  return localIsoDate(bullet.createdAt) === localIsoDate(new Date().toISOString());
}

export function usePocket() {
  const [items, setItems] = useState<PocketItem[]>(load);
  const { mutateAsync: create } = useCreateBulletMutation();
  const { mutateAsync: migrate } = useMigrateBulletMutation();
  const { mutateAsync: cancelBullet } = useCancelBulletMutation();

  useEffect(() => save(items), [items]);

  const hiddenIds = useMemo(() => new Set(items.map((i) => i.id)), [items]);

  const pocketBullet = useCallback((bullet: Bullet) => {
    setItems((prev) => (prev.some((i) => i.id === bullet.id) ? prev : [...prev, bullet]));
  }, []);

  // Перенос кармана из первой версии прототипа: задачи, которые жили только в браузере,
  // создаём на сервере; убранные из дней берём как есть.
  useEffect(() => {
    const legacy = takeLegacy();
    if (legacy.length === 0) return;
    void (async () => {
      for (const old of legacy) {
        try {
          pocketBullet(old.bullet ?? (await create({ content: old.content, kind: "KIND_TASK" })));
        } catch {
          // не удалось создать — задача теряется, для прототипа приемлемо
        }
      }
    })();
  }, [create, pocketBullet]);

  // Убираем из кармана сразу (оптимистично) и выполняем действие на сервере;
  // при ошибке кладём задачу обратно на то же место.
  const removeWith = useCallback(async (item: PocketItem, action: () => Promise<unknown>) => {
    let index = -1;
    setItems((prev) => {
      index = prev.findIndex((i) => i.id === item.id);
      return prev.filter((i) => i.id !== item.id);
    });

    try {
      await action();
    } catch (err) {
      setItems((prev) => {
        if (prev.some((i) => i.id === item.id)) return prev;
        const next = [...prev];
        next.splice(index < 0 ? next.length : index, 0, item);
        return next;
      });
      throw err;
    }
  }, []);

  // Достать: задача, записанная сегодня, и так лежит в сегодняшнем дне — достаточно
  // перестать её прятать; более старую переносим в сегодня.
  const takeToday = useCallback(
    (item: PocketItem) =>
      removeWith(item, async () => {
        if (!createdToday(item)) {
          await migrate({ bullet: item, target: "today" });
        }
      }),
    [removeWith, migrate],
  );

  // Отменить: обычный cancel — задача остаётся в журнале зачёркнутой в своём дне.
  const cancel = useCallback(
    (item: PocketItem) => removeWith(item, () => cancelBullet(item)),
    [removeWith, cancelBullet],
  );

  return { items, hiddenIds, pocketBullet, takeToday, cancel };
}
