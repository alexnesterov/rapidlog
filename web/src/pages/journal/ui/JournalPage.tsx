import { useCallback, useMemo, useState } from 'react';
import type { Bullet, MigrateTarget } from '../model/bullet';
import { useBulletsQuery, useCancelBulletMutation, useCompleteBulletMutation, useMigrateBulletMutation } from '../api/bulletsQueries';
import { useRouter } from '@tanstack/react-router';
import { useLogout, useSession } from '@/entities/session';
import { todayIsoDate } from '@/shared/lib';
import { Logo } from '@/shared/ui/logo';
import { Splash } from '@/shared/ui/splash';
import { DaySection } from './DaySection';
import { PocketSection } from './PocketSection';
import { usePocket, type PocketItem } from '../model/pocket';
import './JournalPage.css';

export function JournalPage() {
  const { data: days = [], isPending, isFetching, isError: listError } = useBulletsQuery();
  const completeMutation = useCompleteBulletMutation();
  const cancelMutation = useCancelBulletMutation();
  const migrateMutation = useMigrateBulletMutation();
  const { data: session } = useSession();
  const router = useRouter();
  const logoutMutation = useLogout(() => router.navigate({ to: '/login', search: {} }));
  const pocket = usePocket();
  const { takeToday, cancel: cancelInPocket } = pocket;

  const [error, setError] = useState<string | null>(null);

  const completeBullet = useCallback(
    (bullet: Bullet) => {
      completeMutation.mutate(bullet, {
        onSuccess: () => setError(null),
        onError: () => setError('не удалось обновить запись'),
      });
    },
    [completeMutation],
  );

  const moveBullet = useCallback(
    (bullet: Bullet, target: MigrateTarget) => {
      if (target !== 'today') return;
      migrateMutation.mutate(
        { bullet, target },
        {
          onSuccess: () => setError(null),
          onError: () => setError('не удалось перенести запись'),
        },
      );
    },
    [migrateMutation],
  );

  const cancelBulletEntry = useCallback(
    (bullet: Bullet) => {
      cancelMutation.mutate(bullet, {
        onSuccess: () => setError(null),
        onError: () => setError('не удалось отменить запись'),
      });
    },
    [cancelMutation],
  );

  const takeFromPocket = useCallback(
    (item: PocketItem) => {
      takeToday(item).then(
        () => setError(null),
        () => setError('не удалось перенести задачу из кармана'),
      );
    },
    [takeToday],
  );

  const cancelFromPocket = useCallback(
    (item: PocketItem) => {
      cancelInPocket(item).then(
        () => setError(null),
        () => setError('не удалось отменить задачу из кармана'),
      );
    },
    [cancelInPocket],
  );

  const today = todayIsoDate();

  // server-задачи, убранные в карман, прячем из их дней; опустевший прошлый день не показываем
  const visibleDays = useMemo(
    () =>
      days
        .map((day) => ({ ...day, bullets: day.bullets.filter((b) => !pocket.hiddenIds.has(b.id)) }))
        .filter((day) => day.day === today || day.bullets.length > 0),
    [days, pocket.hiddenIds, today],
  );

  const displayError = listError ? 'не удалось загрузить записи' : error;

  // шрифты ждёт корневой loader, записи — loader маршрута; сюда попадаем, только если
  // записей в кэше всё же нет (например, prefetch ещё не успел)
  if (isPending) {
    return <Splash />;
  }

  return (
    <div className="page">
      <header className="page__header">
        <h1>
          <Logo />
        </h1>
        {session && (
          <div className="page__account">
            <span className="page__email">{session.email}</span>
            <button
              type="button"
              className="page__logout"
              onClick={() => logoutMutation.mutate()}
              disabled={logoutMutation.isPending}
            >
              выйти
            </button>
          </div>
        )}
      </header>

      {isFetching && <div className="sync-indicator" aria-hidden="true" />}

      {displayError && <p className="log-state log-state--error">{displayError}</p>}

      {!displayError && (
        <div className="days">
          {visibleDays.map((day, index) => (
            <DaySection
              key={day.day}
              date={day.day}
              bullets={day.bullets}
              isToday={day.day === today}
              onComplete={completeBullet}
              onMigrate={moveBullet}
              onCancel={cancelBulletEntry}
              onPocket={pocket.pocketBullet}
              onAddToPocket={pocket.pocketBullet}
              animationDelay={index * 70}
            >
              <PocketSection items={pocket.items} onTakeToday={takeFromPocket} onCancel={cancelFromPocket} />
            </DaySection>
          ))}
        </div>
      )}

    </div>
  );
}

