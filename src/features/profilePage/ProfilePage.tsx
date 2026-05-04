import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { useUserSessions } from '../../hooks/useUserSessions';
import { formatSession } from '../../utils/formatSession';
import styles from './ProfilePage.module.css';

const API_BASE = 'http://localhost:8080/api/clients';

// ─── Типы ────────────────────────────────────────────────────────────────────
interface UserStats {
  balance: number;
  hours: number;
  totalBookings: number;
}

interface EditForm {
  name: string;
  email: string;
  phone: string;
  password: string;
}

// ─── Компонент ───────────────────────────────────────────────────────────────
const ProfilePage: React.FC = () => {
  const { user, logout, loading, updateUser } = useAuth();
  const navigate = useNavigate();
  const { sessions, loading: sessionsLoading, error: sessionsError } = useUserSessions(user?.id);

  // Статистика с сервера
  const [stats, setStats] = useState<UserStats | null>(null);
  const [statsLoading, setStatsLoading] = useState(true);
  const [statsError, setStatsError] = useState<string | null>(null);

  // Режим редактирования
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState<EditForm>({ name: '', email: '', phone: '', password: '' });
  const [editLoading, setEditLoading] = useState(false);
  const [editError, setEditError] = useState<string | null>(null);
  const [editSuccess, setEditSuccess] = useState(false);

  // Загрузка статистики по /api/clients/{id}
  useEffect(() => {
    if (!user?.id) return;
    let cancelled = false;
    setStatsLoading(true);
    setStatsError(null);

    fetch(`${API_BASE}/${user.id}`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(data => {
        if (!cancelled) {
          setStats({
            balance: data.balance ?? 0,
            hours: data.hours ?? 0,
            totalBookings: data.computerSessions?.length ?? 0,
          });
        }
      })
      .catch(err => {
        if (!cancelled) setStatsError(err.message);
      })
      .finally(() => {
        if (!cancelled) setStatsLoading(false);
      });

    return () => { cancelled = true; };
  }, [user?.id]);

  // Открыть форму редактирования
  const handleOpenEdit = () => {
    setEditForm({
      name: user?.name ?? '',
      email: user?.email ?? '',
      phone: user?.phone ?? '',
      password: '',
    });
    setEditError(null);
    setEditSuccess(false);
    setIsEditing(true);
  };

  // Закрыть форму
  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditError(null);
    setEditSuccess(false);
  };

  // Сохранить изменения → PUT /api/clients?id=
  const handleSaveEdit = async () => {
    if (!user?.id) return;
    setEditLoading(true);
    setEditError(null);
    setEditSuccess(false);

    // Отправляем только заполненные поля
    const body: Record<string, string> = {};
    if (editForm.name.trim())     body.name     = editForm.name.trim();
    if (editForm.email.trim())    body.email    = editForm.email.trim();
    if (editForm.phone.trim())    body.phone    = editForm.phone.trim();
    if (editForm.password.trim()) body.password = editForm.password.trim();

    try {
      const res = await fetch(`${API_BASE}?id=${user.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `HTTP ${res.status}`);
      }

      const updated = await res.json();

      // Обновляем контекст если есть метод updateUser
      if (typeof updateUser === 'function') {
        updateUser(updated);
      }

      setEditSuccess(true);
      setTimeout(() => {
        setIsEditing(false);
        setEditSuccess(false);
      }, 1500);
    } catch (err: any) {
      setEditError(err.message ?? 'Ошибка при сохранении');
    } finally {
      setEditLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  // Ранг по балансу
  const getUserRank = () => {
    const balance = stats?.balance ?? user?.balance ?? 0;
    if (balance >= 5000) return { label: 'VIP', className: styles.rankVip };
    if (balance >= 2000) return { label: 'PRO', className: styles.rankPro };
    return { label: 'НОВИЧОК', className: styles.rankNewbie };
  };
  const rank = getUserRank();
  const userInitials = user?.name?.slice(0, 2).toUpperCase() || '??';

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.loader}>Загрузка...</div>
      </div>
    );
  }

  if (!user) {
    navigate('/');
    return null;
  }

  return (
    <div className={styles.profilePage}>
      <div className={styles.profileBgText}>PROFILE</div>
      <div className={styles.profileInner}>

        {/* ── Верхняя секция ── */}
        <div className={styles.profileTop}>
          <div className={styles.profileAvatarWrap}>
            <div className={styles.profileAvatar}>
              {userInitials}
              <div className={styles.avatarRing} />
            </div>
          </div>
          <div className={styles.profileMeta}>
            <h1 className={styles.profileName}>{user.name}</h1>
            <div className={styles.profileRank}>
              <span>Участник клуба</span>
              <span className={`${styles.rankBadge} ${rank.className}`}>{rank.label}</span>
            </div>
          </div>
          <div className={styles.profileBackBtn}>
            <button className={styles.btnProfileBack} onClick={() => navigate('/')}>
              ← На главную
            </button>
            <button className={styles.btnProfileBack} onClick={() => navigate('/leaderboard')}>
              ◈ Рейтинг
            </button>
          </div>
        </div>

        {/* ── Статистика (с сервера) ── */}
        <div className={styles.profileStatsRow}>
          <div className={styles.statCard}>
            <span className={styles.statNumber}>
              {statsLoading ? '—' : (stats?.totalBookings ?? 0)}
            </span>
            <span className={styles.statLabel}>Бронирований</span>
          </div>
          <div className={styles.statCard}>
            <span className={`${styles.statNumber} ${styles.textGreen}`}>
              {statsLoading ? '—' : `${stats?.balance ?? 0} ₽`}
            </span>
            <span className={styles.statLabel}>Баланс</span>
          </div>
          <div className={styles.statCard}>
            <span className={styles.statNumber}>
              {statsLoading ? '—' : `${stats?.hours ?? 0}ч`}
            </span>
            <span className={styles.statLabel}>Часов в клубе</span>
          </div>
        </div>

        {statsError && (
          <div className={styles.errorMessage} style={{ marginBottom: 24 }}>
            Не удалось загрузить статистику: {statsError}
          </div>
        )}

        {/* ── Две колонки ── */}
        <div className={styles.profileGrid}>

          {/* Личные данные / форма редактирования */}
          <div className={styles.profileBlock}>
            <h3 className={styles.blockTitle}>Личные данные</h3>

            {!isEditing ? (
              // ── VIEW MODE ──
              <>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>Никнейм</span>
                  <span className={styles.infoValue}>{user.name}</span>
                </div>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>Email</span>
                  <span className={styles.infoValue}>{user.email}</span>
                </div>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>Телефон</span>
                  <span className={styles.infoValue}>{user.phone || '—'}</span>
                </div>
              </>
            ) : (
              // ── EDIT MODE ──
              <div className={styles.editForm}>
                {editSuccess && (
                  <div className={styles.successMessage}>✓ Данные успешно сохранены</div>
                )}
                {editError && (
                  <div className={styles.errorMessage}>{editError}</div>
                )}

                <div className={styles.editField}>
                  <label className={styles.editLabel}>Никнейм</label>
                  <input
                    className={styles.editInput}
                    name="name"
                    value={editForm.name}
                    onChange={handleChange}
                    placeholder="Ваш никнейм"
                    autoComplete="off"
                  />
                </div>
                <div className={styles.editField}>
                  <label className={styles.editLabel}>Email</label>
                  <input
                    className={styles.editInput}
                    name="email"
                    type="email"
                    value={editForm.email}
                    onChange={handleChange}
                    placeholder="example@mail.com"
                    autoComplete="off"
                  />
                </div>
                <div className={styles.editField}>
                  <label className={styles.editLabel}>Телефон</label>
                  <input
                    className={styles.editInput}
                    name="phone"
                    type="tel"
                    value={editForm.phone}
                    onChange={handleChange}
                    placeholder="+7 (900) 000-00-00"
                  />
                </div>
                <div className={styles.editField}>
                  <label className={styles.editLabel}>Новый пароль</label>
                  <input
                    className={styles.editInput}
                    name="password"
                    type="password"
                    value={editForm.password}
                    onChange={handleChange}
                    placeholder="Оставьте пустым чтобы не менять"
                    autoComplete="new-password"
                  />
                </div>

                <div className={styles.editActions}>
                  <button
                    className={styles.editSaveBtn}
                    onClick={handleSaveEdit}
                    disabled={editLoading}
                  >
                    {editLoading ? 'Сохранение...' : 'Сохранить'}
                  </button>
                  <button
                    className={styles.editCancelBtn}
                    onClick={handleCancelEdit}
                    disabled={editLoading}
                  >
                    Отмена
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* История бронирований */}
          <div className={styles.profileBlock}>
            <h3 className={styles.blockTitle}>Последние бронирования</h3>
            {sessionsLoading && (
              <div className={styles.loaderSmall}>Загрузка истории...</div>
            )}
            {sessionsError && (
              <div className={styles.errorMessage}>Ошибка: {sessionsError}</div>
            )}
            {!sessionsLoading && !sessionsError && sessions.length === 0 && (
              <div className={styles.emptyMessage}>У вас пока нет бронирований</div>
            )}
            {!sessionsLoading && sessions.map(session => {
              const { id, tableName, date, hours, amount, icon } = formatSession(session);
              return (
                <div key={id} className={styles.recentBooking}>
                  <span className={styles.bookingIcon}>{icon}</span>
                  <div className={styles.bookingInfo}>
                    <div className={styles.bookingName}>{tableName}</div>
                    <div className={styles.bookingDate}>{date} · {hours} ч</div>
                  </div>
                  <span className={styles.bookingSum}>{amount} ₽</span>
                </div>
              );
            })}
          </div>
        </div>

        {/* ── Кнопки действий ── */}
        <div className={styles.profileActions}>
          <button
            className={`${styles.actionBtn} ${styles.primary}`}
            onClick={() => navigate('/tables')}
          >
            Забронировать стол
          </button>
          <button
            className={`${styles.actionBtn} ${styles.outline}`}
            onClick={isEditing ? handleCancelEdit : handleOpenEdit}
          >
            {isEditing ? 'Отменить редактирование' : 'Редактировать профиль'}
          </button>
          <button
            className={`${styles.actionBtn} ${styles.danger}`}
            onClick={() => { logout(); navigate('/'); }}
          >
            Выйти
          </button>
        </div>

      </div>
    </div>
  );
};

export default ProfilePage;