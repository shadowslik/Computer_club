import React from 'react';
import Header from '../header/Header';
import PriceList from '../priceList/PriceList';
import GameList from '../gameList/GameList';
import Contacts from '../contacts/Contaacts';
import styles from './HomePage.module.css';

interface HomePageProps {
  onLoginClick: () => void;
  onBookNowClick: () => void;
}

const HomePage: React.FC<HomePageProps> = ({ onLoginClick, onBookNowClick }) => {
  return (
    <div className={styles.appWrapper}>
      <Header onLoginClick={onLoginClick} />

      <main>
        {/* ── HERO ── */}
        <section className={styles.heroContainer}>
          <div className={styles.heroBg} />
          <div className={styles.heroGradient} />
          <div className={styles.heroGrid} />
          <div className={styles.heroVlines}>
            <span /><span /><span /><span />
          </div>

          <div className={styles.heroContent}>
            <div className={styles.heroBadge}>
              Москва · Премиальный Gaming Club
            </div>

            <h1 className={styles.heroTitle}>
              ТВОЙ ПРЕДЕЛ
              <span className={styles.accent}>ТОЛЬКО ВЫСОКИЙ</span>
              FPS
            </h1>

            <p className={styles.heroSub}>
              Premium Gaming · 24 / 7 · Круглосуточно
            </p>

            <div className={styles.heroBtns}>
              <button className={styles.btnPrimary} onClick={onBookNowClick}>
                Забронировать стол
              </button>
              <button
                className={styles.btnOutline}
                onClick={() => {
                  document.getElementById('prices')?.scrollIntoView({ behavior: 'smooth' });
                }}
              >
                Смотреть тарифы
              </button>
            </div>
          </div>

          <div className={styles.heroScroll}>
            <div className={styles.scrollBar} />
            <span>Scroll</span>
          </div>
        </section>

        {/* ── STATS ── */}
        <div className={styles.statsBar}>
          <div className={styles.statItem}>
            <span className={styles.statNum}>48</span>
            <span className={styles.statLabel}>Игровых мест</span>
          </div>
          <div className={styles.statItem}>
            <span className={styles.statNum}>360</span>
            <span className={styles.statLabel}>Гц мониторы</span>
          </div>
          <div className={styles.statItem}>
            <span className={styles.statNum}>100+</span>
            <span className={styles.statLabel}>Игр в библиотеке</span>
          </div>
          <div className={styles.statItem}>
            <span className={styles.statNum}>24/7</span>
            <span className={styles.statLabel}>Без выходных</span>
          </div>
        </div>

        <PriceList />
        <GameList />
        <Contacts />
      </main>
    </div>
  );
};

export default HomePage;
