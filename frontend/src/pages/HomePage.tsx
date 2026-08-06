import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as moviesApi from "../api/movies";
import { useLocation as useCinemaLocation } from "../context/LocationContext";
import MovieRail from "../components/MovieRail";
import type { Movie } from "../types";
import styles from "./HomePage.module.css";

export default function HomePage() {
  const { selectedCity } = useCinemaLocation();
  const [nowShowing, setNowShowing] = useState<Movie[] | null>(null);
  const [comingSoon, setComingSoon] = useState<Movie[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([moviesApi.listMovies("now_showing"), moviesApi.listMovies("coming_soon")])
      .then(([now, soon]) => {
        setNowShowing(now ?? []);
        setComingSoon(soon ?? []);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load movies"));
  }, []);

  return (
    <div className="page container fadeIn">
      <section className={styles.hero}>
        <h1 className={styles.heroTitle}>Book your next big-screen moment.</h1>
        <p className={styles.heroSubtitle}>
          Browse what's playing now and what's coming soon, pick your seats, and get an instant receipt.
        </p>
        <div className={styles.heroActions}>
          <Link to="/select-location" className="btn btnPrimary">
            {selectedCity ? `Change city (${selectedCity.name})` : "Choose your city"}
          </Link>
        </div>
      </section>

      {!selectedCity && (
        <div className={`card ${styles.locationBanner} fadeInUp`}>
          <div>
            <strong>Pick a city to see showtimes near you.</strong>
            <div style={{ color: "var(--color-text-muted)", fontSize: 13, marginTop: 4 }}>
              You can also use your current location for a one-tap match.
            </div>
          </div>
          <Link to="/select-location" className="btn btnSecondary">
            Select location
          </Link>
        </div>
      )}

      {error && <p className="errorText">{error}</p>}

      <h2 className="sectionTitle">Now Showing</h2>
      {nowShowing ? (
        <MovieRail movies={nowShowing} emptyText="No movies are showing right now." />
      ) : (
        <RailSkeleton />
      )}

      <h2 className="sectionTitle">Coming Soon</h2>
      {comingSoon ? (
        <MovieRail movies={comingSoon} emptyText="No upcoming movies yet." />
      ) : (
        <RailSkeleton />
      )}
    </div>
  );
}

function RailSkeleton() {
  return (
    <div style={{ display: "flex", gap: 18 }}>
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="skeleton" style={{ width: 170, aspectRatio: "2 / 3", borderRadius: 16 }} />
      ))}
    </div>
  );
}
