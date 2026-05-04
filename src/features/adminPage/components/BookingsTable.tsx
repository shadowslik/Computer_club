import React, { useState } from 'react';
import styles from '../AdminPage.module.css';

const mockBookings = [
  { id: 'B-2841', user: 'Gamer_777',   table: 'Pro Zone #1',   type: 'pro',      dt: '30.04.2025 14:00', hours: 2, sum: '500₽',  status: 'active' },
  { id: 'B-2840', user: 'VoidWalker',  table: 'VIP Lounge #1', type: 'vip',      dt: '30.04.2025 12:00', hours: 3, sum: '1200₽', status: 'active' },
  { id: 'B-2839', user: 'ProPlayer99', table: 'Стандарт #1',   type: 'standard', dt: '30.04.2025 10:00', hours: 1, sum: '150₽',  status: 'done' },
  { id: 'B-2838', user: 'IronFist_23', table: 'Pro Zone #2',   type: 'pro',      dt: '29.04.2025 20:00', hours: 4, sum: '1000₽', status: 'cancelled' },
  { id: 'B-2837', user: 'ShadowSniper',table: 'Стандарт #2',   type: 'standard', dt: '29.04.2025 18:00', hours: 2, sum: '300₽',  status: 'cancelled' },
  { id: 'B-2836', user: 'NightStalker',table: 'VIP Lounge #2', type: 'vip',      dt: '01.05.2025 16:00', hours: 2, sum: '800₽',  status: 'pending' },
];

const statusMap = {
  active:    { cls: styles.badgeGreen,  label: 'Активно' },
  done:      { cls: styles.badgeGray,   label: 'Завершено' },
  cancelled: { cls: styles.badgeRed,    label: 'Отменено' },
  pending:   { cls: styles.badgeYellow, label: 'Ожидает' },
} as const;

const typeMap = {
  vip:      styles.badgePurple,
  pro:      styles.badgeBlue,
  standard: styles.badgeGray,
} as const;

const filters = ['Все', 'Активные', 'Ожидают', 'Завершены', 'Отменены'];

const BookingsTable: React.FC = () => {
  const [activeFilter, setActiveFilter] = useState('Все');

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Бронирования</div>
          <h2 className={styles.pageTitle}>Бронирования</h2>
        </div>
        <span className={`${styles.badge} ${styles.badgeYellow}`} style={{ padding: '8px 16px', fontSize: '10px' }}>
          Сегодня: 47
        </span>
      </div>

      <div className={styles.filterRow}>
        {filters.map(f => (
          <button
            key={f}
            className={`${styles.filterBtn} ${activeFilter === f ? styles.active : ''}`}
            onClick={() => setActiveFilter(f)}
          >
            {f}
          </button>
        ))}
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.dataTable}>
          <thead>
            <tr>
              <th>ID</th><th>Пользователь</th><th>Стол</th>
              <th>Тариф</th><th>Дата / Время</th><th>Длит.</th><th>Сумма</th><th>Статус</th>
            </tr>
          </thead>
          <tbody>
            {mockBookings.map(b => {
              const s = statusMap[b.status as keyof typeof statusMap];
              const tc = typeMap[b.type as keyof typeof typeMap];
              return (
                <tr key={b.id}>
                  <td style={{ color: 'rgba(255,255,255,0.35)' }}>{b.id}</td>
                  <td style={{ color: '#fff', fontWeight: 500 }}>{b.user}</td>
                  <td style={{ color: 'rgba(255,255,255,0.65)' }}>{b.table}</td>
                  <td><span className={`${styles.badge} ${tc}`}>{b.type.toUpperCase()}</span></td>
                  <td style={{ color: 'rgba(255,255,255,0.5)' }}>{b.dt}</td>
                  <td>{b.hours}ч</td>
                  <td style={{ fontWeight: 600 }}>{b.sum}</td>
                  <td><span className={`${styles.badge} ${s.cls}`}>{s.label}</span></td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default BookingsTable;
