import { useState, useEffect } from 'react';

interface IPCheckResponse {
  ip: string;
  list: 'white' | 'black' | 'gray' | 'none';
  message: string;
}


async function getMyIP(): Promise<string> {
  const response = await fetch('http://localhost:8080/ip/check');
  if (!response.ok) throw new Error('Не удалось получить IP');
  const data: { ip: string } = await response.json();
  return data.ip;
}

async function checkIPOnServer(ip: string): Promise<IPCheckResponse> {
  const response = await fetch(`http://localhost:8080/ip/check?ip=${encodeURIComponent(ip)}`);
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Ошибка HTTP ${response.status}`);
  }
  return response.json();
}

export function useIPAccess() {
  const [loading, setLoading] = useState(true);
  const [listType, setListType] = useState<'white' | 'black' | 'gray' | 'none' | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    const checkAccess = async () => {
      try {
        const ip = await getMyIP();              
        const result = await checkIPOnServer(ip); 
        if (isMounted) {
          setListType(result.list);
          setError(null);
        }
      } catch (err) {
        if (isMounted) {
          setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
          setListType(null);
        }
      } finally {
        if (isMounted) setLoading(false);
      }
    };
    checkAccess();
    return () => { isMounted = false; };
  }, []);

  const isAllowed = listType === 'white';
  return { loading, isAllowed, listType, error };
}