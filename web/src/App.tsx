import { useCallback, useEffect, useRef, useState } from 'react';
import type { Bullet, MigrateTarget } from './types/bullet';
import { useBulletsQuery, useCancelBulletMutation, useCompleteBulletMutation, useMigrateBulletMutation } from './api/bulletsQueries';
import { todayIsoDate } from './lib/date';
import { waitForFonts } from './lib/fonts';
import { DaySection } from './components/DaySection';
import './App.css';

function App() {
  const { data: days = [], isPending, isFetching, isError: listError } = useBulletsQuery();
  const completeMutation = useCompleteBulletMutation();
  const cancelMutation = useCancelBulletMutation();
  const migrateMutation = useMigrateBulletMutation();

  const [error, setError] = useState<string | null>(null);
  const [fontsLoaded, setFontsLoaded] = useState(false);
  const fontsReady = useRef(waitForFonts()).current;

  useEffect(() => {
    fontsReady.then(() => setFontsLoaded(true));
  }, [fontsReady]);

  const initialized = fontsLoaded && !isPending;

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

  if (!initialized) {
    return (
      <div className="splash">
        <div className="splash__mark" aria-hidden="true">
          <span className="splash__dot" />
        </div>
        <div className="splash__spinner" aria-hidden="true" />
      </div>
    );
  }

  return (
    <div className="page">
      <div className="page__stamp" aria-hidden="true">
        preview
      </div>

      <header className="page__header">
        <h1>
          Rapid<span className="page__accent">Log</span>
        </h1>
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

export default App;
