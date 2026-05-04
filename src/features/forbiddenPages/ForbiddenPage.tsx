// ForbiddenPage.tsx
import { Link } from 'react-router-dom';
import { useIPAccess } from '../../hooks/checkIp';
import styles from './ForbiddenPage.module.css';

export const ForbiddenPage = () => {
  const {listType,loading, error } = useIPAccess();


  const getMessage = () => {
    if (loading) return '⏳ Проверка вашего IP...';
    if (error) return '⚠️ Ошибка соединения с сервером. Попробуйте позже.';
    switch (listType) {
      case 'black':
        return '🚫 Ваш IP в чёрном списке. Доступ запрещён навсегда.';
      case 'gray':
        return '⚙️ Ваш IP в сером списке. Требуется ручная верификация.';
      case 'none':
        return '❓ Ваш IP не найден в списках. Доступ только по белому списку.';
      default:
        return '🔒 Доступ запрещён. Ваш IP не в белом списке.';
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.errorCode}>403</div>
        <h1 className={styles.title}>Доступ запрещён</h1>
        <p className={styles.message}>{getMessage()}</p>
    

        <Link to="/" className={styles.homeLink}>
          ← Вернуться на главную
        </Link>
      </div>
    </div>
  );
};