package httpserver

import (
		"net/http"

		"github.com/go-chi/chi/v5"
		chimiddleware "github.com/go-chi/chi/v5/middleware"
		"github.com/go-chi/cors"

		"movieticketbooking/internal/auth"
		"movieticketbooking/internal/handlers"
	)

func NewRouter(d *handlers.Deps, allowedOrigin string) http.Handler {
		r := chi.NewRouter()

		origins := []string{"http://localhost:5173"}
		if allowedOrigin != "" {
					origins = append(origins, allowedOrigin)
				}

		r.Use(chimiddleware.Logger)
		r.Use(chimiddleware.Recoverer)
		r.Use(cors.Handler(cors.Options{
					AllowedOrigins:   origins,
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
					AllowCredentials: true,
				}))

		r.Get("/api/healthz", handlers.Healthz)

		r.Route("/api/auth", func(r chi.Router) {
					r.Post("/register", d.Register)
					r.Post("/login", d.Login)
					r.With(auth.RequireAuth(d.TokenIssuer)).Get("/me", d.Me)
				})

		r.Route("/api/cities", func(r chi.Router) {
					r.Get("/", d.ListCities)
					r.Get("/nearest", d.NearestCity)
					r.Get("/{cityID}/theaters", d.ListTheatersByCity)
				})

		r.Route("/api/theaters", func(r chi.Router) {
					r.Get("/nearest", d.NearestTheater)
					r.Get("/{theaterID}", d.GetTheater)
				})

		r.Route("/api/movies", func(r chi.Router) {
					r.Get("/", d.ListMovies)
					r.Get("/{movieID}", d.GetMovie)
					r.Get("/{movieID}/shows", d.ListShowsForMovie)
				})

		r.Route("/api/shows", func(r chi.Router) {
					r.Get("/{showID}", d.GetShow)
					r.Get("/{showID}/seats", d.ListSeatsForShow)
				})

		r.Route("/api/bookings", func(r chi.Router) {
					r.Use(auth.RequireAuth(d.TokenIssuer))
					r.Post("/", d.CreateBooking)
					r.Get("/", d.ListMyBookings)
					r.Get("/{bookingID}", d.GetBooking)
					r.Get("/{bookingID}/receipt", d.GetBookingReceipt)
					r.Post("/{bookingID}/pay", d.PayForBooking)
				})

		return r
}
