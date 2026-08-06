package handlers

import (
	"movieticketbooking/internal/auth"
	"movieticketbooking/internal/repository"
)

// Deps bundles everything handlers need: repositories, the JWT issuer, and pricing config.
type Deps struct {
	Cities   *repository.CityRepository
	Theaters *repository.TheaterRepository
	Movies   *repository.MovieRepository
	Shows    *repository.ShowRepository
	Seats    *repository.SeatRepository
	Users    *repository.UserRepository
	Bookings *repository.BookingRepository

	TokenIssuer    *auth.TokenIssuer
	GSTRatePercent float64
}
