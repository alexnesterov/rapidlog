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

export function useCompleteBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (bullet: Bullet) => markBulletDone(bullet.id),
    onSuccess: (updated) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => (prev ? replaceBullet(prev, updated) : prev));
    },
  });
}

export function useCancelBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (bullet: Bullet) => cancelBullet(bullet.id),
    onSuccess: (updated) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => (prev ? replaceBullet(prev, updated) : prev));
    },
  });
}

export function useMigrateBulletMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (vars: { bullet: Bullet; target: MigrateTarget }) => migrateBullet(vars.bullet.id),
    onSuccess: (created, vars) => {
      queryClient.setQueryData<BulletDayGroup[]>(bulletsQueryKey, (prev) => {
        if (!prev) return prev;
        const closed = replaceBullet(prev, { ...vars.bullet, signifier: "migrated" });
        return appendToToday(closed, created);
      });
    },
  });
}
