import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import styles from './TablesPage.module.css';

interface Computer {
  id: number;
  type: string;
  status: string;
  pricePerHour: number;
}

const TablesPage: React.FC = () => {
  const navigate = useNavigate();
  const [computers, setComputers] = useState<Computer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('http://localhost:8080/api/computers')
      .then(res => {
        if (!res.ok) throw new Error('Ошибка загрузки');
        return res.json();
      })
      .then(data => {
        // Бэкенд возвращает массив, а не объект { content: [...] }
        const computersList = Array.isArray(data) ? data : [];
        setComputers(computersList);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError(err.message);
        setLoading(false);
      });
  }, []);

  const handleBookTable = (computer: Computer) => {
    if (computer.status === 'Не занят') {
      navigate(`/booking?computerId=${computer.id}`);
    } else {
      alert('Этот компьютер сейчас занят');
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type?.toLowerCase()) {
      case 'standard': return '◈';
      case 'pro': return '❖';
      case 'vip': return '🔱';
      default: return '🖥️';
    }
  };

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'Не занят': return styles.available;
      case 'Занят': return styles.occupied;
      default: return styles.soon;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'Не занят': return 'Свободен';
      case 'Занят': return 'Занят';
      default: return status;
    }
  };

  if (loading) return <div className={styles.loader}>Загрузка...</div>;
  if (error) return <div className={styles.error}>Ошибка: {error}</div>;

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.container}>
        <Link to="/" className={styles.backLink}>← На главную</Link>
        <h1 className={styles.pageTitle}>
          Выберите компьютер <span>/ Select PC</span>
        </h1>
        <div className={styles.tablesGrid}>
          {computers.map(computer => (
            <div key={computer.id} className={styles.tableCard} onClick={() => handleBookTable(computer)}>
              <div className={styles.tableIcon}>{getTypeIcon(computer.type)}</div>
              <h3 className={styles.tableType}>
                {computer.type?.toUpperCase() || 'PC'} #{computer.id}
              </h3>
              <ul className={styles.tableSpecs}>
                <li>{computer.pricePerHour}₽ / час</li>
              </ul>
              <span className={`${styles.tableStatus} ${getStatusClass(computer.status)}`}>
                {getStatusText(computer.status)}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TablesPage;