import { ResponseComputerSession } from '../api/userSessions';

export const formatSession = (session: ResponseComputerSession) => {
  const startDate = new Date(session.start);
  const endDate = new Date(session.end);
  const durationHours = (endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60);
  const hoursRounded = Math.round(durationHours * 10) / 10;

  let icon = '◈';
  const typeLower = session.computer.type.toLowerCase();
  if (typeLower.includes('pro')) icon = '❖';
  else if (typeLower.includes('vip')) icon = '🔱';

  return {
    id: session.id,
    tableName: `${session.computer.type} #${session.computer.id}`,
    date: startDate.toLocaleDateString('ru-RU'),
    hours: hoursRounded,
    amount: session.total,
    icon,
  };
};