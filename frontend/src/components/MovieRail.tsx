import type { Movie } from "../types";
import MovieCard from "./MovieCard";
import styles from "./MovieRail.module.css";

export default function MovieRail({ movies, emptyText }: { movies: Movie[]; emptyText: string }) {
  if (movies.length === 0) {
    return <div className={styles.empty}>{emptyText}</div>;
  }
  return (
    <div className={styles.rail}>
      {movies.map((m, i) => (
        <div key={m.id} className="fadeInUp" style={{ animationDelay: `${Math.min(i, 8) * 45}ms` }}>
          <MovieCard movie={m} />
        </div>
      ))}
    </div>
  );
}
