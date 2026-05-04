import React, { useState, useEffect } from 'react';
import { useSearchParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import styles from './BookingPage.module.css';

const BookingPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const computerId = searchParams.get('computerId');
  const { user, updateUser } = useAuth();
  const navigate = useNavigate();

  const [durationHours, setDurationHours] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [computerInfo, setComputerInfo] = useState<{ type: string; pricePerHour: number } | null>(null);
  const [loadingComputer, setLoadingComputer] = useState(true);

  useEffect(() => {
    if (!computerId) return;
    fetch(`http://localhost:8080/api/computers/id/${computerId}`)
      .then(res => {
        if (!res.ok) throw new Error('Компьютер не найден');
        return res.json();
      })
      .then(data => {
        setComputerInfo(data);
        setLoadingComputer(false);
      })
      .catch(err => {
        setError(err.message);
        setLoadingComputer(false);
      });
  }, [computerId]);

  const getTypeDisplay = (type: string) => {
    switch (type?.toLowerCase()) {
      case 'standard': return { name: 'Стандарт', icon: '◈' };
      case 'pro': return { name: 'Pro Zone', icon: '❖' };
      case 'vip': return { name: 'VIP Lounge', icon: '🔱' };
      default: return { name: type || 'Компьютер', icon: '🖥️' };
    }
  };

  const typeInfo = computerInfo ? getTypeDisplay(computerInfo.type) : { name: 'Компьютер', icon: '🖥️' };

  const totalPrice = computerInfo ? computerInfo.pricePerHour * durationHours : 0;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) {
      alert('Войдите в аккаунт');
      return;
    }
    setError('');
    setIsLoading(true);

    try {
      const response = await fetch('http://localhost:8080/api/computer_sessions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          computer_id: Number(computerId),
          user_id: user.id,
          durationHours: durationHours,
        }),
      });

      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Ошибка бронирования');
      }

      const sessionData = await response.json();
      // Обновляем пользователя
      const userResponse = await fetch(`http://localhost:8080/api/clients/${user.id}`);
      const updatedUser = await userResponse.json();
      if (updateUser) updateUser(updatedUser);

      alert(`Бронирование успешно! Списано: ${sessionData.total}₽`);
      navigate('/profile');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (!computerId) return <div className={styles.error}>Не указан ID компьютера</div>;
  if (loadingComputer) return <div className={styles.loader}>Загрузка...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.container}>
        <div className={styles.bookingCard}>
          <h2 className={styles.pageTitle}>Бронирование</h2>
          <p className={styles.pageSubtitle}>Выберите время игры</p>

          <div className={styles.typeIndicator}>
            <span className={styles.typeIcon}>{typeInfo.icon}</span>
            <span className={styles.typeName}>{typeInfo.name}</span>
            {computerInfo && (
              <span className={styles.price}>{computerInfo.pricePerHour}₽/час</span>
            )}
          </div>

          <form className={styles.form} onSubmit={handleSubmit}>
            <div className={styles.formGroup}>
              <label>Сколько часов?</label>
              <select
                value={durationHours}
                onChange={(e) => setDurationHours(Number(e.target.value))}
              >
                {[1, 2, 3, 4, 5, 6, 7, 8].map((h) => (
                  <option key={h} value={h}>
                    {h} час(а/ов)
                  </option>
                ))}
              </select>
            </div>

            {error && <div className={styles.errorMessage}>{error}</div>}

            <button type="submit" className={styles.submitBtn} disabled={isLoading}>
              {isLoading
                ? 'Бронирование...'
                : `Забронировать (${totalPrice}₽)`}
            </button>
          </form>

          <div className={styles.backLink}>
            <Link to="/tables">← Назад к выбору стола</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BookingPage;