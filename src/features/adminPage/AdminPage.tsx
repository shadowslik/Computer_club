import React, { useState } from 'react';
import Sidebar from './components/Sidebar';
import MetricsPanel from './components/MetricsPanel';
import EndpointsTable from './components/EndpointsTable';
import RedisPanel from './components/RedisPanel';
import UsersTable from './components/UsersTable';
import BookingsTable from './components/BookingsTable';
import LiveLogs from './components/LiveLogs';
import styles from './AdminPage.module.css';

export type AdminSection = 'metrics' | 'endpoints' | 'redis' | 'users' | 'bookings' | 'logs';

const AdminPage: React.FC = () => {
  const [active, setActive] = useState<AdminSection>('metrics');

  const content: Record<AdminSection, React.ReactNode> = {
    metrics:   <MetricsPanel />,
    endpoints: <EndpointsTable />,
    redis:     <RedisPanel />,
    users:     <UsersTable />,
    bookings:  <BookingsTable />,
    logs:      <LiveLogs />,
  };

  return (
    <div className={styles.adminLayout}>
      <Sidebar activeSection={active} onSelect={setActive} />
      <main className={styles.adminMain}>
        {content[active]}
      </main>
    </div>
  );
};

export default AdminPage;
