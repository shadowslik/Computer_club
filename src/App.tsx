import React, { useState } from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom';
import HomePage from './features/mainPage/previe/HomePage';
import TablesPage from './features/tablesPage/TablesPage';
import BookingPage from './features/bookingPage/BookingPage';
import RegistrationModal from './features/registrationModal/RegistrationModal';
import LoginModal from './features/loginModal/LoginModal';
import { ForbiddenPage } from './features/forbiddenPages/ForbiddenPage';
import { AuthProvider } from './contexts/AuthContext';
import ProfilePage from './features/profilePage/ProfilePage';
import ProtectedRoute from './features/ProtectedRoute';
import LeaderboardPage from './features/leaderboardPage/LeaderboardPage';



import AdminPage from './features/adminPage/AdminPage';

type ModalType = 'register' | 'login' | null;

const AppContent: React.FC = () => {
  const [modalType, setModalType] = useState<ModalType>(null);
  const navigate = useNavigate();

  const openLogin = () => setModalType('login');
  const closeModal = () => setModalType(null);
  const switchToRegister = () => setModalType('register');
  const switchToLogin = () => setModalType('login');

  const handleBookNow = () => {
    navigate('/tables');
  };

  return (
    <>
      <Routes>
        <Route
          path="/"
          element={
            <HomePage
              onLoginClick={openLogin}
              onBookNowClick={handleBookNow}
            />
          }
        />
        <Route path="/forbidden" element={<ForbiddenPage />} />
        <Route path="/tables" element={<TablesPage />} />
        <Route path="/booking" element={<BookingPage />} />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />
        {/* Админ-панель — только для пользователей с ролью admin */}
          <Route
            path="/admin"
            element={
              <ProtectedRoute>
                <AdminPage />
              </ProtectedRoute>
            }
          />
          <Route path="/leaderboard" element={<LeaderboardPage />} />
      </Routes>

      <RegistrationModal
        isOpen={modalType === 'register'}
        onClose={closeModal}
        onSwitchToLogin={switchToLogin}
      />
      <LoginModal
        isOpen={modalType === 'login'}
        onClose={closeModal}
        onSwitchToRegister={switchToRegister}
      />
    </>
  );
};

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;