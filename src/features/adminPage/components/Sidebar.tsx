import React from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Sidebar.module.css';
import { AdminSection } from '../AdminPage';

interface SidebarProps {
  activeSection: AdminSection;
  onSelect: (section: AdminSection) => void;
}

const nav: { section: string; items: { id: AdminSection; label: string; icon: string }[] }[] = [
  {
    section: 'Мониторинг',
    items: [
      { id: 'metrics',   label: 'Метрики',   icon: '◈' },
      { id: 'logs',      label: 'Live логи',  icon: '≡' },
    ],
  },
  {
    section: 'Gateway',
    items: [
      { id: 'endpoints', label: 'Эндпоинты', icon: '⟐' },
      { id: 'redis',     label: 'Redis Cache', icon: '⬡' },
    ],
  },
  {
    section: 'Данные',
    items: [
      { id: 'users',     label: 'Пользователи', icon: '◎' },
      { id: 'bookings',  label: 'Бронирования', icon: '◷' },
    ],
  },
];

const Sidebar: React.FC<SidebarProps> = ({ activeSection, onSelect }) => {
  const navigate = useNavigate();

  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>
        <div className={styles.wordmark}>CYBER<span>ADMIN</span></div>
        <span className={styles.goBadge}>● Go Gateway · 8080</span>
      </div>

      <nav>
        {nav.map(group => (
          <div key={group.section} className={styles.navSection}>
            <span className={styles.navLabel}>{group.section}</span>
            {group.items.map(item => (
              <button
                key={item.id}
                className={`${styles.navItem} ${activeSection === item.id ? styles.active : ''}`}
                onClick={() => onSelect(item.id)}
              >
                <span className={styles.navIcon}>{item.icon}</span>
                <span>{item.label}</span>
              </button>
            ))}
          </div>
        ))}
      </nav>

      <div className={styles.bottom}>
        <div className={styles.statusRow}>
          <span className={styles.statusDot} />
          <span>gateway online</span>
        </div>
        <button className={styles.toSiteBtn} onClick={() => navigate('/')}>
          ← Вернуться на сайт
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
