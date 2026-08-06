CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cities (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    state TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS theaters (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL REFERENCES cities(id),
    name TEXT NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    UNIQUE(city_id, name)
);

CREATE TABLE IF NOT EXISTS movies (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    genres TEXT[] NOT NULL DEFAULT '{}',
    language TEXT NOT NULL,
    duration_minutes INT NOT NULL,
    description TEXT,
    poster_url TEXT,
    status TEXT NOT NULL CHECK (status IN ('now_showing', 'coming_soon')),
    release_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS shows (
    id BIGSERIAL PRIMARY KEY,
    movie_id BIGINT NOT NULL REFERENCES movies(id),
    theater_id BIGINT NOT NULL REFERENCES theaters(id),
    starts_at TIMESTAMPTZ NOT NULL,
    price_cents INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(movie_id, theater_id, starts_at)
);

CREATE TABLE IF NOT EXISTS seats (
    id BIGSERIAL PRIMARY KEY,
    show_id BIGINT NOT NULL REFERENCES shows(id),
    seat_label TEXT NOT NULL,
    is_booked BOOLEAN NOT NULL DEFAULT false,
    UNIQUE(show_id, seat_label)
);

CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    show_id BIGINT NOT NULL REFERENCES shows(id),
    status TEXT NOT NULL CHECK (status IN ('pending_payment', 'confirmed', 'cancelled')) DEFAULT 'pending_payment',
    seat_count INT NOT NULL,
    subtotal_cents INT NOT NULL,
    gst_rate_percent NUMERIC(5,2) NOT NULL,
    gst_amount_cents INT NOT NULL,
    total_amount_cents INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS booking_seats (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id),
    seat_id BIGINT NOT NULL REFERENCES seats(id)
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL UNIQUE REFERENCES bookings(id),
    amount_cents INT NOT NULL,
    card_last4 TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('success', 'failed')),
    failure_reason TEXT,
    paid_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_theaters_city_id ON theaters(city_id);
CREATE INDEX IF NOT EXISTS idx_shows_movie_id ON shows(movie_id);
CREATE INDEX IF NOT EXISTS idx_shows_theater_id ON shows(theater_id);
CREATE INDEX IF NOT EXISTS idx_seats_show_id ON seats(show_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_booking_seats_booking_id ON booking_seats(booking_id);
