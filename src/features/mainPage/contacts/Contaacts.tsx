import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  FaTelegramPlane, FaVk,
  FaMapMarkerAlt, FaPhoneAlt, FaEnvelope, FaClock,
} from 'react-icons/fa';
import styles from './Contacts.module.css';

const Contacts: React.FC = () => {
  const navigate = useNavigate();
  const address = 'Москва, ул. Пушкина';
  const mapSrc = `https://yandex.ru/map-widget/v1/?um=constructor%3A&source=constructor&mode=search&text=${encodeURIComponent(address)}&z=17`;
  const mapLink = `https://yandex.ru/maps/?text=${encodeURIComponent(address)}`;

  return (
    <>
      <section className={styles.contactsSection} id="contacts">
        <div className={styles.container}>
          <p className={styles.sectionEyebrow}>03 — Контакты</p>
          <h2 className={styles.sectionTitle}>
            Найди нас <span>/ Find Us</span>
          </h2>

          <div className={styles.contactsGrid}>

            {/* ── Info ── */}
            <div>
              <div className={styles.infoBlock}>
                <div className={styles.infoLabel}>
                  <FaMapMarkerAlt className={styles.infoIcon} /> Адрес
                </div>
                <div className={styles.infoValue}>{address}</div>
              </div>

              <div className={styles.infoBlock}>
                <div className={styles.infoLabel}>
                  <FaPhoneAlt className={styles.infoIcon} /> Телефон
                </div>
                <div className={styles.infoValue}>
                  <a href="tel:+74951234567">+7 (495) 123-45-67</a>
                  <a href="tel:+78005553535">8 (800) 555-35-35 (бесплатно)</a>
                </div>
              </div>

              <div className={styles.infoBlock}>
                <div className={styles.infoLabel}>
                  <FaEnvelope className={styles.infoIcon} /> Email
                </div>
                <div className={styles.infoValue}>
                  <a href="mailto:info@cyberlounge.ru">info@cyberlounge.ru</a>
                </div>
              </div>

              <div className={styles.infoBlock}>
                <div className={styles.infoLabel}>
                  <FaClock className={styles.infoIcon} /> Режим работы
                </div>
                <div className={styles.infoValue}>Ежедневно: 10:00 – 02:00</div>
              </div>

              <div className={styles.infoBlock}>
                <div className={styles.infoLabel}>Мы в соцсетях</div>
                <div className={styles.socialLinks}>
                  <a href="https://t.me/yourclub" target="_blank" rel="noopener noreferrer" aria-label="Telegram">
                    <FaTelegramPlane />
                  </a>
                  <a href="https://vk.com/yourclub" target="_blank" rel="noopener noreferrer" aria-label="VK">
                    <FaVk />
                  </a>
                </div>
              </div>
            </div>

            {/* ── Map ── */}
            <div>
              <div className={styles.mapHeading}>Как нас найти</div>
              <div className={styles.mapContainer}>
                <iframe
                  title="Карта клуба"
                  src={mapSrc}
                  allowFullScreen
                  loading="lazy"
                />
              </div>
              <p className={styles.mapNote}>
                <a href={mapLink} target="_blank" rel="noopener noreferrer">
                  Открыть в Яндекс.Картах →
                </a>
              </p>
            </div>

          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className={styles.siteFooter}>
        <div className={styles.footerLogo} onClick={() => navigate('/')}>
          CYBER<span>LOUNGE</span>
        </div>
        <div className={styles.footerCopy}>© 2025 CyberLounge · All rights reserved</div>
        <div className={styles.footerLinks}>
          <a href="https://t.me/yourclub" target="_blank" rel="noopener noreferrer">Telegram</a>
          <a href="https://vk.com/yourclub" target="_blank" rel="noopener noreferrer">VK</a>
          <a href="https://discord.gg/yourclub" target="_blank" rel="noopener noreferrer">Discord</a>
        </div>
      </footer>
    </>
  );
};

export default Contacts;
