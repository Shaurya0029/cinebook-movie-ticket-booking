# CineBook — Movie Ticket Booking System

A full-stack movie ticket booking app: browse movies currently showing and coming soon, pick a city/theater (including "use current location"), choose a showtime and seats, pay (simulated), and get a receipt. Built as a learning project on top of Go.

## Tech stack / languages used

- **Go** — backend REST API (`net/http` via the `chi` router)
- **TypeScript** — frontend application logic
- **React** (with TSX/JSX) — frontend UI
- **CSS** (CSS Modules) — styling
- **SQL** — PostgreSQL schema and queries
- **HTML** — app entry point (`index.html`)
- **YAML** — `docker-compose.yml`
- **Bash** (shell scripting) — none checked in, but used during setup

Key libraries: `go-chi/chi`, `jackc/pgx`, `golang-jwt/jwt`, `golang.org/x/crypto/bcrypt` on the backend; `react`, `react-router-dom`, `vite` on the frontend.

## Project structure

```
MovieTicketBooking/
  cmd/api/        Go API server entrypoint
  cmd/seed/       Database seed script
  internal/       Go backend packages (auth, db, handlers, models, repository, pricing, geo)
  frontend/       React + TypeScript + Vite frontend
  docker-compose.yml   PostgreSQL for local development
```

## Features

- Browse "Now Showing" and "Coming Soon" movies
- City/theater picker with browser geolocation ("use current location")
- Showtime picker with live seat map (transactional seat locking — no double-booking)
- Price breakdown with GST calculated per booking
- Simulated payment flow (a documented test card always declines, everything else succeeds)
- Booking receipts and a full booking history per user
- JWT-based auth (register/login)

## Running locally

Requires Go, Node.js, and Docker.

```bash
# 1. Start Postgres
docker compose up -d

# 2. Seed demo data (cities, theaters, movies, shows, seats)
go run ./cmd/seed

# 3. Start the API server (http://localhost:8080)
go run ./cmd/api

# 4. Start the frontend (http://localhost:5173)
cd frontend
npm install
npm run dev
```

## Known limitations

This is a learning project, not a production system:

- Payments are entirely simulated — no real payment gateway is integrated.
- If a user selects seats but never completes payment, those seats remain locked with no expiry/cleanup job.
- Location matching uses a small seeded set of cities/theaters with plain-math (Haversine) distance — no real geocoding/maps API.
