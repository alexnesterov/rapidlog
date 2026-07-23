import { useCallback, useEffect, useState } from "react";
import { listBullets } from "./api/bulletsApi";
import type { Bullet } from "./types/bullet";
import { BulletForm } from "./components/BulletForm";
import { BulletList } from "./components/BulletList";
import "./App.css";

function App() {
  const [bullets, setBullets] = useState<Bullet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(() => {
    setLoading(true);
    setError(null);
    listBullets()
      .then((res) => setBullets(res.bullets))
      .catch(() => setError("не удалось загрузить записи"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    reload();
  }, [reload]);

  const openCount = bullets.filter((b) => b.status === "OPEN").length;

  return (
    <div className="page">
      <header className="page__header">
        <h1>
          Rapid<span className="page__accent">Log</span>
        </h1>
        <p className="page__subtitle">
          {loading ? "синхронизация…" : `${openCount} открыто из ${bullets.length}`}
        </p>
      </header>

      <main className="page__main">
        <BulletForm onCreated={reload} />
        <BulletList bullets={bullets} loading={loading} error={error} />
      </main>

      <footer className="page__footer">данные не сохраняются — демо-режим на моках</footer>
    </div>
  );
}

export default App;
