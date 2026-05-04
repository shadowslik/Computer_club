import React, { useState } from 'react';
import styles from '../AdminPage.module.css';

interface RedisKey { key: string; type: string; ttl: string; size: string; }

const initial: RedisKey[] = [
  { key: 'user:session:a1b2c3',   type: 'STRING', ttl: '1800s',  size: '248B' },
  { key: 'tables:availability',   type: 'STRING', ttl: '30s',    size: '1.2KB' },
  { key: 'prices:all',            type: 'STRING', ttl: '3600s',  size: '856B' },
  { key: 'booking:list:2025',     type: 'LIST',   ttl: '120s',   size: '4.8KB' },
  { key: 'user:profile:1042',     type: 'HASH',   ttl: '900s',   size: '512B' },
  { key: 'ratelimit:185.42.x.x',  type: 'STRING', ttl: '58s',    size: '32B' },
  { key: 'games:catalog',         type: 'STRING', ttl: '86400s', size: '12KB' },
  { key: 'metrics:snapshot',      type: 'HASH',   ttl: '10s',    size: '2.1KB' },
];

const RedisPanel: React.FC = () => {
  const [keys, setKeys] = useState(initial);

  const del = (i: number) => setKeys(p => p.filter((_, idx) => idx !== i));
  const flush = () => {
    if (window.confirm('Очистить весь кэш Redis?')) { setKeys([]); }
  };

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Redis Cache</div>
          <h2 className={styles.pageTitle}>Redis Cache</h2>
        </div>
        <button className={styles.dangerBtn} onClick={flush}>⚠ Flush All</button>
      </div>

      <div className={styles.redisStats}>
        <div className={styles.statBox}>
          <div className={`${styles.statValue}`}>{keys.length}</div>
          <div className={styles.statLabel}>Ключей в кэше</div>
        </div>
        <div className={styles.statBox}>
          <div className={`${styles.statValue} ${styles.clrGreen}`}>87%</div>
          <div className={styles.statLabel}>Hit Rate</div>
        </div>
        <div className={styles.statBox}>
          <div className={styles.statValue}>42MB</div>
          <div className={styles.statLabel}>Использовано</div>
        </div>
        <div className={styles.statBox}>
          <div className={`${styles.statValue} ${styles.clrBlue}`}>256MB</div>
          <div className={styles.statLabel}>Лимит памяти</div>
        </div>
      </div>

      <div className={styles.tableWrap}>
        <div className={styles.tableHead}>
          <span className={styles.tableHeadTitle}>Ключи кэша</span>
          <span className={styles.countBadge}>Показано {keys.length} из {initial.length}</span>
        </div>
        <table className={styles.dataTable}>
          <thead>
            <tr>
              <th>Ключ</th><th>Тип</th><th>TTL</th><th>Размер</th><th></th>
            </tr>
          </thead>
          <tbody>
            {keys.map((k, i) => (
              <tr key={i}>
                <td style={{ color: 'rgba(255,255,255,0.72)' }}>{k.key}</td>
                <td><span className={`${styles.badge} ${styles.badgeGray}`}>{k.type}</span></td>
                <td>{k.ttl}</td>
                <td>{k.size}</td>
                <td>
                  <button className={styles.dangerBtnSmall} onClick={() => del(i)}>DEL</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default RedisPanel;
