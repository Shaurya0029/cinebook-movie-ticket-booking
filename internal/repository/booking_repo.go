package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
	"movieticketbooking/internal/pricing"
)

var ErrSeatUnavailable = errors.New("one or more selected seats are no longer available")

type BookingRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{pool: pool}
}

// CreateBooking locks the requested seats and inserts a pending_payment booking,
// all inside one transaction so two users can never both grab the same seat.
func (r *BookingRepository) CreateBooking(ctx context.Context, userID, showID int64, seatIDs []int64, priceCents int, gstRatePercent float64) (*models.Booking, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	seats, err := LockSeatsForBooking(ctx, tx, showID, seatIDs)
	if err != nil {
		return nil, fmt.Errorf("lock seats: %w", err)
	}
	if len(seats) != len(seatIDs) {
		return nil, ErrSeatUnavailable
	}
	for _, s := range seats {
		if s.IsBooked {
			return nil, ErrSeatUnavailable
		}
	}

	breakdown := pricing.Calculate(priceCents, len(seatIDs), gstRatePercent)

	var b models.Booking
	err = tx.QueryRow(ctx, `
		INSERT INTO bookings (user_id, show_id, status, seat_count, subtotal_cents, gst_rate_percent, gst_amount_cents, total_amount_cents)
		VALUES ($1, $2, 'pending_payment', $3, $4, $5, $6, $7)
		RETURNING id, user_id, show_id, status, seat_count, subtotal_cents, gst_rate_percent, gst_amount_cents, total_amount_cents, created_at`,
		userID, showID, len(seatIDs), breakdown.SubtotalCents, breakdown.GSTRatePercent, breakdown.GSTAmountCents, breakdown.TotalAmountCents).
		Scan(&b.ID, &b.UserID, &b.ShowID, &b.Status, &b.SeatCount, &b.SubtotalCents, &b.GSTRatePercent, &b.GSTAmountCents, &b.TotalAmountCents, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert booking: %w", err)
	}

	for _, seatID := range seatIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO booking_seats (booking_id, seat_id) VALUES ($1, $2)`, b.ID, seatID); err != nil {
			return nil, fmt.Errorf("insert booking_seats: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return &b, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	var b models.Booking
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, show_id, status, seat_count, subtotal_cents, gst_rate_percent, gst_amount_cents, total_amount_cents, created_at
		FROM bookings WHERE id = $1`, id).
		Scan(&b.ID, &b.UserID, &b.ShowID, &b.Status, &b.SeatCount, &b.SubtotalCents, &b.GSTRatePercent, &b.GSTAmountCents, &b.TotalAmountCents, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) ListForUser(ctx context.Context, userID int64) ([]models.Booking, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, show_id, status, seat_count, subtotal_cents, gst_rate_percent, gst_amount_cents, total_amount_cents, created_at
		FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ShowID, &b.Status, &b.SeatCount, &b.SubtotalCents, &b.GSTRatePercent, &b.GSTAmountCents, &b.TotalAmountCents, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) SeatLabelsForBooking(ctx context.Context, bookingID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.seat_label FROM booking_seats bs
		JOIN seats s ON s.id = bs.seat_id
		WHERE bs.booking_id = $1 ORDER BY s.seat_label`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	labels := []string{}
	for rows.Next() {
		var l string
		if err := rows.Scan(&l); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

// ConfirmBooking marks a booking confirmed and records a successful payment, in one transaction.
func (r *BookingRepository) ConfirmBooking(ctx context.Context, bookingID int64, amountCents int, cardLast4 string) (*models.Payment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE bookings SET status = 'confirmed' WHERE id = $1`, bookingID); err != nil {
		return nil, err
	}

	var p models.Payment
	err = tx.QueryRow(ctx, `
		INSERT INTO payments (booking_id, amount_cents, card_last4, status, paid_at)
		VALUES ($1, $2, $3, 'success', now())
		RETURNING id, booking_id, amount_cents, card_last4, status, COALESCE(failure_reason, ''), paid_at`,
		bookingID, amountCents, cardLast4).
		Scan(&p.ID, &p.BookingID, &p.AmountCents, &p.CardLast4, &p.Status, &p.FailureReason, &p.PaidAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &p, nil
}

// FailBooking cancels a booking, releases its seats, and records a failed payment, in one transaction.
func (r *BookingRepository) FailBooking(ctx context.Context, bookingID int64, amountCents int, cardLast4, reason string) (*models.Payment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE bookings SET status = 'cancelled' WHERE id = $1`, bookingID); err != nil {
		return nil, err
	}
	if err := ReleaseSeatsForBooking(ctx, tx, bookingID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM booking_seats WHERE booking_id = $1`, bookingID); err != nil {
		return nil, err
	}

	var p models.Payment
	err = tx.QueryRow(ctx, `
		INSERT INTO payments (booking_id, amount_cents, card_last4, status, failure_reason)
		VALUES ($1, $2, $3, 'failed', $4)
		RETURNING id, booking_id, amount_cents, card_last4, status, COALESCE(failure_reason, ''), paid_at`,
		bookingID, amountCents, cardLast4, reason).
		Scan(&p.ID, &p.BookingID, &p.AmountCents, &p.CardLast4, &p.Status, &p.FailureReason, &p.PaidAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *BookingRepository) GetPaymentForBooking(ctx context.Context, bookingID int64) (*models.Payment, error) {
	var p models.Payment
	err := r.pool.QueryRow(ctx, `
		SELECT id, booking_id, amount_cents, card_last4, status, COALESCE(failure_reason, ''), paid_at
		FROM payments WHERE booking_id = $1`, bookingID).
		Scan(&p.ID, &p.BookingID, &p.AmountCents, &p.CardLast4, &p.Status, &p.FailureReason, &p.PaidAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
