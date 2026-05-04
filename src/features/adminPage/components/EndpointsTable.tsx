import React, { useState } from 'react';
import styles from '../AdminPage.module.css';

interface Endpoint {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  path: string;
  latency: string;
  enabled: boolean;
}

const initial: Endpoint[] = [
  { method: 'GET',    path: '/api/v1/health',          latency: '2ms',  enabled: true },
  { method: 'POST',   path: '/api/v1/auth/login',      latency: '45ms', enabled: true },
  { method: 'POST',   path: '/api/v1/auth/register',   latency: '62ms', enabled: true },
  { method: 'GET',    path: '/api/v1/users/:id',        latency: '18ms', enabled: true },
  { method: 'GET',    path: '/api/v1/tables',           latency: '12ms', enabled: true },
  { method: 'POST',   path: '/api/v1/bookings',         latency: '34ms', enabled: true },
  { method: 'GET',    path: '/api/v1/bookings/:id',     latency: '15ms', enabled: true },
  { method: 'PUT',    path: '/api/v1/bookings/:id',     latency: '28ms', enabled: true },
  { method: 'DELETE', path: '/api/v1/bookings/:id',     latency: '22ms', enabled: true },
  { method: 'GET',    path: '/api/v1/prices',           latency: '8ms',  enabled: true },
  { method: 'GET',    path: '/metrics',                 latency: '5ms',  enabled: true },
  { method: 'GET',    path: '/api/v1/admin/stats',      latency: '31ms', enabled: false },
];

const EndpointsTable: React.FC = () => {
  const [eps, setEps] = useState(initial);

  const toggle = (i: number) =>
    setEps(prev => prev.map((e, idx) => idx === i ? { ...e, enabled: !e.enabled } : e));

  const test = (ep: Endpoint) =>
    alert(`${ep.method} ${ep.path}\n\nStatus: 200 OK\nLatency: ${ep.latency}\nCache: ${Math.random() > 0.4 ? 'HIT' : 'MISS'}`);

  const active = eps.filter(e => e.enabled).length;

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Эндпоинты</div>
          <h2 className={styles.pageTitle}>Управление API</h2>
        </div>
        <span className={`${styles.badge} ${styles.badgeGreen}`} style={{ padding: '8px 16px', fontSize: '10px' }}>
          {active} / {eps.length} активны
        </span>
      </div>

      <div className={styles.endpointsList}>
        {eps.map((ep, i) => (
          <div key={i} className={styles.endpointRow}>
            <span className={`${styles.epMethod} ${styles[ep.method]}`}>{ep.method}</span>
            <span className={styles.epPath}>{ep.path}</span>
            <span className={styles.epLatency}>{ep.latency}</span>
            <button className={styles.epTestBtn} onClick={() => test(ep)}>Тест</button>
            <button
              className={`${styles.toggle} ${ep.enabled ? styles.on : ''}`}
              onClick={() => toggle(i)}
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default EndpointsTable;
