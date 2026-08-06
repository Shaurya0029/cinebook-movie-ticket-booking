package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

type ShowRepository struct {
	pool *pgxpool.Pool
}

func NewShowRepository(pool *pgxpool.Pool) *ShowRepository {
	return &ShowRepository{pool: pool}
}

// ShowFilter narrows down shows for a movie by theater, city, and/or calendar date.
type ShowFilter struct {
	MovieID   int64
	TheaterID int64 // 0 = any
	CityID    int64 // 0 = any
	Date      *time.Time
}

func (r *ShowRepository) ListForMovie(ctx context.Context, f ShowFilter) ([]models.Show, error) {
	query := `
		SELECT s.id, s.movie_id, s.theater_id, s.starts_at, s.price_cents,
		       m.title, t.name, c.name
		FROM shows s
		JOIN movies m ON m.id = s.movie_id
		JOIN theaters t ON t.id = s.theater_id
		JOIN cities c ON c.id = t.city_id
		WHERE s.movie_id = $1`
	args := []any{f.MovieID}

	if f.TheaterID != 0 {
		args = append(args, f.TheaterID)
		query += ` AND s.theater_id = $` + strconv.Itoa(len(args))
	}
	if f.CityID != 0 {
		args = append(args, f.CityID)
		query += ` AND t.city_id = $` + strconv.Itoa(len(args))
	}
	if f.Date != nil {
		args = append(args, *f.Date)
		query += ` AND s.starts_at::date = $` + strconv.Itoa(len(args)) + `::date`
	}
	query += ` ORDER BY s.starts_at`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var s models.Show
		if err := rows.Scan(&s.ID, &s.MovieID, &s.TheaterID, &s.StartsAt, &s.PriceCents, &s.MovieTitle, &s.TheaterName, &s.CityName); err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}
	return shows, rows.Err()
}

func (r *ShowRepository) GetByID(ctx context.Context, id int64) (*models.Show, error) {
	var s models.Show
	err := r.pool.QueryRow(ctx, `
		SELECT s.id, s.movie_id, s.theater_id, s.starts_at, s.price_cents,
		       m.title, t.name, c.name
		FROM shows s
		JOIN movies m ON m.id = s.movie_id
		JOIN theaters t ON t.id = s.theater_id
		JOIN cities c ON c.id = t.city_id
		WHERE s.id = $1`, id).
		Scan(&s.ID, &s.MovieID, &s.TheaterID, &s.StartsAt, &s.PriceCents, &s.MovieTitle, &s.TheaterName, &s.CityName)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
