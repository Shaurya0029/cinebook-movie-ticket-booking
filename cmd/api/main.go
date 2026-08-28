package main

import (
		"context"
		"log"
		"net/http"

		"movieticketbooking/internal/auth"
		"movieticketbooking/internal/config"
		"movieticketbooking/internal/db"
		"movieticketbooking/internal/handlers"
		"movieticketbooking/internal/httpserver"
		"movieticketbooking/internal/repository"
	)

func main() {
		cfg := config.Load()
		ctx := context.Background()

		pool, err := db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
					log.Fatalf("failed to connect to database: %v", err)
				}
		defer pool.Close()

		if err := db.EnsureSchema(ctx, pool); err != nil {
					log.Fatalf("failed to ensure schema: %v", err)
				}

		deps := &handlers.Deps{
					Cities:         repository.NewCityRepository(pool),
					Theaters:       repository.NewTheaterRepository(pool),
					Movies:         repository.NewMovieRepository(pool),
					Shows:          repository.NewShowRepository(pool),
					Seats:          repository.NewSeatRepository(pool),
					Users:          repository.NewUserRepository(pool),
					Bookings:       repository.NewBookingRepository(pool),
					TokenIssuer:    auth.NewTokenIssuer(cfg.JWTSecret),
					GSTRatePercent: cfg.GSTRatePercent,
				}

		router := httpserver.NewRouter(deps, cfg.AllowedOrigin)

		log.Printf("movieticketbooking API listening on :%s", cfg.Port)
		if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
					log.Fatalf("server failed: %v", err)
				}
}
