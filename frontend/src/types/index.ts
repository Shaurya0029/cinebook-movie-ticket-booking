export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface City {
  id: number;
  name: string;
  state?: string;
  latitude: number;
  longitude: number;
}

export interface Theater {
  id: number;
  city_id: number;
  name: string;
  address?: string;
  latitude: number;
  longitude: number;
}

export type MovieStatus = "now_showing" | "coming_soon";

export interface Movie {
  id: number;
  title: string;
  genres: string[];
  language: string;
  duration_minutes: number;
  description?: string;
  poster_url?: string;
  status: MovieStatus;
  release_date?: string;
}

export interface Show {
  id: number;
  movie_id: number;
  theater_id: number;
  starts_at: string;
  price_cents: number;
  movie_title?: string;
  theater_name?: string;
  city_name?: string;
}

export interface Seat {
  id: number;
  show_id: number;
  seat_label: string;
  is_booked: boolean;
}

export type BookingStatus = "pending_payment" | "confirmed" | "cancelled";

export interface Booking {
  id: number;
  user_id: number;
  show_id: number;
  status: BookingStatus;
  seat_count: number;
  subtotal_cents: number;
  gst_rate_percent: number;
  gst_amount_cents: number;
  total_amount_cents: number;
  created_at: string;
}

export type PaymentStatus = "success" | "failed";

export interface Payment {
  id: number;
  booking_id: number;
  amount_cents: number;
  card_last4: string;
  status: PaymentStatus;
  failure_reason?: string;
  paid_at?: string;
}

export interface Receipt {
  booking: Booking;
  payment?: Payment;
  movie_title: string;
  theater_name: string;
  city_name: string;
  show_starts_at: string;
  price_cents: number;
  user_first_name: string;
  user_last_name: string;
  seat_labels: string[] | null;
}

export interface PayResponse {
  payment: Payment;
  booking_status: BookingStatus;
}

export interface ApiError {
  error: string;
}
