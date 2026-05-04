import { useState, useEffect } from 'react';
import { fetchUserSessions, ResponseComputerSession } from '../api/userSessions';

export const useUserSessions = (userId: number | undefined) => {
  const [sessions, setSessions] = useState<ResponseComputerSession[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!userId) {
      setLoading(false);
      return;
    }

    let cancelled = false;
    const load = async () => {
      setLoading(true);
      try {
        const data = await fetchUserSessions(userId);
        if (!cancelled) {
          setSessions(data);
          setError(null);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Ошибка загрузки');
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => { cancelled = true; };
  }, [userId]);

  return { sessions, loading, error };
};