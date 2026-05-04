import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useIPAccess } from '../../../hooks/checkIp';
import styles from './PriceList.module.css';
import standartBg from './priceBG/standart.jpg';
import proBg from './priceBG/pro.jpg';
import vipBg from './priceBG/vip.jpg';

interface PriceOption {
  id: number;
  icon: string;
  label: string;
  title: string;
  price: string;
  duration: string;
  features: string[];
  backgroundImage: string;
  typeSlug: 'standard' | 'pro' | 'vip';
}

const prices: PriceOption[] = [
  {
    id: 1, icon: '◈', label: 'Тариф 01', title: 'Стандарт',
    price: '150₽', duration: 'за 1 час',
    features: ['RTX 3060 Ti', 'Monitor 144Hz · 1080p', 'Mechanical Keyboard', 'Gaming Headset'],
    backgroundImage: standartBg, typeSlug: 'standard',
  },
  {
    id: 2, icon: '❖', label: 'Тариф 02', title: 'Pro Zone',
    price: '250₽', duration: 'за 1 час',
    features: ['RTX 4070', 'Monitor 240Hz · 1440p', 'Pro Peripherals', 'Priority Support'],
    backgroundImage: proBg, typeSlug: 'pro',
  },
  {
    id: 3, icon: '🔱', label: 'Тариф 03', title: 'VIP Lounge',
    price: '400₽', duration: 'за 1 час',
    features: ['RTX 4090', 'Monitor 360Hz · 4K', 'Private Room', 'Personal Waiter'],
    backgroundImage: vipBg, typeSlug: 'vip',
  },
];

const PriceList: React.FC = () => {
  const navigate = useNavigate();
  const { isAllowed, loading } = useIPAccess();

  const handleBook = (typeSlug: string) => {
    if (loading) return;
    navigate(isAllowed ? `/booking?type=${typeSlug}` : '/forbidden');
  };

  return (
    <section className={styles.pricesSection} id="prices">
      <div className={styles.container}>
        <p className={styles.sectionEyebrow}>01 — Тарифы</p>
        <h2 className={styles.sectionTitle}>
          Выбери свой уровень
          <span>/ Pricing Plans</span>
        </h2>

        <div className={styles.priceGrid}>
          {prices.map(item => (
            <div
              key={item.id}
              className={styles.priceCard}
              style={{ '--card-bg-image': `url(${item.backgroundImage})` } as React.CSSProperties}
            >
              <span className={styles.cardIcon}>{item.icon}</span>
              <div className={styles.cardLabel}>{item.label}</div>
              <div className={styles.cardTitle}>{item.title}</div>

              <div className={styles.cardPriceWrap}>
                <div className={styles.cardPrice}>{item.price}</div>
                <div className={styles.cardDur}>{item.duration}</div>
              </div>

              <div className={styles.cardDivider} />

              <ul className={styles.cardFeatures}>
                {item.features.map((f, i) => <li key={i}>{f}</li>)}
              </ul>

              <button
                className={styles.bookBtn}
                onClick={() => handleBook(item.typeSlug)}
                disabled={loading}
              >
                Забронировать
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default PriceList;
