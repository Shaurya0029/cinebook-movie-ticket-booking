import { apiRequest } from "./client";
import type { Movie, MovieStatus } from "../types";

export function listMovies(status?: MovieStatus) {
  const query = status ? `?status=${status}` : "";
  return apiRequest<Movie[]>(`/movies${query}`);
}

export function getMovie(movieId: number) {
  return apiRequest<Movie>(`/movies/${movieId}`);
}
