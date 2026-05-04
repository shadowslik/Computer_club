import React, { useState } from 'react';
import styles from '../AdminPage.module.css';

const mockUsers = [
  { id: 1001, name: 'Gamer_777',   email: 'gamer@mail.com',   balance: '2 400₽', bookings: 12, status: 'active',  reg: '12.01.2025' },
  { id: 1002, name: 'ProPlayer99', email: 'pro@gmail.com',    balance: '850₽',   bookings: 47, status: 'active',  reg: '03.11.2024' },
  { id: 1003, name: 'NightStalker',email: 'night@yandex.ru',  balance: '0₽',     bookings: 3,  status: 'banned',  reg: '28.02.2025' },
  { id: 1004, name: 'VoidWalker',  email: 'void@mail.ru',     balance: '5 200₽', bookings: 89, status: 'vip',     reg: '15.08.2024' },
  { id: 1005, name: 'IronFist_23', email: 'iron@gmail.com',   balance: '320₽',   bookings: 7,  status: 'active',  reg: '01.03.2025' },
  { id: 1006, name: 'ShadowSniper',email: 'shadow@inbox.ru',  balance: '1 100₽', bookings: 21, status: 'active',  reg: '20.12.2024' },
];

const statusMap = {
  active: { cls: styles.badgeGreen,  label: 'Активен' },
  banned: { cls: styles.badgeRed,    label: 'Заблокирован' },
  vip:    { cls: styles.badgePurple, label: 'VIP' },
} as const;

const UsersTable: React.FC = () => {
  const [search, setSearch] = useState('');

  const filtered = mockUsers.filter(u =>
    !search || u.name.toLowerCase().includes(search.toLowerCase()) || u.email.includes(search)
  );

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Пользователи</div>
          <h2 className={styles.pageTitle}>Пользователи</h2>
        </div>
        <span className={`${styles.badge} ${styles.badgeBlue}`} style={{ padding: '8px 16px', fontSize: '10px' }}>
          Всего: 1 842
        </span>
      </div>

      <div className={styles.searchRow}>
        <input
          className={styles.searchInput}
          placeholder="Поиск по имени, email, телефону..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <button className={styles.searchBtn}>Найти</button>
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.dataTable}>
          <thead>
            <tr>
              <th>ID</th><th>Никнейм</th><th>Email</th>
              <th>Баланс</th><th>Бронирований</th><th>Статус</th><th>Регистрация</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(u => {
              const s = statusMap[u.status as keyof typeof statusMap];
              return (
                <tr key={u.id}>
                  <td style={{ color: 'rgba(255,255,255,0.35)' }}>#{u.id}</td>
                  <td style={{ fontWeight: 600, color: '#fff' }}>{u.name}</td>
                  <td style={{ color: 'rgba(255,255,255,0.5)' }}>{u.email}</td>
                  <td>{u.balance}</td>
                  <td style={{ textAlign: 'center' }}>{u.bookings}</td>
                  <td><span className={`${styles.badge} ${s.cls}`}>{s.label}</span></td>
                  <td style={{ color: 'rgba(255,255,255,0.4)' }}>{u.reg}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UsersTable;
