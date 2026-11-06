import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { cancelBullet, createBullet, listBullets, markBulletDone, migrateBullet } from "./bulletsApi";
import type { Bullet, BulletDayGroup, CreateBulletRequest, MigrateTarget } from "../model/bullet";
import { todayIsoDate } from "@/shared/lib";

export const bulletsQueryKey = ["bullets"] as const;

function withToday(days: BulletDayGroup[]): BulletDayGroup[] {
  const today = todayIsoDate();
  if (days.some((day) => day.day === today)) return days;
  return [{ day: today, bullets: [] }, ...days];
}

function replaceBullet(days: BulletDayGroup[], updated: Bullet): BulletDayGroup[] {
  return days.map((day) => ({
    ...day,
    bullets: day.bullets.map((b) => (b.id === updated.id ? updated : b)),
  }));
}

function appendToToday(days: BulletDayGroup[], bullet: Bullet): BulletDayGroup[] {
  const today = todayIsoDate();
  return withToday(days).map((day) => (day.day === today ? { ...day, bullets: [...day.bullets, bullet] } : day));
}

// Перенос не создаёт новую запись: тот же bullet уезжает в сегодняшний день.
// Опустевший прошлый день убираем — сервер пустых групп не отдаёт.
function moveToToday(days: BulletDayGroup[], bullet: Bullet): BulletDayGroup[] {
  const today = todayIsoDate();
  const without = days
    .map((day) => ({ ...day, bullets: day.bullets.filter((b) => b.id !== bullet.id) }))
    .filter((day) => day.day === today || day.bullets.length > 0);
  return appendToToday(without, bullet);
}

export function useBulletsQuery() {
  return useQuery({
    queryKey: bulletsQueryKey,
    queryFn: async () => withToday(await listBullets()),
  });
}

export function useCreateBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (req: CreateBulletRequest) => createBullet(req),
    onSuccess: (created) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => appendToToday(prev ?? [], created));
    },
  });
}

type BulletsSnapshot = { previous: BulletDayGroup[] | undefined };

export function useCompleteBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation<Bullet, unknown, Bullet, BulletsSnapshot>({
    mutationFn: (bullet: Bullet) => markBulletDone(bullet.id),
    onMutate: async (bullet) => {
      await queryClient.cancelQueries({ queryKey: bulletsQueryKey });
      const previous = queryClient.getQueryData<BulletDayGroup[]>(bulletsQueryKey);
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) =>
        prev ? replaceBullet(prev, { ...bullet, signifier: "SIGNIFIER_COMPLETED" }) : prev,
      );
      return { previous };
    },
    onError: (_err, _bullet, context) => {
      if (context) queryClient.setQueryData(bulletsQueryKey, context.previous);
    },
    onSuccess: (updated) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => (prev ? replaceBullet(prev, updated) : prev));
    },
  });
}

export function useCancelBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation<Bullet, unknown, Bullet, BulletsSnapshot>({
    mutationFn: (bullet: Bullet) => cancelBullet(bullet.id),
    onMutate: async (bullet) => {
      await queryClient.cancelQueries({ queryKey: bulletsQueryKey });
      const previous = queryClient.getQueryData<BulletDayGroup[]>(bulletsQueryKey);
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) =>
        prev ? replaceBullet(prev, { ...bullet, signifier: "SIGNIFIER_CANCELLED" }) : prev,
      );
      return { previous };
    },
    onError: (_err, _bullet, context) => {
      if (context) queryClient.setQueryData(bulletsQueryKey, context.previous);
    },
    onSuccess: (updated) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => (prev ? replaceBullet(prev, updated) : prev));
    },
  });
}

export function useMigrateBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation<Bullet, unknown, { bullet: Bullet; target: MigrateTarget }, BulletsSnapshot>({
    mutationFn: (vars: { bullet: Bullet; target: MigrateTarget }) => migrateBullet(vars.bullet.id),
    onMutate: async (vars) => {
      await queryClient.cancelQueries({ queryKey: bulletsQueryKey });
      const previous = queryClient.getQueryData<BulletDayGroup[]>(bulletsQueryKey);
      const now = new Date().toISOString();
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) =>
        prev ? moveToToday(prev, { ...vars.bullet, createdAt: now, updatedAt: now }) : prev,
      );
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context) queryClient.setQueryData(bulletsQueryKey, context.previous);
    },
    onSuccess: (migrated) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => (prev ? moveToToday(prev, migrated) : prev));
    },
  });
}
