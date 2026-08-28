package config

import (
		"os"
		"strconv"

		"github.com/joho/godotenv"
	)

type Config struct {
		DatabaseURL    string
		JWTSecret      string
		Port           string
		GSTRatePercent float64
		AllowedOrigin  string
}

// Load reads a .env file if present (dev convenience) and falls back to
// real environment variables. Sensible defaults are used for anything
// missing so the API can still boot in a fresh checkout.
func Load() Config {
		_ = godotenv.Load()

		return Config{
					DatabaseURL:    getEnv("DATABASE_URL", "postgres://movieapp:movieapp@localhost:5434/movieticketbooking?sslmode=disable"),
					JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me"),
					Port:           getEnv("PORT", "8080"),
					GSTRatePercent: getEnvFloat("GST_RATE_PERCENT", 18.00),
					AllowedOrigin:  getEnv("ALLOWED_ORIGIN", ""),
				}
}

func getEnv(key, fallback string) string {
		if v := os.Getenv(key); v != "" {
					return v
				}
		return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
		v := os.Getenv(key)
		if v == "" {
					return fallback
				}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
					return fallback
				}
		return f
}
