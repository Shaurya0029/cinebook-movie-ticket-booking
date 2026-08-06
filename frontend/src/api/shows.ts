import { apiRequest } from "./client";
import type { Seat, Show } from "../types";

interface ShowFilters {
  theaterId?: number;
  cityId?: number;
  date?: string; // YYYY-MM-DD
}

export function listShowsForMovie(movieId: number, filters: ShowFilters = {}) {
  const params = new URLSearchParams();
  if (filters.theaterId) params.set("theater_id", String(filters.theaterId));
  if (filters.cityId) params.set("city_id", String(filters.cityId));
  if (filters.date) params.set("date", filters.date);
  const query = params.toString() ? `?${params.toString()}` : "";
  return apiRequest<Show[]>(`/movies/${movieId}/shows${query}`);
}

export function getShow(showId: number) {
  return apiRequest<Show>(`/shows/${showId}`);
}

export function listSeatsForShow(showId: number) {
  return apiRequest<Seat[]>(`/shows/${showId}/seats`);
}
