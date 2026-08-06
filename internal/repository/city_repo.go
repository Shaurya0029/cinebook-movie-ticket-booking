package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

type CityRepository struct {
	pool *pgxpool.Pool
}

func NewCityRepository(pool *pgxpool.Pool) *CityRepository {
	return &CityRepository{pool: pool}
}

func (r *CityRepository) List(ctx context.Context) ([]models.City, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(state, ''), latitude, longitude FROM cities ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var c models.City
		if err := rows.Scan(&c.ID, &c.Name, &c.State, &c.Latitude, &c.Longitude); err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, rows.Err()
}

func (r *CityRepository) GetByID(ctx context.Context, id int64) (*models.City, error) {
	var c models.City
	err := r.pool.QueryRow(ctx, `SELECT id, name, COALESCE(state, ''), latitude, longitude FROM cities WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.State, &c.Latitude, &c.Longitude)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
