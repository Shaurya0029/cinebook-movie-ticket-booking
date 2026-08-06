package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

type MovieRepository struct {
	pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{pool: pool}
}

// List returns movies, optionally filtered by status ("now_showing" or "coming_soon").
// An empty status returns all movies.
func (r *MovieRepository) List(ctx context.Context, status string) ([]models.Movie, error) {
	query := `SELECT id, title, genres, language, duration_minutes, COALESCE(description, ''), COALESCE(poster_url, ''), status, release_date FROM movies`
	args := []any{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY release_date DESC NULLS LAST, title`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Genres, &m.Language, &m.DurationMinutes, &m.Description, &m.PosterURL, &m.Status, &m.ReleaseDate); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, rows.Err()
}

func (r *MovieRepository) GetByID(ctx context.Context, id int64) (*models.Movie, error) {
	var m models.Movie
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, genres, language, duration_minutes, COALESCE(description, ''), COALESCE(poster_url, ''), status, release_date FROM movies WHERE id = $1`,
		id).Scan(&m.ID, &m.Title, &m.Genres, &m.Language, &m.DurationMinutes, &m.Description, &m.PosterURL, &m.Status, &m.ReleaseDate)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
