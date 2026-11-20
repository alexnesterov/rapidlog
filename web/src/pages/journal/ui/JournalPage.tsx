import { useCallback, useState } from 'react';
import type { Bullet, MigrateTarget } from '../model/bullet';
import { useBulletsQuery, useCancelBulletMutation, useCompleteBulletMutation, useMigrateBulletMutation } from '../api/bulletsQueries';
import { useRouter } from '@tanstack/react-router';
import { useLogout, useSession } from '@/entities/session';
import { ThemeMenu } from '@/features/theme-switch';
import { todayIsoDate } from '@/shared/lib';
import { Logo } from '@/shared/ui/logo';
import { Splash } from '@/shared/ui/splash';
import { AccountMenu } from './AccountMenu';
import { DaySection } from './DaySection';
import './JournalPage.css';

export function JournalPage() {
  const { data: days = [], isPending, isFetching, isError: listError } = useBulletsQuery();
  const completeMutation = useCompleteBulletMutation();
  const cancelMutation = useCancelBulletMutation();
  const migrateMutation = useMigrateBulletMutation();
  const { data: session } = useSession();
  const router = useRouter();
  const logoutMutation = useLogout(() => router.navigate({ to: '/login', search: {} }));

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

  const today = todayIsoDate();
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
          <Logo size="compact" />
          <span className="visually-hidden">: журнал</span>
        </h1>
        <div className="page__actions">
          <ThemeMenu />
          {session && (
            <AccountMenu
              email={session.email}
              onLogout={() => logoutMutation.mutate()}
              pending={logoutMutation.isPending}
            />
          )}
        </div>
      </header>

      {isFetching && <div className="sync-indicator" aria-hidden="true" />}

      {displayError && <p className="log-state log-state--error">{displayError}</p>}

      {!displayError && (
        <div className="days">
          {days.map((day, index) => (
            <DaySection
              key={day.day}
              date={day.day}
              bullets={day.bullets}
              isToday={day.day === today}
              onComplete={completeBullet}
              onMigrate={moveBullet}
              onCancel={cancelBulletEntry}
              animationDelay={index * 70}
            />
          ))}
        </div>
      )}

    </div>
  );
}

