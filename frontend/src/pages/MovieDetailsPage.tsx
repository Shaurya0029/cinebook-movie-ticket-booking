import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import * as moviesApi from "../api/movies";
import * as showsApi from "../api/shows";
import { useLocation as useCinemaLocation } from "../context/LocationContext";
import { formatMoney, formatShowDate, formatShowTime } from "../utils/format";
import type { Movie, Show } from "../types";
import styles from "./MovieDetailsPage.module.css";

export default function MovieDetailsPage() {
  const { movieId } = useParams();
  const navigate = useNavigate();
  const { selectedCity } = useCinemaLocation();

  const [movie, setMovie] = useState<Movie | null>(null);
  const [shows, setShows] = useState<Show[] | null>(null);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!movieId) return;
    setMovie(null);
    setShows(null);
    setSelectedDate(null);
    moviesApi
      .getMovie(Number(movieId))
      .then(setMovie)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load movie"));
  }, [movieId]);

  useEffect(() => {
    if (!movieId || !movie || movie.status !== "now_showing") return;
    showsApi
      .listShowsForMovie(Number(movieId), { cityId: selectedCity?.id })
      .then((data) => setShows(data ?? []))
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load showtimes"));
  }, [movieId, movie, selectedCity]);

  const upcomingShows = useMemo(() => (shows ?? []).filter((s) => new Date(s.starts_at).getTime() > Date.now()), [shows]);
  const dateGroups = useMemo(() => groupByDate(upcomingShows), [upcomingShows]);
  const activeDate = selectedDate ?? dateGroups[0]?.[0] ?? null;
  const activeShows = dateGroups.find(([d]) => d === activeDate)?.[1] ?? [];
  const theaterGroups = useMemo(() => groupByTheater(activeShows), [activeShows]);

  if (error) {
    return (
      <div className="page container">
        <p className="errorText">{error}</p>
      </div>
    );
  }

  if (!movie) {
    return (
      <div className="page container">
        <div className="skeleton" style={{ height: 400, borderRadius: 16 }} />
      </div>
    );
  }

  return (
    <div className="page container fadeIn">
      <div className={styles.layout}>
        {movie.poster_url && <img className={styles.poster} src={movie.poster_url} alt={movie.title} />}

        <div>
          <h1 className={styles.title}>{movie.title}</h1>
          <div className={styles.metaRow}>
            {movie.genres.map((g) => (
              <span key={g} className={styles.chip}>
                {g}
              </span>
            ))}
            <span className={styles.chip}>{movie.language}</span>
            <span className={styles.chip}>{movie.duration_minutes} min</span>
          </div>
          <p className={styles.description}>{movie.description}</p>

          {movie.status === "coming_soon" ? (
            <div className={`card ${styles.comingSoonNotice}`}>
              Coming soon{movie.release_date ? ` — releasing ${new Date(movie.release_date).toLocaleDateString("en-IN", { day: "numeric", month: "long", year: "numeric" })}` : ""}.
              Check back once it's released to book tickets.
            </div>
          ) : (
            <>
              {!selectedCity && (
                <div className={`card ${styles.comingSoonNotice}`} style={{ marginBottom: 20 }}>
                  Showing all cities. <a onClick={() => navigate("/select-location")} style={{ color: "var(--color-accent)", cursor: "pointer" }}>Pick a city</a> to narrow down theaters near you.
                </div>
              )}

              {shows === null ? (
                <div className="skeleton" style={{ height: 160, borderRadius: 16 }} />
              ) : dateGroups.length === 0 ? (
                <p style={{ color: "var(--color-text-muted)" }}>No showtimes available right now.</p>
              ) : (
                <>
                  <div className={styles.dateTabs}>
                    {dateGroups.map(([date]) => (
                      <button
                        key={date}
                        className={`${styles.dateTab} ${date === activeDate ? styles.active : ""}`}
                        onClick={() => setSelectedDate(date)}
                      >
                        {formatShowDate(date)}
                      </button>
                    ))}
                  </div>

                  {theaterGroups.map(([theaterName, cityName, theaterShows], i) => (
                    <div key={theaterName} className={`card ${styles.theaterGroup} fadeInUp`} style={{ animationDelay: `${i * 40}ms` }}>
                      <div className={styles.theaterHeader}>
                        <span className={styles.theaterName}>{theaterName}</span>
                        <span className={styles.cityTag}>{cityName}</span>
                      </div>
                      <div className={styles.timeGrid}>
                        {theaterShows.map((show) => (
                          <button
                            key={show.id}
                            className={styles.timeSlot}
                            onClick={() => navigate(`/shows/${show.id}/seats`)}
                          >
                            {formatShowTime(show.starts_at)}
                            <span className={styles.timePrice}>{formatMoney(show.price_cents)}</span>
                          </button>
                        ))}
                      </div>
                    </div>
                  ))}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function groupByDate(shows: Show[]): [string, Show[]][] {
  const map = new Map<string, Show[]>();
  for (const show of shows) {
    const date = show.starts_at.slice(0, 10);
    if (!map.has(date)) map.set(date, []);
    map.get(date)!.push(show);
  }
  return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0]));
}

function groupByTheater(shows: Show[]): [string, string, Show[]][] {
  const map = new Map<string, { city: string; shows: Show[] }>();
  for (const show of shows) {
    const key = show.theater_name ?? `Theater ${show.theater_id}`;
    if (!map.has(key)) map.set(key, { city: show.city_name ?? "", shows: [] });
    map.get(key)!.shows.push(show);
  }
  for (const entry of map.values()) {
    entry.shows.sort((a, b) => a.starts_at.localeCompare(b.starts_at));
  }
  return Array.from(map.entries())
    .sort((a, b) => a[0].localeCompare(b[0]))
    .map(([name, { city, shows }]) => [name, city, shows]);
}
