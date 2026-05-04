import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './LeaderboardPage.module.css';

// ─── Types ────────────────────────────────────────────────────────────────────
interface ResponseUser {
  id: number;
  name: string;
  email: string;
  phone: string;
  balance: number;
  hours: number;
  computerSessions: { id: number; tableType: string; hours: number; date: string }[];
}

type SortMode = 'hours' | 'sessions' | 'balance';

const API_BASE = 'http://localhost:8080/api/clients';
const PAGE_SIZE = 20;

// ─── Helpers ─────────────────────────────────────────────────────────────────
const getEndpoint = (mode: SortMode) => {
  switch (mode) {
    case 'hours':     return `${API_BASE}/sort_by_count_hours`;
    case 'sessions':  return `${API_BASE}/sort_by_count_sessions`;
    case 'balance':   return `${API_BASE}/sort_by_balance`;
  }
};

const getTypeIcon = (sessions: ResponseUser['computerSessions']) => {
  if (!sessions?.length) return '◈';
  const last = sessions[sessions.length - 1]?.tableType?.toLowerCase() ?? '';
  if (last.includes('vip'))  return '🔱';
  if (last.includes('pro'))  return '❖';
  return '◈';
};

const getUserRank = (user: ResponseUser): { label: string; cls: string } => {
  if (user.balance >= 5000) return { label: 'VIP',     cls: styles.rankVip };
  if (user.balance >= 2000) return { label: 'PRO',     cls: styles.rankPro };
  return                          { label: 'НОВИЧОК',  cls: styles.rankDefault };
};

const sortTabs: { id: SortMode; label: string; sublabel: string }[] = [
  { id: 'hours',    label: 'По часам',       sublabel: 'Всего часов в клубе' },
  { id: 'sessions', label: 'По сессиям',     sublabel: 'Количество бронирований' },
  { id: 'balance',  label: 'По балансу',     sublabel: 'Остаток на счёте' },
];

// ─── Component ───────────────────────────────────────────────────────────────
const LeaderboardPage: React.FC = () => {
  const navigate = useNavigate();

  const [sortMode, setSortMode]   = useState<SortMode>('hours');
  const [users, setUsers]         = useState<ResponseUser[]>([]);
  const [loading, setLoading]     = useState(true);
  const [error, setError]         = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    fetch(`${getEndpoint(sortMode)}?page=0&size=${PAGE_SIZE}`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json() as Promise<ResponseUser[]>;
      })
      .then(data => {
        if (!cancelled) setUsers(data);
      })
      .catch(err => {
        if (!cancelled) setError(err.message ?? 'Ошибка загрузки');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => { cancelled = true; };
  }, [sortMode]);

  const getStatValue = (user: ResponseUser) => {
    switch (sortMode) {
      case 'hours':    return `${user.hours ?? 0}ч`;
      case 'sessions': return `${user.computerSessions?.length ?? 0} сесс.`;
      case 'balance':  return `${user.balance?.toLocaleString('ru-RU') ?? 0} ₽`;
    }
  };

  const top3  = users.slice(0, 3);
  const rest  = users.slice(3);

  // podium order: 2nd, 1st, 3rd
  const podiumOrder = [top3[1], top3[0], top3[2]].filter(Boolean) as ResponseUser[];
  const podiumClass = (_u: ResponseUser, idx: number) => {
    if (idx === 1) return styles.podiumFirst;
    if (idx === 0) return styles.podiumSecond;
    return styles.podiumThird;
  };
  const podiumPlace = (_u: ResponseUser, idx: number) => (idx === 1 ? 1 : idx === 0 ? 2 : 3);

  const monthName = new Date().toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' });

  return (
    <div className={styles.page}>
      {/* BG decoration */}
      <div className={styles.bgGrid} />
      <div className={styles.bgText}>TOP</div>

      <div className={styles.inner}>
        {/* Back */}
        <button className={styles.backBtn} onClick={() => navigate('/')}>
          ← На главную
        </button>

        {/* Header */}
        <p className={styles.eyebrow}>{monthName}</p>
        <h1 className={styles.title}>
          Таблица рейтинга <span className={styles.titleMuted}>/ Top Players</span>
        </h1>

        {/* Sort tabs */}
        <div className={styles.tabs}>
          {sortTabs.map(tab => (
            <button
              key={tab.id}
              className={`${styles.tab} ${sortMode === tab.id ? styles.tabActive : ''}`}
              onClick={() => setSortMode(tab.id)}
            >
              <span className={styles.tabLabel}>{tab.label}</span>
              <span className={styles.tabSub}>{tab.sublabel}</span>
            </button>
          ))}
        </div>

        {/* Loading */}
        {loading && (
          <div className={styles.stateBox}>
            <div className={styles.spinner} />
            <span>Загрузка рейтинга...</span>
          </div>
        )}

        {/* Error */}
        {!loading && error && (
          <div className={styles.stateBox}>
            <div className={styles.errorIcon}>⚠</div>
            <span className={styles.errorText}>Ошибка: {error}</span>
            <button className={styles.retryBtn} onClick={() => setSortMode(s => s)}>
              Повторить
            </button>
          </div>
        )}

        {/* Empty */}
        {!loading && !error && users.length === 0 && (
          <div className={styles.stateBox}>
            <span className={styles.emptyText}>Данные ещё не загружены</span>
          </div>
        )}

        {/* Podium (top 3) */}
        {!loading && !error && top3.length >= 1 && (
          <div className={styles.podium}>
            {podiumOrder.map((user, idx) => {
              const place = podiumPlace(user, idx);
              const rank  = getUserRank(user);
              return (
                <div key={user.id} className={`${styles.podiumCard} ${podiumClass(user, idx)}`}>
                  <div className={styles.podiumPlaceBg}>{place}</div>
                  <div className={styles.podiumAvatar}>
                    {user.name.slice(0, 2).toUpperCase()}
                  </div>
                  <div className={styles.podiumName}>{user.name}</div>
                  <div className={styles.podiumMeta}>
                    <span className={`${styles.rankBadge} ${rank.cls}`}>{rank.label}</span>
                  </div>
                  <div className={styles.podiumStat}>{getStatValue(user)}</div>
                  <div className={styles.podiumPlace}>#{place}</div>
                </div>
              );
            })}
          </div>
        )}

        {/* Table (4+) */}
        {!loading && !error && rest.length > 0 && (
          <div className={styles.tableWrap}>
            <div className={styles.tableHead}>
              <span>#</span>
              <span>Игрок</span>
              <span className={styles.colRight}>
                {sortTabs.find(t => t.id === sortMode)?.label}
              </span>
            </div>
            {rest.map((user, idx) => {
              const place = idx + 4;
              const icon  = getTypeIcon(user.computerSessions);
              return (
                <div key={user.id} className={styles.tableRow}>
                  <span className={styles.rowNum}>{place}</span>
                  <div className={styles.rowPlayer}>
                    <div className={styles.rowIcon}>{icon}</div>
                    <div>
                      <div className={styles.rowName}>{user.name}</div>
                      <div className={styles.rowSub}>
                        {user.computerSessions?.length
                          ? `${user.computerSessions.length} сессий · ${user.hours ?? 0}ч`
                          : 'Нет сессий'}
                      </div>
                    </div>
                  </div>
                  <span className={`${styles.rowStat} ${styles.colRight}`}>{getStatValue(user)}</span>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

export default LeaderboardPage;
