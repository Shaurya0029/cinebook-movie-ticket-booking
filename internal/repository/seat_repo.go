package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

type SeatRepository struct {
	pool *pgxpool.Pool
}

func NewSeatRepository(pool *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{pool: pool}
}

func (r *SeatRepository) ListForShow(ctx context.Context, showID int64) ([]models.Seat, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, show_id, seat_label, is_booked FROM seats WHERE show_id = $1 ORDER BY seat_label`, showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		if err := rows.Scan(&s.ID, &s.ShowID, &s.SeatLabel, &s.IsBooked); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, rows.Err()
}

// LockSeatsForBooking re-checks the requested seats are still free and marks them
// booked, all within the caller's transaction, using SELECT ... FOR UPDATE so two
// concurrent bookings can't both grab the same seat.
func LockSeatsForBooking(ctx context.Context, tx pgx.Tx, showID int64, seatIDs []int64) ([]models.Seat, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, show_id, seat_label, is_booked FROM seats WHERE show_id = $1 AND id = ANY($2) FOR UPDATE`,
		showID, seatIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		if err := rows.Scan(&s.ID, &s.ShowID, &s.SeatLabel, &s.IsBooked); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE seats SET is_booked = true WHERE show_id = $1 AND id = ANY($2)`, showID, seatIDs); err != nil {
		return nil, err
	}
	return seats, nil
}

// ReleaseSeatsForBooking marks a booking's seats free again (used on payment failure).
func ReleaseSeatsForBooking(ctx context.Context, tx pgx.Tx, bookingID int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE seats SET is_booked = false
		WHERE id IN (SELECT seat_id FROM booking_seats WHERE booking_id = $1)`, bookingID)
	return err
}
