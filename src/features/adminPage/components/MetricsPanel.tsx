import React, { useEffect, useState, useCallback } from 'react';
import styles from '../AdminPage.module.css';

interface Metrics {
  rps: number; latency: number; cpu: number; ram: number; goroutines: number;
}

const randomBetween = (min: number, max: number) =>
  Math.floor(Math.random() * (max - min + 1)) + min;

const MetricsPanel: React.FC = () => {
  const [metrics, setMetrics] = useState<Metrics>({
    rps: 342, latency: 18, cpu: 24, ram: 312, goroutines: 148,
  });
  const [clock, setClock] = useState('');
  const [bars, setBars] = useState<number[]>(
    Array.from({ length: 48 }, () => randomBetween(25, 100))
  );

  useEffect(() => {
    const tick = () => setClock(new Date().toLocaleTimeString('ru-RU'));
    tick();
    const t = setInterval(tick, 1000);
    return () => clearInterval(t);
  }, []);

  useEffect(() => {
    const t = setInterval(() => {
      setMetrics({
        rps: randomBetween(270, 420),
        latency: randomBetween(10, 28),
        cpu: randomBetween(14, 38),
        ram: randomBetween(275, 360),
        goroutines: randomBetween(100, 180),
      });
      setBars(prev => [...prev.slice(1), randomBetween(25, 100)]);
    }, 2500);
    return () => clearInterval(t);
  }, []);

  return (
    <div className={styles.section}>
      <div className={styles.topbar}>
        <div className={styles.topbarLeft}>
          <div className={styles.breadcrumb}>Dashboard / Метрики</div>
          <h2 className={styles.pageTitle}>Метрики системы</h2>
        </div>
        <div className={styles.clock}>{clock}</div>
      </div>

      {/* Primary metrics */}
      <div className={styles.metricsGrid}>
        <MetCard label="RPS" value={metrics.rps} unit="req/s" accent="Green" color="clrGreen" trend="trendUp" trendText={`↑ ${metrics.rps > 350 ? '+14%' : '+8%'} vs avg`} />
        <MetCard label="P99 Latency"  value={`${metrics.latency}ms`} unit="" accent="Blue"   color="clrBlue"   trend="trendUp"   trendText="↓ улучшение -3ms" />
        <MetCard label="CPU Usage"    value={`${metrics.cpu}%`}      unit="" accent="Yellow" color="clrYellow" trend="trendFlat" trendText="— стабильно" />
        <MetCard label="RAM"          value={`${metrics.ram}MB`}     unit="" accent="Gray"   color=""          trend="trendUp"   trendText="↓ 15% heap free" />
      </div>

      {/* RPS Chart */}
      <div className={styles.chartArea}>
        <div className={styles.chartHeader}>
          <span className={styles.chartTitle}>RPS — последние 24 часа</span>
          <div className={styles.chartLegend}>
            <span><i className={styles.dotGreen} />Запросы</span>
            <span><i className={styles.dotRed} />Ошибки</span>
          </div>
        </div>
        <div className={styles.chartBars}>
          {bars.map((h, i) => (
            <div
              key={i}
              className={styles.bar}
              style={{ height: `${h}%`, background: `rgba(0,232,122,${0.25 + (h / 100) * 0.55})` }}
            >
              <div className={styles.barTip}>{Math.round(h * 4)} rps</div>
            </div>
          ))}
        </div>
      </div>

      {/* Secondary */}
      <div className={styles.metricsGrid3}>
        <MetCard label="Cache Hit Rate"    value="87%"                    unit="" accent="Green" color="clrGreen" />
        <MetCard label="Goroutines"        value={metrics.goroutines}     unit="" accent="Blue"  color="clrBlue" />
        <MetCard label="Uptime"            value="14д 6ч"                 unit="" accent="Gray"  color="" />
      </div>
    </div>
  );
};

interface MetCardProps {
  label: string; value: string | number; unit: string;
  accent: string; color: string; trend?: string; trendText?: string;
}
const MetCard: React.FC<MetCardProps> = ({ label, value, unit, accent, color, trend, trendText }) => (
  <div className={`${styles.metricCard} ${styles[`accent${accent}`]}`}>
    <div className={styles.metricLabel}>{label}</div>
    <div className={`${styles.metricValue} ${color ? styles[color] : ''}`}>
      {value}{unit && <span className={styles.metricUnit}>{unit}</span>}
    </div>
    {trendText && (
      <div className={`${styles.metricTrend} ${trend ? styles[trend] : ''}`}>{trendText}</div>
    )}
  </div>
);

export default MetricsPanel;
