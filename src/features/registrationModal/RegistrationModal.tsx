import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaTimes } from 'react-icons/fa';
import { useAuth } from '../../contexts/AuthContext'; // путь под ваш проект
import styles from './RegistrationModal.module.css';

interface RegistrationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSwitchToLogin: () => void;
}

const RegistrationModal: React.FC<RegistrationModalProps> = ({ isOpen, onClose, onSwitchToLogin }) => {
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const { register } = useAuth();   // получаем функцию регистрации из контекста
  const navigate = useNavigate();

  if (!isOpen) return null;

  const handleContentClick = (e: React.MouseEvent) => e.stopPropagation();

  const handleSwitch = (e: React.MouseEvent) => {
    e.preventDefault();
    onClose();
    onSwitchToLogin();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await register(name, email, phone, password);
      onClose();              // закрываем модальное окно
      navigate('/profile');   // переходим в личный кабинет
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modalContent} onClick={handleContentClick}>
        <button className={styles.closeModal} onClick={onClose}>
          <FaTimes />
        </button>

        <h2>РЕГИСТРАЦИЯ</h2>
        <p>Присоединяйся к сообществу лучших игроков</p>

        {error && <div className={styles.errorMessage}>{error}</div>}

        <form className={styles.authForm} onSubmit={handleSubmit}>
          <div className={styles.inputGroup}>
            <label>Ваш никнейм</label>
            <input
              type="text"
              placeholder="Gamer_777"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label>Телефон</label>
            <input
              type="tel"
              placeholder="+7 (900) 000-00-00"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label>Почта</label>
            <input
              type="email"
              placeholder="example@mail.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label>Пароль</label>
            <input
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit" className={styles.submitBtn} disabled={isLoading}>
            {isLoading ? 'Регистрация...' : 'Создать аккаунт'}
          </button>
        </form>

        <div className={styles.modalFooter}>
          Уже есть аккаунт?{' '}
          <span className={styles.link} onClick={handleSwitch}>
            Войти
          </span>
        </div>
      </div>
    </div>
  );
};

export default RegistrationModal;