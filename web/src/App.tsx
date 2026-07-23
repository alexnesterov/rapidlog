import { useCallback, useEffect, useState } from "react";
import { listBullets, markBulletDone } from "./api/bulletsApi";
import type { Bullet } from "./types/bullet";
import { todayIsoDate } from "./lib/date";
import { DaySection } from "./components/DaySection";
import "./App.css";

interface DayGroup {
  date: string;
  bullets: Bullet[];
}

function groupByDate(bullets: Bullet[]): DayGroup[] {
  const byDate = new Map<string, Bullet[]>();
  for (const bullet of bullets) {
    const group = byDate.get(bullet.date);
    if (group) {
      group.push(bullet);
    } else {
      byDate.set(bullet.date, [bullet]);
    }
  }

  const today = todayIsoDate();
  if (!byDate.has(today)) {
    byDate.set(today, []);
  }

  return [...byDate.entries()]
    .sort((a, b) => b[0].localeCompare(a[0]))
    .map(([date, dayBullets]) => ({ date, bullets: [...dayBullets].reverse() }));
}

function App() {
  const [bullets, setBullets] = useState<Bullet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [initialized, setInitialized] = useState(false);

  const reload = useCallback(() => {
    setLoading(true);
    setError(null);
    listBullets()
      .then((res) => setBullets(res.bullets))
      .catch(() => setError("не удалось загрузить записи"))
      .finally(() => {
        setLoading(false);
        setInitialized(true);
      });
  }, []);

  useEffect(() => {
    reload();
  }, [reload]);

  const completeBullet = useCallback(
    (bullet: Bullet) => {
      markBulletDone(bullet.id)
        .then(reload)
        .catch(() => setError("не удалось обновить запись"));
    },
    [reload],
  );

  const today = todayIsoDate();
  const days = groupByDate(bullets);

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
              key={day.date}
              date={day.date}
              bullets={day.bullets}
              isToday={day.date === today}
              onCreated={reload}
              onComplete={completeBullet}
              animationDelay={index * 70}
            />
          ))}
        </div>
      )}

      <footer className="page__footer">данные не сохраняются — демо-режим на моках</footer>
    </div>
  );
}

export default App;
