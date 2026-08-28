import { useCallback, useEffect, useRef, useState } from 'react';
import { listBullets, markBulletDone, migrateBullet } from './api/bulletsApi';
import type { Bullet, BulletDayGroup, MigrateTarget } from './types/bullet';
import { todayIsoDate } from './lib/date';
import { waitForFonts } from './lib/fonts';
import { DaySection } from './components/DaySection';
import './App.css';

function withToday(days: BulletDayGroup[]): BulletDayGroup[] {
  const today = todayIsoDate();
  if (days.some((day) => day.day === today)) return days;
  return [{ day: today, bullets: [] }, ...days];
}

function App() {
  const [days, setDays] = useState<BulletDayGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [initialized, setInitialized] = useState(false);
  const fontsReady = useRef(waitForFonts()).current;

  const reload = useCallback(() => {
    setLoading(true);
    setError(null);
    listBullets()
      .then((res) => setDays(withToday(res)))
      .catch(() => setError('не удалось загрузить записи'))
      .finally(() => {
        fontsReady.then(() => {
          setLoading(false);
          setInitialized(true);
        });
      });
  }, [fontsReady]);

  useEffect(() => {
    reload();
  }, [reload]);

  const completeBullet = useCallback(
    (bullet: Bullet) => {
      markBulletDone(bullet.id)
        .then(reload)
        .catch(() => setError('не удалось обновить запись'));
    },
    [reload],
  );

  const moveBullet = useCallback(
    (bullet: Bullet, target: MigrateTarget) => {
      if (target !== 'today') return;
      migrateBullet(bullet.id)
        .then(reload)
        .catch(() => setError('не удалось перенести запись'));
    },
    [reload],
  );

  const today = todayIsoDate();

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

      {loading && <div className="sync-indicator" aria-hidden="true" />}

      {error && <p className="log-state log-state--error">{error}</p>}

      {!error && (
        <div className="days">
          {days.map((day, index) => (
            <DaySection
              key={day.day}
              date={day.day}
              bullets={day.bullets}
              isToday={day.day === today}
              onCreated={reload}
              onComplete={completeBullet}
              onMigrate={moveBullet}
              animationDelay={index * 70}
            />
          ))}
        </div>
      )}

    </div>
  );
}

export default App;
