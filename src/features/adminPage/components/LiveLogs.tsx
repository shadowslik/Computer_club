import React, { useState, useEffect, useRef, useCallback } from 'react';
import styles from '../AdminPage.module.css';

type LogLevel = 'INFO' | 'WARN' | 'ERROR' | 'DEBUG';
interface LogEntry { time: string; level: LogLevel; message: string; }

const LOG_POOL: [LogLevel, string][] = [
  ['INFO',  'GET /api/v1/tables → 200 OK (12ms) [cache:HIT]'],
  ['INFO',  'POST /api/v1/auth/login → 200 OK (45ms) user:1042'],
  ['DEBUG', 'Redis GET tables:availability → HIT (0.3ms)'],
  ['INFO',  'GET /api/v1/prices → 200 OK (8ms) [cache:HIT]'],
  ['WARN',  'Rate limit approaching for IP 185.42.x.x (88/100 req/min)'],
  ['INFO',  'POST /api/v1/bookings → 201 Created (34ms) booking:B-2842'],
  ['DEBUG', 'Redis SET booking:B-2842 TTL=3600'],
  ['INFO',  'GET /metrics → 200 OK (5ms) [Prometheus scrape]'],
  ['ERROR', 'Java backend timeout: /api/v1/users/1099 (>5000ms) — circuit breaker open'],
  ['INFO',  'Circuit breaker reset for /api/v1/users'],
  ['DEBUG', 'Goroutines: 148 | HeapAlloc: 312MB | GC cycles: 7'],
  ['WARN',  'Slow query: GET /api/v1/bookings?date=2025-04-30 (850ms)'],
  ['INFO',  'DELETE /api/v1/bookings/B-2837 → 200 OK (22ms)'],
  ['DEBUG', 'Redis EXPIRE session:a1b2c3 1800'],
];

const LiveLogs: React.FC = () => {
  const [logs, setLogs] = useState<LogEntry[]>(() =>
    Array.from({ length: 20 }, (_, i) => {
      const p = LOG_POOL[Math.floor(Math.random() * LOG_POOL.length)];
      return {
        time: new Date(Date.now() - i * 4000).toLocaleTimeString('ru-RU'),
        level: p[0], message: p[1],
      };
    })
  );
  const [filter, setFilter] = useState<LogLevel | 'ALL'>('ALL');
  const [paused, setPaused] = useState(false);
  const termRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const t = setInterval(() => {
      if (paused) return;
      const p = LOG_POOL[Math.floor(Math.random() * LOG_POOL.length)];
      setLogs(prev => [{
        time: new Date().toLocaleTimeString('ru-RU'),
        level: p[0], message: p[1],
      }, ...prev].slice(0, 150));
    }, 1400);
    return () => clearInterval(t);
  }, [paused]);

  const filtered = filter === 'ALL' ? logs : logs.filter(l => l.level === filter);

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Логи</div>
          <h2 className={styles.pageTitle}>Live Логи</h2>
        </div>
      </div>

      <div className={styles.logFilters}>
        {(['ALL', 'INFO', 'WARN', 'ERROR', 'DEBUG'] as const).map(lvl => (
          <button
            key={lvl}
            className={`${styles.logFilter} ${filter === lvl ? styles.active : ''}`}
            onClick={() => setFilter(lvl)}
          >
            {lvl}
          </button>
        ))}
        <button className={styles.pauseBtn} onClick={() => setPaused(p => !p)}>
          {paused ? '▶ Продолжить' : '⏸ Пауза'}
        </button>
      </div>

      <div className={styles.logContainer} ref={termRef}>
        {filtered.map((log, i) => (
          <div key={i} className={styles.logLine}>
            <span className={styles.logTime}>{log.time}</span>
            <span className={`${styles.logLevel} ${styles[log.level]}`}>{log.level}</span>
            <span className={styles.logMessage}>{log.message}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LiveLogs;
