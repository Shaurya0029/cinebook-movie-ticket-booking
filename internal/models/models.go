package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type City struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	State     string  `json:"state,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Theater struct {
	ID        int64   `json:"id"`
	CityID    int64   `json:"city_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Movie struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Genres          []string   `json:"genres"`
	Language        string     `json:"language"`
	DurationMinutes int        `json:"duration_minutes"`
	Description     string     `json:"description,omitempty"`
	PosterURL       string     `json:"poster_url,omitempty"`
	Status          string     `json:"status"` // now_showing | coming_soon
	ReleaseDate     *time.Time `json:"release_date,omitempty"`
}

type Show struct {
	ID         int64     `json:"id"`
	MovieID    int64     `json:"movie_id"`
	TheaterID  int64     `json:"theater_id"`
	StartsAt   time.Time `json:"starts_at"`
	PriceCents int       `json:"price_cents"`

	// Populated on joined reads for convenience.
	MovieTitle  string `json:"movie_title,omitempty"`
	TheaterName string `json:"theater_name,omitempty"`
	CityName    string `json:"city_name,omitempty"`
}

type Seat struct {
	ID        int64  `json:"id"`
	ShowID    int64  `json:"show_id"`
	SeatLabel string `json:"seat_label"`
	IsBooked  bool   `json:"is_booked"`
}

type Booking struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	ShowID           int64     `json:"show_id"`
	Status           string    `json:"status"` // pending_payment | confirmed | cancelled
	SeatCount        int       `json:"seat_count"`
	SubtotalCents    int       `json:"subtotal_cents"`
	GSTRatePercent   float64   `json:"gst_rate_percent"`
	GSTAmountCents   int       `json:"gst_amount_cents"`
	TotalAmountCents int       `json:"total_amount_cents"`
	CreatedAt        time.Time `json:"created_at"`
}

type Payment struct {
	ID            int64      `json:"id"`
	BookingID     int64      `json:"booking_id"`
	AmountCents   int        `json:"amount_cents"`
	CardLast4     string     `json:"card_last4"`
	Status        string     `json:"status"` // success | failed
	FailureReason string     `json:"failure_reason,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

// Receipt is the fully joined view returned to the client after booking/payment.
type Receipt struct {
	Booking       Booking   `json:"booking"`
	Payment       *Payment  `json:"payment,omitempty"`
	MovieTitle    string    `json:"movie_title"`
	TheaterName   string    `json:"theater_name"`
	CityName      string    `json:"city_name"`
	ShowStartsAt  time.Time `json:"show_starts_at"`
	PriceCents    int       `json:"price_cents"`
	UserFirstName string    `json:"user_first_name"`
	UserLastName  string    `json:"user_last_name"`
	SeatLabels    []string  `json:"seat_labels"`
}
