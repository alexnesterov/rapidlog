import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { cancelBullet, createBullet, listBullets, markBulletDone, migrateBullet } from "./bulletsApi";
import type { Bullet, BulletDayGroup, CreateBulletRequest, MigrateTarget } from "../types/bullet";
import { todayIsoDate } from "../lib/date";

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
        prev ? replaceBullet(prev, { ...bullet, signifier: "completed" }) : prev,
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
        prev ? replaceBullet(prev, { ...bullet, signifier: "cancelled" }) : prev,
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
      // Оптимистично закрываем только исходную строку; новую "сегодняшнюю"
      // запись добавляем в onSuccess, когда известны её реальные id/content.
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) =>
        prev ? replaceBullet(prev, { ...vars.bullet, signifier: "migrated" }) : prev,
      );
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context) queryClient.setQueryData(bulletsQueryKey, context.previous);
    },
    onSuccess: (created, vars) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => {
        if (!prev) return prev;
        const closed = replaceBullet(prev, { ...vars.bullet, signifier: "migrated" });
        return appendToToday(closed, created);
      });
    },
  });
}
