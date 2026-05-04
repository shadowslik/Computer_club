import React from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './GameList.module.css';
import valorantLogo from './gameLogs/valorant.jpg';

interface Game {
  id: number;
  title: string;
  genre: string;
  image: string;
}

const games: Game[] = [
  { id: 1, title: 'Counter-Strike 2',  genre: 'Tactical Shooter', image: 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/730/library_600x900_2x.jpg' },
  { id: 2, title: 'Dota 2',            genre: 'MOBA',             image: 'https://cdn.akamai.steamstatic.com/apps/dota2/images/dota2_social.jpg' },
  { id: 3, title: 'Valorant',          genre: 'Tactical Shooter', image: valorantLogo },
  { id: 4, title: 'Cyberpunk 2077',    genre: 'Action RPG',       image: 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/1091500/library_600x900_2x.jpg' },
  { id: 5, title: 'Apex Legends',      genre: 'Battle Royale',    image: 'https://media.contentapi.ea.com/content/dam/apex-legends/images/2019/01/apex-featured-image-16x9.jpg.adapt.crop191x100.1200w.jpg' },
  { id: 6, title: 'League of Legends', genre: 'MOBA',             image: 'https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Kaisa_0.jpg' },
];

const GameList: React.FC = () => {
  const navigate = useNavigate();

  const handleErr = (e: React.SyntheticEvent<HTMLImageElement>) => {
    e.currentTarget.src = 'https://via.placeholder.com/600x900/111/444?text=Game';
  };

  return (
    <section className={styles.gamesSection} id="games">
      <div className={styles.container}>
        <p className={styles.sectionEyebrow}>02 — Библиотека</p>
        <h2 className={styles.sectionTitle}>
          Библиотека игр <span>/ Games</span>
        </h2>
      </div>

      <div className={styles.gamesGrid}>
        {games.map((game, idx) => (
          <div key={game.id} className={styles.gameCard}>
            <div className={styles.gameImageWrapper}>
              <img
                src={game.image}
                alt={game.title}
                className={styles.gameImage}
                loading="lazy"
                onError={handleErr}
              />
              <div className={styles.gameOverlay}>
                <div className={styles.gameNumber}>{String(idx + 1).padStart(2, '0')}</div>
                <span className={styles.gameGenre}>{game.genre}</span>
                <h3 className={styles.gameName}>{game.title}</h3>
                <div className={styles.gameLine} />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className={styles.gamesFooter}>
        <div className={styles.lbTeaser} onClick={() => navigate('/leaderboard')}>
          <div className={styles.lbTeaserLeft}>
            <div className={styles.lbTeaserTitle}>Таблица рейтинга</div>
            <div className={styles.lbTeaserSub}>Топ игроков клуба за этот месяц</div>
          </div>
          <div className={styles.lbTeaserArrow}>Смотреть →</div>
        </div>

        <div className={styles.gamesMore}>
          <p>И ещё более 100 игр в Steam, Epic Games и Battle.net</p>
        </div>
      </div>
    </section>
  );
};

export default GameList;
