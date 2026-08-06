package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

type TheaterRepository struct {
	pool *pgxpool.Pool
}

func NewTheaterRepository(pool *pgxpool.Pool) *TheaterRepository {
	return &TheaterRepository{pool: pool}
}

func (r *TheaterRepository) ListByCity(ctx context.Context, cityID int64) ([]models.Theater, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, city_id, name, COALESCE(address, ''), latitude, longitude FROM theaters WHERE city_id = $1 ORDER BY name`,
		cityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theaters []models.Theater
	for rows.Next() {
		var t models.Theater
		if err := rows.Scan(&t.ID, &t.CityID, &t.Name, &t.Address, &t.Latitude, &t.Longitude); err != nil {
			return nil, err
		}
		theaters = append(theaters, t)
	}
	return theaters, rows.Err()
}

func (r *TheaterRepository) GetByID(ctx context.Context, id int64) (*models.Theater, error) {
	var t models.Theater
	err := r.pool.QueryRow(ctx,
		`SELECT id, city_id, name, COALESCE(address, ''), latitude, longitude FROM theaters WHERE id = $1`, id).
		Scan(&t.ID, &t.CityID, &t.Name, &t.Address, &t.Latitude, &t.Longitude)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListAllWithCity returns every theater together with its city name, used for
// nearest-theater lookups where we compute distance in Go over the small seeded set.
func (r *TheaterRepository) ListAllWithCity(ctx context.Context) ([]models.Theater, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, city_id, name, COALESCE(address, ''), latitude, longitude FROM theaters ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theaters []models.Theater
	for rows.Next() {
		var t models.Theater
		if err := rows.Scan(&t.ID, &t.CityID, &t.Name, &t.Address, &t.Latitude, &t.Longitude); err != nil {
			return nil, err
		}
		theaters = append(theaters, t)
	}
	return theaters, rows.Err()
}
