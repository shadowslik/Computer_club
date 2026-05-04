export interface Computer {
  id: number;
  type: string;
  status: string;
  pricePerHour: number;
}

export interface ResponseComputerSession {
  id: number;
  userName: string;
  userPhone: string;
  computer: Computer;
  start: string;   // ISO
  end: string;
  total: number;
}

// Для Create React App (webpack) используем process.env
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api';

export const fetchUserSessions = async (userId: number): Promise<ResponseComputerSession[]> => {
  const response = await fetch(`${API_BASE_URL}/computer_sessions/userSessions/${userId}`);
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Ошибка ${response.status}: ${errorText}`);
  }
  const data: ResponseComputerSession[] = await response.json();
  // Сортировка от новых к старым (на случай, если бэк вернёт не по порядку)
  return data.sort((a, b) => new Date(b.start).getTime() - new Date(a.start).getTime());
};