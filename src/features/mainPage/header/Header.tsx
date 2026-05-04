import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaTelegramPlane, FaVk, FaUserAlt } from 'react-icons/fa';
import { useAuth } from '../../../contexts/AuthContext';
import styles from './Header.module.css';

interface HeaderProps {
  onLoginClick: () => void;
}

const Header: React.FC<HeaderProps> = ({ onLoginClick }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 50);
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  const scrollTo = (e: React.MouseEvent<HTMLAnchorElement>, id: string) => {
    e.preventDefault();
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleLogout = () => {
    if (window.confirm('Вы уверены, что хотите выйти?')) {
      logout();
      navigate('/');
    }
  };

  return (
    <header className={`${styles.header} ${scrolled ? styles.scrolled : ''}`}>
      <div className={styles.headerInner}>

        {/* Logo */}
        <div className={styles.logoSection}>
          <span className={styles.logoLink} onClick={() => navigate('/')}>
            CYBER<span>LOUNGE</span>
          </span>
        </div>

        {/* Nav */}
        <nav className={styles.navMenu} aria-label="Главное меню">
          <ul>
            <li><a href="#prices"   onClick={e => scrollTo(e, 'prices')}>Цены</a></li>
            <li><a href="#games"    onClick={e => scrollTo(e, 'games')}>Игры</a></li>
            <li><a href="#"         onClick={e => { e.preventDefault(); navigate('/leaderboard'); }}>Рейтинг</a></li>
            <li><a href="#contacts" onClick={e => scrollTo(e, 'contacts')}>Контакты</a></li>
          </ul>
        </nav>

        {/* Right */}
        <div className={styles.headerRight}>
          <div className={styles.socials}>
            <a href="https://t.me/yourclub"  target="_blank" rel="noopener noreferrer" aria-label="Telegram">
              <FaTelegramPlane />
            </a>
            <a href="https://vk.com/yourclub" target="_blank" rel="noopener noreferrer" aria-label="VK">
              <FaVk />
            </a>
          </div>

          {user ? (
            <div className={styles.userMenu}>
              <div className={styles.userChip} onClick={() => navigate('/profile')}>
                <span className={styles.onlineDot} />
                <FaUserAlt className={styles.userIcon} />
                {user.name}
              </div>
              <button className={styles.logoutBtn} onClick={handleLogout}>
                Выйти
              </button>
            </div>
          ) : (
            <button className={styles.loginBtn} onClick={onLoginClick}>
              <FaUserAlt className={styles.userIcon} />
              Личный кабинет
            </button>
          )}
        </div>

      </div>
    </header>
  );
};

export default Header;
