import { Link } from "react-router-dom";
import type { Movie } from "../types";
import styles from "./MovieCard.module.css";

export default function MovieCard({ movie }: { movie: Movie }) {
  return (
    <Link to={`/movies/${movie.id}`} className={styles.card}>
      <div className={styles.posterWrap}>
        {movie.poster_url && <img className={styles.poster} src={movie.poster_url} alt={movie.title} loading="lazy" />}
        {movie.status === "coming_soon" && <span className={styles.badge}>Coming Soon</span>}
      </div>
      <div className={styles.title}>{movie.title}</div>
      <div className={styles.meta}>{movie.genres.slice(0, 2).join(", ")}</div>
    </Link>
  );
}
